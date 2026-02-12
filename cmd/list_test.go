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

	want := "\033[34mService\033[0m/\033[32mnginx\033[0m\n"
	got := output.String()

	if got != want {
		t.Errorf("s.Find() = %q, want %q", got, want)
	}
}

func TestFindPartialMatch(t *testing.T) {
	input := bytes.NewReader([]byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginxDeployment
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
		BuildFilters("nginx$", ""),
		BuildOutputs(true, &output))

	want := "\033[34mService\033[0m/\033[32mnginx\033[0m\n"
	got := output.String()

	if got != want {
		t.Errorf("s.Find() = %q, want %q", got, want)
	}
}

func TestParsePositionalArg(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		kindFlag   string
		nameFlag   string
		wantKind   string
		wantName   string
	}{
		{
			name:     "Kind/Name syntax",
			args:     []string{"Deployment/nginx"},
			wantKind: "Deployment",
			wantName: "nginx",
		},
		{
			name:     "Kind only",
			args:     []string{"Service"},
			wantKind: "Service",
			wantName: "",
		},
		{
			name:     "Wildcard kind (*/name)",
			args:     []string{"*/nginx"},
			wantKind: "",
			wantName: "nginx",
		},
		{
			name:     "Wildcard name (Kind/*)",
			args:     []string{"Service/*"},
			wantKind: "Service",
			wantName: "",
		},
		{
			name:     "Flag overrides positional kind",
			args:     []string{"Service"},
			kindFlag: "Deployment",
			wantKind: "Deployment",
			wantName: "",
		},
		{
			name:     "Flag overrides positional name",
			args:     []string{"Deployment/web"},
			nameFlag: "nginx",
			wantKind: "Deployment",
			wantName: "nginx",
		},
		{
			name:     "Both flags override",
			args:     []string{"ConfigMap/config"},
			kindFlag: "Deployment",
			nameFlag: "nginx",
			wantKind: "Deployment",
			wantName: "nginx",
		},
		{
			name:     "No args, no flags",
			args:     []string{},
			wantKind: "",
			wantName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start with flag values
			name := tt.nameFlag
			kind := tt.kindFlag

			// Parse positional argument (Kind/Name) - mirrors logic in List()
			if len(tt.args) == 1 {
				parts := splitKindName(tt.args[0])
				if kind == "" && parts[0] != "" && parts[0] != "*" {
					kind = parts[0]
				}
				if len(parts) == 2 && name == "" && parts[1] != "*" {
					name = parts[1]
				}
			}

			if kind != tt.wantKind {
				t.Errorf("kind = %q, want %q", kind, tt.wantKind)
			}
			if name != tt.wantName {
				t.Errorf("name = %q, want %q", name, tt.wantName)
			}
		})
	}
}

func TestListWithKindNameSyntax(t *testing.T) {
	input := `apiVersion: v1
kind: Deployment
metadata:
  name: nginx
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
`

	tests := []struct {
		name string
		kind string
		filterName string
		want string
	}{
		{
			name:       "Kind/Name filters both",
			kind:       "Deployment",
			filterName: "nginx",
			want:       "\033[34mDeployment\033[0m/\033[32mnginx\033[0m\n",
		},
		{
			name:       "Kind only filters by kind",
			kind:       "Service",
			filterName: "",
			want:       "\033[34mService\033[0m/\033[32mnginx\033[0m\n",
		},
		{
			name:       "Name only filters by name",
			kind:       "",
			filterName: "nginx",
			want:       "\033[34mDeployment\033[0m/\033[32mnginx\033[0m\n\033[34mService\033[0m/\033[32mnginx\033[0m\n",
		},
		{
			name:       "No filters returns all",
			kind:       "",
			filterName: "",
			want:       "\033[34mConfigMap\033[0m/\033[32mconfig\033[0m\n\033[34mDeployment\033[0m/\033[32mnginx\033[0m\n\033[34mService\033[0m/\033[32mnginx\033[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer

			Find(
				&kio.ByteReader{Reader: bytes.NewReader([]byte(input))},
				BuildFilters(tt.filterName, tt.kind),
				BuildOutputs(true, &output))

			got := output.String()
			if got != tt.want {
				t.Errorf("output = %q, want %q", got, tt.want)
			}
		})
	}
}

// Helper to split Kind/Name for tests
func splitKindName(arg string) []string {
	for i, c := range arg {
		if c == '/' {
			return []string{arg[:i], arg[i+1:]}
		}
	}
	return []string{arg}
}
