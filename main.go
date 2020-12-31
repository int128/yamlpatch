package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

type jsonPatchOperation struct {
	Op       string    `yaml:"op"`
	JSONPath string    `yaml:"jsonpath"`
	Value    yaml.Node `yaml:"value"`
}

func applyJSONPatch(n *yaml.Node, o jsonPatchOperation) error {
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
		// TODO: fix for a content node
		node.Kind = o.Value.Kind
		node.Style = o.Value.Style
		//log.Printf("%#v -> %#v", node.Value, o.Value.Value)
		node.Value = o.Value.Value
		//log.Printf("%#v -> %#v", node.Content, o.Value.Content)
		node.Content = o.Value.Content
	}
	return nil
}

type options struct {
	jsonPatch string
}

func run() error {
	var o options
	flag.StringVar(&o.jsonPatch, "p", "", "JSON patch to apply")
	flag.Parse()

	var jsonPatch []jsonPatchOperation
	if err := yaml.Unmarshal([]byte(o.jsonPatch), &jsonPatch); err != nil {
		return fmt.Errorf("invalid patch: %w", err)
	}

	d := yaml.NewDecoder(os.Stdin)
	e := yaml.NewEncoder(os.Stdout)
	e.SetIndent(2)
	for {
		var n yaml.Node
		if err := d.Decode(&n); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("could not decode YAML: %w", err)
		}
		for _, operation := range jsonPatch {
			if err := applyJSONPatch(&n, operation); err != nil {
				return fmt.Errorf("could not apply patch: %w", err)
			}
		}
		if err := e.Encode(&n); err != nil {
			return fmt.Errorf("could not write YAML: %w", err)
		}
	}
}

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	if err := run(); err != nil {
		log.Printf("error: %s", err)
	}
}
