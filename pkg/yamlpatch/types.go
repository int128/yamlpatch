// Package yamlpatch provides feature to patch to a YAML node.
package yamlpatch

import (
	"gopkg.in/yaml.v3"
)

// Operation represents an operation to patch.
// This is a subset of JSON Patch defined at https://tools.ietf.org/html/rfc6902.
type Operation struct {
	// Currently supported: "replace"
	Op string `yaml:"op"`

	// JSON pointer.
	// Either path or jsonpath must be set.
	// See https://tools.ietf.org/html/rfc6901.
	JSONPointer string `yaml:"path"`

	// JSON path.
	// Either path or jsonpath must be set.
	// See https://github.com/vmware-labs/yaml-jsonpath for details.
	JSONPath string `yaml:"jsonpath"`

	Value yaml.Node `yaml:"value"`
}
