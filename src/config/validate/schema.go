package validate

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/xeipuuv/gojsonschema"
	yamlv3 "gopkg.in/yaml.v3"
)

type LineResolver struct {
	root  *yamlv3.Node
	lines []string
}

func NewLineResolver(data []byte) (*LineResolver, error) {
	var root yamlv3.Node
	if err := yamlv3.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	return &LineResolver{
		root:  &root,
		lines: strings.Split(text, "\n"),
	}, nil
}

func (r *LineResolver) Annotate(path, message string) string {
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

func NormalizeSchemaPath(path string) string {
	if path == "(root)" {
		return ""
	}
	return path
}

func ValidateSchemaBytes(schema, data []byte, source string) error {
	var document any
	if err := yaml.Unmarshal(data, &document); err != nil {
		return fmt.Errorf("failed to parse config file: %s", err)
	}

	result, err := gojsonschema.Validate(
		gojsonschema.NewBytesLoader(schema),
		gojsonschema.NewGoLoader(document),
	)
	if err != nil {
		return err
	}

	if result.Valid() {
		return nil
	}

	lineResolver, lineErr := NewLineResolver(data)
	if lineErr != nil {
		lineResolver = nil
	}

	validationErrors := make([]string, 0, len(result.Errors()))
	for _, item := range result.Errors() {
		fieldPath := NormalizeSchemaPath(item.Field())
		if lineResolver != nil {
			validationErrors = append(validationErrors, lineResolver.Annotate(fieldPath, item.String()))
			continue
		}
		validationErrors = append(validationErrors, item.String())
	}
	slices.Sort(validationErrors)

	return fmt.Errorf("config schema validation failed (%s):\n- %s", source, strings.Join(validationErrors, "\n- "))
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
