package yamlpatch_test

import (
	"os"

	"github.com/int128/yamlpatch/pkg/yamlpatch"
	"gopkg.in/yaml.v3"
)

func ExampleApply() {
	const input = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3 # at least 3
  template:
    spec:
      containers:
        - name: nginx # just an example
          image: nginx
`
	const patch = `
- op: replace
  jsonpath: $.spec.template.spec.containers[*].image
  value: nginx:1.14
`
	// unmarshal the input and patch
	var n yaml.Node
	if err := yaml.Unmarshal([]byte(input), &n); err != nil {
		panic(err)
	}
	var ops []yamlpatch.Operation
	if err := yaml.Unmarshal([]byte(patch), &ops); err != nil {
		panic(err)
	}
	// apply the patch
	if err := yamlpatch.Apply(&n, ops); err != nil {
		panic(err)
	}
	// write the result
	e := yaml.NewEncoder(os.Stdout)
	e.SetIndent(2)
	if err := e.Encode(&n); err != nil {
		panic(err)
	}
	// Output:
	// apiVersion: apps/v1
	// kind: Deployment
	// metadata:
	//   name: nginx-deployment
	// spec:
	//   replicas: 3 # at least 3
	//   template:
	//     spec:
	//       containers:
	//         - name: nginx # just an example
	//           image: nginx:1.14
}
