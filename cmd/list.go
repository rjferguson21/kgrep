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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		List(cmd, args)
	},
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

	path := args[0]

	input := kio.LocalPackageReader{
		PackagePath:       path,
		MatchFilesGlob:    []string{"*.yaml"},
		PreserveSeqIndent: true,
		WrapBareSeqNode:   true,
	}

	fltrs := []kio.Filter{&filters.GrepFilter{
		Path:  []string{"metadata", "name"},
		Value: name,
	}}

	if kind != "" {
		fltrs = append(fltrs, &filters.GrepFilter{
			Path:  []string{"kind"},
			Value: kind,
		})
	}

	var outputs []kio.Writer

	if summary {
		outputs = append(outputs, printwriter.PrintWriter{
			Writer: os.Stdout,
			Root:   ".",
		})
	} else {
		outputs = append(outputs, kio.ByteWriter{Writer: os.Stdout})
	}

	err = kio.Pipeline{
		Inputs:  []kio.Reader{input},
		Outputs: outputs,
		Filters: fltrs,
	}.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
