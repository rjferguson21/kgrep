package cmd

import (
	"log"
	"os"

	"github.com/rjferguson21/kpretty/printwriter"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
)

func init() {
	listCmd.Flags().Bool("summary", false, "use summary list")
	listCmd.Flags().String("kind", "", "query by kind")
	listCmd.Flags().String("name", "", "query by name")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Args:  cobra.MaximumNArgs(1),
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

	var output kio.Writer

	if summary {
		output = printwriter.PrintWriter{
			Writer: os.Stdout,
			Root:   ".",
		}
	} else {
		output = kio.ByteWriter{Writer: os.Stdout}
	}

	Find(input, BuildFilters(name, kind), output)
}
