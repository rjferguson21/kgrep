package cmd

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/kyaml/kio"
)

func TestFindName(t *testing.T) {
	input := bytes.NewReader([]byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: dep-nginx
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

	Find(
		&kio.ByteReader{Reader: input},
		BuildFilters("dep-nginx", ""),
		BuildOutputs(false, &output))

	data := make(map[interface{}]interface{})
	err := yaml.Unmarshal(output.Bytes(), &data)

	if err != nil {
		t.Error("Error parsing YAML")
	}

	want := "dep-nginx"
	got := data["metadata"].(map[interface{}]interface{})["name"]

	if got != want {
		t.Errorf("s.Find() = %q, want %q", got, want)
	}
}

func TestFindKind(t *testing.T) {
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

	Find(
		&kio.ByteReader{Reader: input},
		BuildFilters("", "Deployment"),
		BuildOutputs(false, &output))

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

func TestFindSummary(t *testing.T) {
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

	Find(
		&kio.ByteReader{Reader: input},
		BuildFilters("", "Service"),
		BuildOutputs(true, &output))

	want := "Service/nginx\n"
	got := output.String()

	if got != want {
		t.Errorf("s.Find() = %q, want %q", got, want)
	}
}
