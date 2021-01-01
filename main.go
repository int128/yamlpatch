package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/int128/yamlpatch/pkg/yamlpatch"
	"gopkg.in/yaml.v3"
)

type options struct {
	patch string
}

func run() error {
	var o options
	flag.StringVar(&o.patch, "p", "", "JSON patch to apply")
	flag.Parse()

	var ops []yamlpatch.Operation
	if err := yaml.Unmarshal([]byte(o.patch), &ops); err != nil {
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
		if err := yamlpatch.Apply(&n, ops); err != nil {
			return fmt.Errorf("could not apply patch: %w", err)
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
