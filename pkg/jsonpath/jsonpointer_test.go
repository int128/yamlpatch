package jsonpath

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

func TestFromJSONPointer(t *testing.T) {
	const input = `{
  "foo": ["bar", "baz"],
  "": 0,
  "a/b": 1,
  "c%d": 2,
  "e^f": 3,
  "g|h": 4,
  "i\\j": 5,
  "k\"l": 6,
  " ": 7,
  "m~n": 8
}`

	type testcase struct {
		jsonPointer string
		want        string
	}
	testcases := []testcase{
		// examples from https://tools.ietf.org/html/rfc6901#section-5
		{`/foo`, `["bar", "baz"]`},
		{`/foo/0`, `"bar"`},
		{`/`, `0`},
		{`/a~1b`, `1`},
		{`/c%d`, `2`},
		{`/e^f`, `3`},
		{`/g|h`, `4`},
		{`/i\j`, `5`},
		{`/k"l`, `6`},
		{`/ `, `7`},
		{`/m~0n`, `8`},
	}

	for _, tc := range testcases {
		t.Run(tc.jsonPointer, func(t *testing.T) {
			var n yaml.Node
			if err := yaml.Unmarshal([]byte(input), &n); err != nil {
				t.Fatalf("unmarshal error: %s", err)
			}
			jsonPath := FromJSONPointer(tc.jsonPointer)
			p, err := yamlpath.NewPath(jsonPath)
			if err != nil {
				t.Fatalf("invalid path in patch: %s", err)
			}
			found, err := p.Find(&n)
			if err != nil {
				t.Errorf("could not find the path in YAML: %s", err)
			}
			if len(found) == 0 {
				t.Fatalf("node not found")
			}
			if len(found) > 1 {
				t.Errorf("found wants [1] but was %+v", found)
			}
			b, err := yaml.Marshal(found[0])
			if err != nil {
				t.Errorf("could not marshal the node: %s", err)
			}
			got := strings.TrimSpace(string(b))
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("mismatch (-got +want)\n%s", diff)
			}
		})
	}
}
