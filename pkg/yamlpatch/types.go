// Package yamlpatch provides feature to apply a JSON Patch to a YAML document.
//
// It supports both JSON Pointer and JSON Path expressions.
//
package yamlpatch

import (
	"gopkg.in/yaml.v3"
)

// Operation represents an operation to patch.
// This is compatible with JSON Patch defined at https://tools.ietf.org/html/rfc6902.
type Operation struct {
	// Currently supported: "replace"
	Op string `yaml:"op"`

	// An expression of JSON Pointer.
	// Either path or jsonpath must be set.
	// See https://tools.ietf.org/html/rfc6901.
	JSONPointer string `yaml:"path"`

	// An expression of JSON Path.
	// Either path or jsonpath must be set.
	// See https://github.com/vmware-labs/yaml-jsonpath for details.
	JSONPath string `yaml:"jsonpath"`

	Value yaml.Node `yaml:"value"`
}
