package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/xeipuuv/gojsonschema"
	yamlv3 "gopkg.in/yaml.v3"
)

//go:embed schema.json
var configSchema []byte

func ValidateConfig(configPath string) error {
	if err := validateSchemaStrict(configPath); err != nil {
		return err
	}

	aliae, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	resolvedYAML, err := renderResolvedYAML(aliae)
	if err != nil {
		return err
	}

	lineResolver, err := newYAMLLineResolver(resolvedYAML)
	if err != nil {
		return err
	}

	result, err := gojsonschema.Validate(
		gojsonschema.NewBytesLoader(configSchema),
		gojsonschema.NewGoLoader(buildSchemaDocument(aliae)),
	)
	if err != nil {
		return err
	}

	if result.Valid() {
		if err := aliae.computeVars(lineResolver); err != nil {
			return err
		}

		if err := validateIfExpressions(aliae, lineResolver); err != nil {
			return err
		}

		if err := validateProgress(aliae, lineResolver); err != nil {
			return err
		}

		if err := validateScriptWeights(aliae, lineResolver); err != nil {
			return err
		}

		return validateScriptStates(aliae, lineResolver)
	}

	validationErrors := make([]string, 0, len(result.Errors()))
	for _, item := range result.Errors() {
		fieldPath := normalizeSchemaPath(item.Field())
		validationErrors = append(validationErrors, lineResolver.annotate(fieldPath, item.String()))
	}
	slices.Sort(validationErrors)

	return fmt.Errorf("config schema validation failed:\n- %s", strings.Join(validationErrors, "\n- "))
}

func validateSchemaStrict(configPath string) error {
	resolvedConfigPath := resolveConfigPath(configPath)
	if isRemoteConfigPath(resolvedConfigPath) {
		aliae, err := LoadConfig(configPath)
		if err != nil {
			return err
		}

		resolvedYAML, err := renderResolvedYAML(aliae)
		if err != nil {
			return err
		}

		return validateSchemaBytes(resolvedYAML, resolvedConfigPath)
	}

	return validateSchemaStrictLocalRecursive(resolvedConfigPath, nil, 1)
}

func validateSchemaStrictLocalRecursive(configPath string, stack []string, depth int) error {
	if depth > maxExtendsDepth {
		return fmt.Errorf("extends depth limit exceeded: max %d", maxExtendsDepth)
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		absPath = configPath
	}

	absPath = filepath.Clean(absPath)
	if slices.Contains(stack, absPath) {
		return fmt.Errorf("extends cycle detected for %s", absPath)
	}

	if stat, statErr := os.Stat(absPath); os.IsNotExist(statErr) || (statErr == nil && stat.IsDir()) {
		return fmt.Errorf("config file not found: %s", absPath)
	}

	previousConfigPath := configPathCache
	configPathCache = absPath
	defer func() {
		configPathCache = previousConfigPath
	}()

	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	includedData, err := includeUnmarshaler(data)
	if err != nil {
		return err
	}

	if err := validateSchemaBytes(includedData, absPath); err != nil {
		return err
	}

	extends, err := parseExtends(includedData)
	if err != nil {
		return err
	}
	lineResolver, lineErr := newYAMLLineResolver(includedData)
	if lineErr != nil {
		lineResolver = nil
	}
	if err := validateExtendsIfExpressions(extends, lineResolver); err != nil {
		return err
	}

	nextStack := make([]string, len(stack)+1)
	copy(nextStack, stack)
	nextStack[len(stack)] = absPath

	for _, item := range extends {
		if item.If.Ignore() {
			continue
		}

		paths, pathErr := resolveExtendsPaths(item, absPath)
		if pathErr != nil {
			return pathErr
		}

		for _, path := range paths {
			if err := validateSchemaStrictLocalRecursive(path, nextStack, depth+1); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateSchemaBytes(data []byte, source string) error {
	var document any
	if err := yaml.Unmarshal(data, &document); err != nil {
		return fmt.Errorf("failed to parse config file: %s", err)
	}

	result, err := gojsonschema.Validate(
		gojsonschema.NewBytesLoader(configSchema),
		gojsonschema.NewGoLoader(document),
	)
	if err != nil {
		return err
	}

	if result.Valid() {
		return nil
	}

	lineResolver, lineErr := newYAMLLineResolver(data)
	if lineErr != nil {
		lineResolver = nil
	}

	validationErrors := make([]string, 0, len(result.Errors()))
	for _, item := range result.Errors() {
		fieldPath := normalizeSchemaPath(item.Field())
		if lineResolver != nil {
			validationErrors = append(validationErrors, lineResolver.annotate(fieldPath, item.String()))
			continue
		}
		validationErrors = append(validationErrors, item.String())
	}
	slices.Sort(validationErrors)

	return fmt.Errorf("config schema validation failed (%s):\n- %s", source, strings.Join(validationErrors, "\n- "))
}

func validateIfExpressions(aliae *Aliae, lineResolver *yamlLineResolver) error {
	collector := newValidationCollector(lineResolver)
	ifExpressionValidationVisitor{}.Visit(aliae, collector)
	return collector.err("config if expression validation failed")
}

func validateExtendsIfExpressions(extends []Extend, lineResolver *yamlLineResolver) error {
	collector := newValidationCollector(lineResolver)
	extendsIfValidationVisitor{}.Visit(&Aliae{Extends: extends}, collector)
	return collector.err("config if expression validation failed")
}

func renderResolvedYAML(aliae *Aliae) ([]byte, error) {
	if aliae == nil {
		return []byte{}, nil
	}

	data, err := yaml.Marshal(aliae)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func normalizeSchemaPath(path string) string {
	if path == "(root)" {
		return ""
	}
	return path
}

type yamlLineResolver struct {
	root  *yamlv3.Node
	lines []string
}

func newYAMLLineResolver(data []byte) (*yamlLineResolver, error) {
	var root yamlv3.Node
	if err := yamlv3.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	return &yamlLineResolver{
		root:  &root,
		lines: strings.Split(text, "\n"),
	}, nil
}

func (r *yamlLineResolver) annotate(path, message string) string {
	if r == nil {
		return message
	}

	node := findYAMLNode(r.root, path)
	if node == nil || node.Line <= 0 || node.Line > len(r.lines) {
		return message
	}

	lineText := strings.TrimSpace(r.lines[node.Line-1])
	return fmt.Sprintf("%s (line %d: %s)", message, node.Line, lineText)
}

func findYAMLNode(root *yamlv3.Node, path string) *yamlv3.Node {
	if root == nil {
		return nil
	}

	current := root
	if current.Kind == yamlv3.DocumentNode && len(current.Content) > 0 {
		current = current.Content[0]
	}

	if len(path) == 0 {
		return current
	}

	for _, token := range strings.Split(path, ".") {
		if current == nil {
			return nil
		}

		switch current.Kind {
		case yamlv3.DocumentNode:
			if len(current.Content) == 0 {
				return nil
			}
			current = current.Content[0]
		case yamlv3.MappingNode:
			next := (*yamlv3.Node)(nil)
			for i := 0; i < len(current.Content)-1; i += 2 {
				if current.Content[i].Value == token {
					next = current.Content[i+1]
					break
				}
			}
			current = next
		case yamlv3.SequenceNode:
			index, err := strconv.Atoi(token)
			if err != nil || index < 0 || index >= len(current.Content) {
				return nil
			}
			current = current.Content[index]
		case yamlv3.ScalarNode, yamlv3.AliasNode:
			return nil
		default:
			return nil
		}
	}

	return current
}

func progressSchemaItem(aliae *Aliae) (map[string]any, bool) {
	if aliae == nil || !aliae.Progress.Enabled {
		return nil, false
	}

	end := any(aliae.Progress.EndPercentage.Value)
	if aliae.Progress.EndPercentage.Reset {
		end = "reset"
	}

	return map[string]any{
		"start_percentage": aliae.Progress.StartPercentage,
		"internal":         aliae.Progress.Internal,
		"end_percentage":   end,
	}, true
}

func validateProgress(aliae *Aliae, lineResolver *yamlLineResolver) error {
	collector := newValidationCollector(lineResolver)
	progressValidationVisitor{}.Visit(aliae, collector)
	return collector.err("config progress validation failed")
}

func validateScriptWeights(aliae *Aliae, lineResolver *yamlLineResolver) error {
	collector := newValidationCollector(lineResolver)
	scriptWeightValidationVisitor{}.Visit(aliae, collector)
	return collector.err("config script validation failed")
}

func validateScriptStates(aliae *Aliae, lineResolver *yamlLineResolver) error {
	collector := newValidationCollector(lineResolver)
	scriptStateValidationVisitor{}.Visit(aliae, collector)
	return collector.err("config script state validation failed")
}
