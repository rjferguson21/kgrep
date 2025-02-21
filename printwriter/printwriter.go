package printwriter

import (
	"fmt"
	"io"
	"sort"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// ANSI color codes
const (
	blue   = "\033[34m"
	green  = "\033[32m"
	reset  = "\033[0m"
)

type PrintWriter struct {
	Writer         io.Writer
	Root           string
	FunctionConfig *yaml.RNode
}

// resourceInfo holds the kind and name of a resource for sorting
type resourceInfo struct {
	kind string
	name string
}

func (p PrintWriter) Write(nodes []*yaml.RNode) error {
	// Convert nodes to resourceInfo for sorting
	resources := make([]resourceInfo, len(nodes))
	for i, v := range nodes {
		resources[i] = resourceInfo{
			kind: v.GetKind(),
			name: v.GetName(),
		}
	}

	// Sort resources by kind
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].kind < resources[j].kind
	})

	// Print sorted resources
	for _, r := range resources {
		fmt.Fprintf(p.Writer, "%s%s%s/%s%s%s\n",
			blue, r.kind, reset,
			green, r.name, reset)
	}
	return nil
}
