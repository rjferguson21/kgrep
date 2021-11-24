package prettywriter

import (
	"container/list"
	"io"
	"log"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

type PrintWriter struct {
	Writer         io.Writer
	Root           string
	FunctionConfig *kyaml.RNode
}

func resultsToString(results *list.List) []string {
	var pretty []string = make([]string, 0)
	for el := results.Front(); el != nil; el = el.Next() {
		n := el.Value.(*yqlib.CandidateNode)
		pretty = append(pretty, yqlib.NodeToString(n))
	}
	return pretty
}

func (p PrintWriter) Write(nodes []*kyaml.RNode) error {
	level, _ := logging.LogLevel("WARN")
	logging.SetLevel(level, "")

	format, err := yqlib.OutputFormatFromString("yaml")
	if err != nil {
		log.Fatal(err)
	}

	var yamlNodes []*yaml.Node
	for _, v := range nodes {
		var n yaml.Node
		b := []byte(v.MustString())
		err := yaml.Unmarshal(b, &n)
		if err != nil {
			log.Fatal(err)
		}
		yamlNodes = append(yamlNodes, &n)
	}

	printerWriter := yqlib.NewSinglePrinterWriter(os.Stdout)
	printer := yqlib.NewPrinter(printerWriter, format, false, true, 2, false)

	allAtOnceEvaluator := yqlib.NewAllAtOnceEvaluator()
	list, err := allAtOnceEvaluator.EvaluateNodes("", yamlNodes...)

	printer.PrintResults(list)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
