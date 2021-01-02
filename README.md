# yamlpatch ![test](https://github.com/int128/yamlpatch/workflows/test/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/int128/yamlpatch/pkg/yamlpatch.svg)](https://pkg.go.dev/github.com/int128/yamlpatch/pkg/yamlpatch)

This is a command line tool to apply a JSON Patch to a YAML Document preserving position and comments.


## Features

- Support both JSON Pointer and JSON Path (depends on [vmware-labs/yaml-jsonpath](https://github.com/vmware-labs/yaml-jsonpath))
- Passed the [conformance tests](https://github.com/json-patch/json-patch-tests) of JSON Patch
- Single binary

**Note**: currently only `op=replace` mode is implemented


## Getting Started

TODO: install

### Example: Replace a field in Kubernetes YAML

Input:

```yaml
# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3 # at least 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
        - name: envoy
          image: envoyproxy/envoy:v1.16.2
          args:
            - --bootstrap-version
            - "3" # required for v3 API
```

Apply a patch:

```sh
yamlpatch -p '[{ "op": "replace", "jsonpath": "$.spec.template.spec.containers[0].image", "value": nginx:1.19 }]' < testdata/fixture1.yaml
```

Result:

```yaml
# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3 # at least 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.19
        - name: envoy
          image: envoyproxy/envoy:v1.16.2
          args:
            - --bootstrap-version
            - "3" # required for v3 API
```
