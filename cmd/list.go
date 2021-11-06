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
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		List(cmd, args)
	},
}

func List(cmd *cobra.Command, args []string) {
	summary, err := cmd.Flags().GetBool("summary")
	query := args[0]

	input := kio.LocalPackageReader{
		PackagePath:       ".",
		MatchFilesGlob:    []string{"*.yaml"},
		PreserveSeqIndent: true,
		WrapBareSeqNode:   true,
	}

	fltrs := []kio.Filter{&filters.GrepFilter{
		Path:  []string{"metadata", "name"},
		Value: query,
	}}

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
