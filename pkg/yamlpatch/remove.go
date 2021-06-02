package yamlpatch

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func applyRemove(n *yaml.Node, o Operation) error {
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
	removeNodesFromYAML(n, nodes)
	return nil
}

func removeNodesFromYAML(n *yaml.Node, removeNodes []*yaml.Node) {
	if n.Content == nil || len(n.Content) == 0 {
		return
	}

	var newContent []*yaml.Node
	for _, child := range n.Content {
		if containsNode(removeNodes, child) {
			continue
		}
		newContent = append(newContent, child)
	}
	n.Content = newContent

	for _, child := range n.Content {
		removeNodesFromYAML(child, removeNodes)
	}
}

func containsNode(nodes []*yaml.Node, target *yaml.Node) bool {
	for _, node := range nodes {
		if node == target {
			return true
		}
	}
	return false
}
