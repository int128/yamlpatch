package yamlpatch

import (
	"fmt"

	"github.com/int128/yamlpatch/pkg/jsonpath"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

// Apply applies the set of operations to the YAML node in order.
func Apply(n *yaml.Node, ops []Operation) error {
	for _, operation := range ops {
		if err := apply(n, operation); err != nil {
			return fmt.Errorf("yamlpatch.Apply error: %w", err)
		}
	}
	return nil
}

func apply(n *yaml.Node, o Operation) error {
	if o.Op == "replace" {
		return applyReplace(n, o)
	}
	if o.Op == "remove" {
		return applyRemove(n, o)
	}
	return fmt.Errorf("invalid op %s (currently supported: replace)", o.Op)
}

func applyReplace(n *yaml.Node, o Operation) error {
	if o.Value.IsZero() {
		return fmt.Errorf("missing value in patch (op=%s)", o.Op)
	}

	targetPath, err := compilePath(o)
	if err != nil {
		return fmt.Errorf("invalid patch: %w", err)
	}
	nodes, err := targetPath.Find(n)
	if err != nil {
		return fmt.Errorf("could not find the path in YAML: %w", err)
	}
	if len(nodes) == 0 {
		return fmt.Errorf("any node did not match (path=%s, jsonpath=%s)", o.JSONPointer, o.JSONPath)
	}
	for _, node := range nodes {
		node.Kind = o.Value.Kind
		node.Style = o.Value.Style
		node.Value = o.Value.Value
		node.Tag = o.Value.Tag
		node.Content = o.Value.Content
	}
	return nil
}

func compilePath(o Operation) (*yamlpath.Path, error) {
	if o.JSONPath != "" && o.JSONPointer != "" {
		return nil, fmt.Errorf("do not set both path and jsonpath (path=%s, jsonpath=%s)", o.JSONPointer, o.JSONPath)
	}

	if o.JSONPath != "" {
		compiled, err := yamlpath.NewPath(o.JSONPath)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON Path (jsonpath=%s): %w", o.JSONPath, err)
		}
		return compiled, nil
	}

	jsonPath := jsonpath.FromJSONPointer(o.JSONPointer)
	compiled, err := yamlpath.NewPath(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON Path (path=%s) -> (jsonpath=%s): %w", o.JSONPointer, jsonPath, err)
	}
	return compiled, nil
}
