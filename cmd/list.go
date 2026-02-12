package cmd

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/rjferguson21/kgrep/printwriter"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
)

func init() {
	listCmd.Flags().BoolP("summary", "s", false, "use summary list")
	listCmd.Flags().StringP("kind", "k", "", "query by kind")
	listCmd.Flags().StringP("name", "n", "", "query by name")
	// rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "kgrep [kind[/name]]",
	Short: "Find kubernetes resources",
	Long: `Search through a list of kubernetes resources for a specific resource, or resources.

Examples:
  # Search for all Deployments
  helm template chart | kgrep Deployment

  # Search for a specific Service by name
  cat manifests.yaml | kgrep Service/nginx

  # Search for any resource named nginx
  kubectl get all -o yaml | kgrep '*/nginx'

  # Use flags for more control
  kubectl get all -o yaml | kgrep --kind Pod --name "nginx.*" --summary
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		List(cmd, args)
	},
}

func Find(input kio.Reader, filters []kio.Filter, output kio.Writer) {
	err := kio.Pipeline{
		Inputs:  []kio.Reader{input},
		Outputs: []kio.Writer{output},
		Filters: filters,
	}.Execute()

	if err != nil {
		log.Fatal(err)
	}
}

func BuildOutputs(summary bool, writer io.Writer) kio.Writer {
	var output kio.Writer

	if summary {
		output = printwriter.PrintWriter{
			Writer: writer,
			Root:   ".",
		}
	} else {
		output = kio.ByteWriter{Writer: writer}
	}

	return output
}

func BuildInputs(cmd *cobra.Command) kio.Reader {
	return &kio.ByteReader{Reader: cmd.InOrStdin()}
}

func BuildFilters(name string, kind string) []kio.Filter {
	grepFilters := []kio.Filter{&filters.GrepFilter{
		Path:  []string{"metadata", "name"},
		Value: name,
	}}

	if kind != "" {
		grepFilters = append(grepFilters, &filters.GrepFilter{
			Path:  []string{"kind"},
			Value: kind,
		})
	}

	return grepFilters
}

func List(cmd *cobra.Command, args []string) {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Fatal(err)
	}

	summary, err := cmd.Flags().GetBool("summary")
	if err != nil {
		log.Fatal(err)
	}

	kind, err := cmd.Flags().GetString("kind")
	if err != nil {
		log.Fatal(err)
	}

	// Parse positional argument (Kind/Name)
	if len(args) == 1 {
		parts := strings.SplitN(args[0], "/", 2)
		// Positional kind (if not overridden by flag), "*" means any
		if kind == "" && parts[0] != "" && parts[0] != "*" {
			kind = parts[0]
		}
		// Positional name (if present and not overridden by flag), "*" means any
		if len(parts) == 2 && name == "" && parts[1] != "*" {
			name = parts[1]
		}
	}

	Find(
		BuildInputs(cmd),
		BuildFilters(name, kind),
		BuildOutputs(summary, os.Stdout))
}
