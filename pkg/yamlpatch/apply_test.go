package yamlpatch

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func TestApply_jsonpatch_conformance(t *testing.T) {
	f, err := os.Open("testdata/jsonpatch_conformance_tests.yaml")
	if err != nil {
		t.Fatalf("could not open testdata: %s", err)
	}
	defer f.Close()

	type testcase struct {
		Doc      yaml.Node   `yaml:"doc"`
		Patch    []Operation `yaml:"patch"`
		Expected yaml.Node   `yaml:"expected"`
		Error    string      `yaml:"error"`
		Comment  string      `yaml:"comment"`
	}
	var testcases []testcase
	if err := yaml.NewDecoder(f).Decode(&testcases); err != nil {
		t.Fatalf("could not decode testdata: %s", err)
	}

	cmpOptsNode := cmpopts.IgnoreFields(yaml.Node{}, "Line", "Column")
	for _, tc := range testcases {
		t.Run(tc.Comment, func(t *testing.T) {
			err := Apply(&tc.Doc, tc.Patch)
			if tc.Error != "" {
				if err == nil {
					t.Errorf("apply wants error but was nil (%s)", tc.Error)
				}
				return
			}
			if err != nil {
				t.Fatalf("apply wants non-error but was error: %s", err)
			}
			if diff := cmp.Diff(tc.Doc, tc.Expected, cmpOptsNode); diff != "" {
				t.Errorf("mismatch (-got +want)\n%s", diff)
			}
		})
	}
}
