package printwriter

import (
	"fmt"
	"io"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type PrintWriter struct {
	Writer         io.Writer
	Root           string
	FunctionConfig *yaml.RNode
}

func (p PrintWriter) Write(nodes []*yaml.RNode) error {
	for _, v := range nodes {
		fmt.Println(v.GetKind(), v.GetName())
	}
	return nil
}
