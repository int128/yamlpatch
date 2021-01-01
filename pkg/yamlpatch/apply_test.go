package yamlpatch

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestApply(t *testing.T) {
	const deploymentFixture = `# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3 # at least 3
  template:
    spec:
      containers:
        - name: nginx
          # example
          image: nginx:1.14.2
`

	type testcase struct {
		input string
		patch string
		want  string
	}
	testcases := map[string]testcase{
		"replaceScalarString": {
			input: deploymentFixture,
			patch: `
- op: replace
  jsonpath: $.spec.template.spec.containers[*].image
  value: foo
`,
			want: `# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3 # at least 3
  template:
    spec:
      containers:
        - name: nginx
          # example
          image: foo
`,
		},
		"replaceScalerInt": {
			input: deploymentFixture,
			patch: `
- op: replace
  jsonpath: $.spec.replicas
  value: 100
`,
			want: `# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 100 # at least 3
  template:
    spec:
      containers:
        - name: nginx
          # example
          image: nginx:1.14.2
`,
		},
		"replaceMapping": {
			input: deploymentFixture,
			patch: `
- op: replace
  jsonpath: $.spec.template.spec
  value:
    containers:
      - image: busybox
`,
			want: `# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3 # at least 3
  template:
    spec:
      containers:
        - image: busybox
`,
		},
		"replaceSequence": {
			input: deploymentFixture,
			patch: `
- op: replace
  jsonpath: $.spec.template.spec.containers
  value:
    - image: busybox
    - image: envoyproxy/envoy:v1.16-latest
`,
			want: `# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3 # at least 3
  template:
    spec:
      containers:
        - image: busybox
        - image: envoyproxy/envoy:v1.16-latest
`,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			var n yaml.Node
			if err := yaml.Unmarshal([]byte(tc.input), &n); err != nil {
				t.Fatalf("unmarshal error: %s", err)
			}
			var ops []Operation
			if err := yaml.Unmarshal([]byte(tc.patch), &ops); err != nil {
				t.Fatalf("unmarshal error: %s", err)
			}
			if err := Apply(&n, ops); err != nil {
				t.Errorf("apply error: %s", err)
			}
			var b strings.Builder
			e := yaml.NewEncoder(&b)
			e.SetIndent(2)
			if err := e.Encode(&n); err != nil {
				t.Fatalf("marshal error: %s", err)
			}
			got := b.String()
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("mismatch (-got +want)\n%s", diff)
			}
		})
	}
}
