package cmd

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/kyaml/kio"
)

func TestFind(t *testing.T) {
	input := bytes.NewReader([]byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
  app: nginx
spec:
  replicas: 3
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
    image: nginx:1.7.9
    ports:
    - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  selector:
  app: nginx
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
`))

	var output bytes.Buffer

	Find(&kio.ByteReader{Reader: input}, BuildFilters("nginx", "Deployment"), &kio.ByteWriter{Writer: &output})

	data := make(map[interface{}]interface{})
	err := yaml.Unmarshal(output.Bytes(), &data)

	if err != nil {
		t.Error("Error parsing YAML")
	}

	want := "Deployment"

	if data["kind"] != "Deployment" {
		t.Errorf("s.Find() = %q, want %q", data["kind"], want)
	}
}
