package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tigrannajaryan/exp-lazyproto/generator"
)

type flagList []string

func (f *flagList) String() string {
	return strings.Join(*f, " ")
}

func (f *flagList) Set(value string) error {
	*f = append(*f, value)
	return nil
}

var files flagList
var outDir string
var includePaths flagList

func main() {
	options := generator.Options{}

	//flag.Var(&files, "input", "Input file list.")
	flag.Var(
		&includePaths, "proto_path",
		"Directory relative to which all .proto files are found.",
	)
	flag.StringVar(&outDir, "go_out", "", "Output directory.")
	flag.BoolVar(
		&options.WithPresence, "with_presence", false, "Generate presence methods.",
	)
	flag.Parse()

	files := flag.Args()

	if len(files) == 0 {
		fmt.Println("Use --input option to specify input files.")
		os.Exit(-1)
	}

	if err := generator.Generate(includePaths, files, outDir, options); err != nil {
		fmt.Println(err.Error())
		os.Exit(-2)
	}
}
