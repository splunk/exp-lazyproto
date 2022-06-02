package generator

import "github.com/yoheimuta/go-protoparser/v4/parser"

type File struct {
	InputFilePath string
	PackageName   string
	Messages      map[string]*Message
}

type Message struct {
	Name   string
	Fields map[string]*parser.Field
}
