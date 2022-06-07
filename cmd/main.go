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
var protoPath string

func main() {
	flag.Var(&files, "input", "Input file list.")
	flag.StringVar(
		&protoPath, "proto_path", "",
		"Directory relative to which all .proto files are found.",
	)
	flag.StringVar(&outDir, "out", "", "Output directory.")
	flag.Parse()

	if len(files) == 0 {
		fmt.Println("Use --input option to specify input files.")
		os.Exit(-1)
	}

	if err := generator.Generate(protoPath, files, outDir); err != nil {
		fmt.Println(err.Error())
		os.Exit(-2)
	}
}
