package yamlpatch

import (
	"gopkg.in/yaml.v3"
)

type Operation struct {
	Op       string    `yaml:"op"`
	JSONPath string    `yaml:"jsonpath"`
	Value    yaml.Node `yaml:"value"`
}
