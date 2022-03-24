package cmd

import (
	"io"
	"log"
	"os"

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
	Use:   "kgrep [resources.yaml]",
	Short: "Find kubernetes resources",
	Long: `Search through a list of kubernetes resources for a specific resource, or resources. For example:

# Search for services within all.yaml
kgrep --kind Service all.yaml

# Search for a Deployment named foo within helm chart
helm template chart | kgrep --kind Deployment --name foo
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

func BuildInputs(cmd *cobra.Command, args []string) kio.Reader {
	var input kio.Reader

	if len(args) == 0 {
		input = &kio.ByteReader{Reader: cmd.InOrStdin()}
	} else {
		input = kio.LocalPackageReader{
			PackagePath:       args[0],
			MatchFilesGlob:    []string{"*.yaml"},
			PreserveSeqIndent: true,
			WrapBareSeqNode:   true,
		}
	}
	return input
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

	Find(
		BuildInputs(cmd, args),
		BuildFilters(name, kind),
		BuildOutputs(summary, os.Stdout))
}
