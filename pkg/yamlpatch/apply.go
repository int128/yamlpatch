package yamlpatch

import (
	"fmt"

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
	if o.Op != "replace" {
		return fmt.Errorf("invalid op %s (currently supported: replace)", o.Op)
	}

	path, err := yamlpath.NewPath(o.JSONPath)
	if err != nil {
		return fmt.Errorf("invalid path in patch: %w", err)
	}
	nodes, err := path.Find(n)
	if err != nil {
		return fmt.Errorf("could not find the path in YAML: %w", err)
	}
	for _, node := range nodes {
		node.Kind = o.Value.Kind
		node.Style = o.Value.Style
		node.Value = o.Value.Value
		node.Content = o.Value.Content
	}
	return nil
}
