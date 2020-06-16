// +build generation

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	code = `
	package models

	type User struct {
		ID     string %s
		Name   string
		Age    int      %s
		Email  string   %s
		Role   UserRole %s
		Phones []string %s
	}
	type Response struct {
		Code int    %s
		Body string
	}
	`
)

func TestCodeParsing(t *testing.T) {
	t.Run("should parse package name correctly", func(t *testing.T) {
		f, _ := parser.ParseFile(token.NewFileSet(), "", unitParseTemplate(code), parser.AllErrors)
		metadata := prepareMetadata(f)
		require.Equal(t, "models", unitStringify(metadata.Package))
	})

	t.Run("should parse model without 'validate' tags correctly", func(t *testing.T) {
		f, _ := parser.ParseFile(token.NewFileSet(), "", unitParseTemplate(code), parser.AllErrors)
		metadata := prepareMetadata(f)

		require.Equal(t, 0, len(metadata.Models))
	})

	t.Run("should parse model correctly", func(t *testing.T) {
		f, _ := parser.ParseFile(token.NewFileSet(), "", unitParseTemplate(code, "`validate:\"len:36\"`"), parser.AllErrors)
		metadata := prepareMetadata(f)

		require.Equal(t, 1, len(metadata.Models))
		model := metadata.Models[0]
		require.Equal(t, "User", unitStringify(model.Name))
	})

	t.Run("should parse tag correctly", func(t *testing.T) {
		f, _ := parser.ParseFile(token.NewFileSet(), "", unitParseTemplate(code, "`validate:\"len:36\"`"), parser.AllErrors)
		metadata := prepareMetadata(f)

		model := metadata.Models[0]
		require.Equal(t, 1, len(model.Declarations))
		require.Equal(t, "ID", unitStringify(model.Declarations[0].Field))
		require.Equal(t, "36", model.Declarations[0].Param)
		require.Equal(t, "len", model.Declarations[0].Type)
		require.Equal(t, false, model.Declarations[0].IsArrayField)
	})

	t.Run("should parse multiple tags correctly", func(t *testing.T) {
		f, _ := parser.ParseFile(token.NewFileSet(), "", unitParseTemplate(code, "", "`validate:\"min:18|max:50\"`"), parser.AllErrors)
		metadata := prepareMetadata(f)

		model := metadata.Models[0]
		require.Equal(t, 2, len(model.Declarations))
		require.Equal(t, "Age", unitStringify(model.Declarations[0].Field))
		require.Equal(t, "18", model.Declarations[0].Param)
		require.Equal(t, "min", model.Declarations[0].Type)
		require.Equal(t, false, model.Declarations[0].IsArrayField)

		require.Equal(t, "Age", unitStringify(model.Declarations[1].Field))
		require.Equal(t, "50", model.Declarations[1].Param)
		require.Equal(t, "max", model.Declarations[1].Type)
		require.Equal(t, false, model.Declarations[1].IsArrayField)
	})

	t.Run("should parse array field correctly", func(t *testing.T) {
		f, _ := parser.ParseFile(token.NewFileSet(), "", unitParseTemplate(code, "", "", "", "", "`validate:\"len:11\"`"), parser.AllErrors)
		metadata := prepareMetadata(f)

		model := metadata.Models[0]
		require.Equal(t, 1, len(model.Declarations))
		require.Equal(t, "Phones", unitStringify(model.Declarations[0].Field))
		require.Equal(t, "11", model.Declarations[0].Param)
		require.Equal(t, "len", model.Declarations[0].Type)
		require.Equal(t, true, model.Declarations[0].IsArrayField)
	})

	t.Run("should recognize tags, which are depended on type correctly", func(t *testing.T) {
		f, _ := parser.ParseFile(token.NewFileSet(), "", unitParseTemplate(code, "", "`validate:\"in:200,404,500\"`", "`validate:\"in:admin,stuff\"`"), parser.AllErrors)
		metadata := prepareMetadata(f)

		model := metadata.Models[0]
		require.Equal(t, 2, len(model.Declarations))
		require.Equal(t, "Age", unitStringify(model.Declarations[0].Field))
		require.Equal(t, "200,404,500", model.Declarations[0].Param)
		require.Equal(t, "inInt", model.Declarations[0].Type)
		require.Equal(t, false, model.Declarations[0].IsArrayField)
		require.Equal(t, "Email", unitStringify(model.Declarations[1].Field))
		require.Equal(t, "admin,stuff", model.Declarations[1].Param)
		require.Equal(t, "inString", model.Declarations[1].Type)
		require.Equal(t, false, model.Declarations[1].IsArrayField)
	})

	t.Run("should parse multiple models correctly", func(t *testing.T) {
		f, _ := parser.ParseFile(token.NewFileSet(), "", unitParseTemplate(code, "`validate:\"len:36\"`", "", "", "", "", "`validate:\"in:200,404,500\"`"), parser.AllErrors)
		metadata := prepareMetadata(f)

		require.Equal(t, 2, len(metadata.Models))
		require.Equal(t, 1, len(metadata.Models[0].Declarations))
		require.Equal(t, "ID", unitStringify(metadata.Models[0].Declarations[0].Field))
		require.Equal(t, "36", metadata.Models[0].Declarations[0].Param)
		require.Equal(t, "len", metadata.Models[0].Declarations[0].Type)
		require.Equal(t, false, metadata.Models[0].Declarations[0].IsArrayField)

		require.Equal(t, 1, len(metadata.Models[1].Declarations))
		require.Equal(t, "Code", unitStringify(metadata.Models[1].Declarations[0].Field))
		require.Equal(t, "200,404,500", metadata.Models[1].Declarations[0].Param)
		require.Equal(t, "inInt", metadata.Models[1].Declarations[0].Type)
		require.Equal(t, false, metadata.Models[1].Declarations[0].IsArrayField)
	})
}

func unitParseTemplate(code string, args ...interface{}) string {
	stringsInCode := strings.Count(code, "%s")
	rest := make([]interface{}, stringsInCode-len(args))
	for index, _ := range rest {
		rest[index] = ""
	}
	params := append(args, rest...)
	return fmt.Sprintf(code, params...)
}

func unitStringify(node ast.Node) string {
	var typeNameBuf bytes.Buffer
	fset := token.NewFileSet()
	printer.Fprint(&typeNameBuf, fset, node)
	return typeNameBuf.String()
}
