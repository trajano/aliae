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
		gojsonschema.NewGoLoader(schemaDocument(aliae)),
	)
	if err != nil {
		return err
	}

	if result.Valid() {
		if err := validateIfExpressions(aliae, lineResolver); err != nil {
			return err
		}

		if err := validateProgress(aliae, lineResolver); err != nil {
			return err
		}

		return validateScriptWeights(aliae, lineResolver)
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

	nextStack := make([]string, len(stack)+1)
	copy(nextStack, stack)
	nextStack[len(stack)] = absPath

	for _, item := range extends {
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
	if aliae == nil {
		return nil
	}

	validationErrors := make([]string, 0)
	for i, alias := range aliae.Aliae {
		if alias == nil {
			continue
		}
		if err := alias.If.Validate(); err != nil {
			path := fmt.Sprintf("alias.%d.if", i)
			validationErrors = append(validationErrors, lineResolver.annotate(path, fmt.Sprintf("alias[%d].if: %s", i, err)))
		}
	}

	for i, env := range aliae.Envs {
		if env == nil {
			continue
		}
		if err := env.If.Validate(); err != nil {
			path := fmt.Sprintf("env.%d.if", i)
			validationErrors = append(validationErrors, lineResolver.annotate(path, fmt.Sprintf("env[%d].if: %s", i, err)))
		}
	}

	for i, path := range aliae.Paths {
		if path == nil {
			continue
		}
		if err := path.If.Validate(); err != nil {
			pathField := fmt.Sprintf("path.%d.if", i)
			validationErrors = append(validationErrors, lineResolver.annotate(pathField, fmt.Sprintf("path[%d].if: %s", i, err)))
		}
	}

	for i, path := range aliae.CDPaths {
		if path == nil {
			continue
		}
		if err := path.If.Validate(); err != nil {
			pathField := fmt.Sprintf("cdpath.%d.if", i)
			validationErrors = append(validationErrors, lineResolver.annotate(pathField, fmt.Sprintf("cdpath[%d].if: %s", i, err)))
		}
	}

	for i, script := range aliae.Scripts {
		if script == nil {
			continue
		}
		if err := script.If.Validate(); err != nil {
			path := fmt.Sprintf("script.%d.if", i)
			validationErrors = append(validationErrors, lineResolver.annotate(path, fmt.Sprintf("script[%d].if: %s", i, err)))
		}
	}

	for i, link := range aliae.Links {
		if link == nil {
			continue
		}
		if err := link.If.Validate(); err != nil {
			path := fmt.Sprintf("link.%d.if", i)
			validationErrors = append(validationErrors, lineResolver.annotate(path, fmt.Sprintf("link[%d].if: %s", i, err)))
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	slices.Sort(validationErrors)
	return fmt.Errorf("config if expression validation failed:\n- %s", strings.Join(validationErrors, "\n- "))
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

func schemaDocument(aliae *Aliae) map[string]any {
	document := make(map[string]any)
	if aliae == nil {
		return document
	}

	if aliases := aliasSchemaItems(aliae); len(aliases) > 0 {
		document["alias"] = aliases
	}

	if envs := envSchemaItems(aliae); len(envs) > 0 {
		document["env"] = envs
	}

	if paths := pathSchemaItems(aliae); len(paths) > 0 {
		document["path"] = paths
	}

	if paths := cdpathSchemaItems(aliae); len(paths) > 0 {
		document["cdpath"] = paths
	}

	if scripts := scriptSchemaItems(aliae); len(scripts) > 0 {
		document["script"] = scripts
	}

	if links := linkSchemaItems(aliae); len(links) > 0 {
		document["link"] = links
	}

	if progress, ok := progressSchemaItem(aliae); ok {
		document["progress"] = progress
	}

	if aliae.StatTimeout > 0 {
		document["stat_timeout"] = aliae.StatTimeout.String()
	}

	return document
}

func aliasSchemaItems(aliae *Aliae) []map[string]any {
	aliases := make([]map[string]any, 0, len(aliae.Aliae))
	for _, alias := range aliae.Aliae {
		if alias == nil {
			continue
		}

		item := map[string]any{
			"name":  alias.Name,
			"value": string(alias.Value),
		}
		if len(alias.Type) > 0 {
			item["type"] = string(alias.Type)
		}
		if len(alias.If) > 0 {
			item["if"] = string(alias.If)
		}
		aliases = append(aliases, item)
	}
	return aliases
}

func envSchemaItems(aliae *Aliae) []map[string]any {
	envs := make([]map[string]any, 0, len(aliae.Envs))
	for _, env := range aliae.Envs {
		if env == nil {
			continue
		}

		item := map[string]any{
			"name":  env.Name,
			"value": env.Value,
		}
		if len(env.Delimiter) > 0 {
			item["delimiter"] = string(env.Delimiter)
		}
		if len(env.If) > 0 {
			item["if"] = string(env.If)
		}
		if len(env.Type) > 0 {
			item["type"] = string(env.Type)
		}
		if env.IsPath {
			item["isPath"] = true
		}
		if env.IfExists {
			item["ifExists"] = true
		}
		if env.Persist {
			item["persist"] = true
		}
		envs = append(envs, item)
	}
	return envs
}

func pathSchemaItems(aliae *Aliae) []map[string]any {
	paths := make([]map[string]any, 0, len(aliae.Paths))
	for _, path := range aliae.Paths {
		if path == nil {
			continue
		}

		item := map[string]any{
			"value": string(path.Value),
		}
		if len(path.If) > 0 {
			item["if"] = string(path.If)
		}
		if path.Persist {
			item["persist"] = true
		}
		if path.Force {
			item["force"] = true
		}
		if path.IfExists {
			item["ifExists"] = true
		}
		paths = append(paths, item)
	}
	return paths
}

func cdpathSchemaItems(aliae *Aliae) []map[string]any {
	paths := make([]map[string]any, 0, len(aliae.CDPaths))
	for _, path := range aliae.CDPaths {
		if path == nil {
			continue
		}

		item := map[string]any{
			"value": string(path.Value),
		}
		if len(path.If) > 0 {
			item["if"] = string(path.If)
		}
		if path.Force {
			item["force"] = true
		}
		if path.IfExists {
			item["ifExists"] = true
		}
		paths = append(paths, item)
	}
	return paths
}

func scriptSchemaItems(aliae *Aliae) []map[string]any {
	scripts := make([]map[string]any, 0, len(aliae.Scripts))
	for _, script := range aliae.Scripts {
		if script == nil {
			continue
		}

		item := map[string]any{
			"value": string(script.Value),
		}
		if len(script.If) > 0 {
			item["if"] = string(script.If)
		}
		if script.Weight > 0 {
			item["weight"] = script.Weight
		}
		scripts = append(scripts, item)
	}
	return scripts
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
		"end_percentage":   end,
	}, true
}

func validateProgress(aliae *Aliae, lineResolver *yamlLineResolver) error {
	if aliae == nil || !aliae.Progress.Enabled {
		return nil
	}

	start := aliae.Progress.StartPercentage
	end := aliae.Progress.EndPercentage.Value

	validationErrors := make([]string, 0, 3)
	if start < 0 || start > 100 {
		validationErrors = append(
			validationErrors,
			lineResolver.annotate("progress.start_percentage", "progress.start_percentage must be between 0 and 100"),
		)
	}

	if end < 0 || end > 100 {
		validationErrors = append(
			validationErrors,
			lineResolver.annotate("progress.end_percentage", "progress.end_percentage must be between 0 and 100"),
		)
	}

	if !aliae.Progress.EndPercentage.Reset && end <= start {
		validationErrors = append(
			validationErrors,
			lineResolver.annotate("progress.end_percentage", "progress.end_percentage must be greater than progress.start_percentage"),
		)
	}

	if len(validationErrors) == 0 {
		return nil
	}

	slices.Sort(validationErrors)
	return fmt.Errorf("config progress validation failed:\n- %s", strings.Join(validationErrors, "\n- "))
}

func validateScriptWeights(aliae *Aliae, lineResolver *yamlLineResolver) error {
	if aliae == nil {
		return nil
	}

	validationErrors := make([]string, 0)
	for i, script := range aliae.Scripts {
		if script == nil || script.Weight == 0 {
			continue
		}

		if script.Weight < 1 {
			path := fmt.Sprintf("script.%d.weight", i)
			validationErrors = append(
				validationErrors,
				lineResolver.annotate(path, fmt.Sprintf("script[%d].weight must be greater than or equal to 1", i)),
			)
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	slices.Sort(validationErrors)
	return fmt.Errorf("config script validation failed:\n- %s", strings.Join(validationErrors, "\n- "))
}

func linkSchemaItems(aliae *Aliae) []map[string]any {
	links := make([]map[string]any, 0, len(aliae.Links))
	for _, link := range aliae.Links {
		if link == nil {
			continue
		}

		item := map[string]any{
			"name":   string(link.Name),
			"target": string(link.Target),
		}
		if len(link.If) > 0 {
			item["if"] = string(link.If)
		}
		if link.MkDir {
			item["mkdir"] = true
		}
		links = append(links, item)
	}
	return links
}
