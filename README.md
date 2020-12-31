# yamlpatch

A tool to apply JSON patch to YAML preserving comments.


## Getting Started

```console
% yamlpatch -p '[{ "op": "replace", "jsonpath": "$.spec.template.spec.containers[*].image", "value": "hello" }]' < testdata/fixture1.yaml
# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
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
          image: "hello"
          ports:
            - containerPort: 80
        - name: envoy
          image: "hello"
          command:
            - /bin/sh
            # dummy command
            - -c
            - uname
---
# https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/
apiVersion: v1
kind: Service
metadata:
  name: my-nginx
  labels:
    run: my-nginx
spec:
  ports:
    # http
    - port: 80
      protocol: TCP
    # grpc
    - port: 10000
      protocol: TCP
  selector:
    run: my-nginx
```
