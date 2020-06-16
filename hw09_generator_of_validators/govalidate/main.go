package main

import (
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Filename was not found")
	}
	fileName := os.Args[1]
	input, _ := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
	buf, err := parseAstFile(input)
	if err != nil {
		log.Fatal(err)
	}
	// write in writer end of pipe
	code, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal("Cannot prettify code", err)
	}
	err = ioutil.WriteFile(strings.Replace(fileName, ".go", "_validation_generated.go", 1), code, os.ModePerm)
	if err != nil {
		log.Fatal("Cannot write formatted file", err)
	}
}
