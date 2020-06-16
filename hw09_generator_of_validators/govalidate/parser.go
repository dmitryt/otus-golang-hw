package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"regexp"
	"strings"
	"text/template"
)

var validateRegexp = regexp.MustCompile(`validate:"(.*)?"`)
var typeDependentValidations = map[string]bool{
	"in": true,
}

type declaration struct {
	Method       string
	ModelName    *ast.Ident
	Field        *ast.Ident
	Param        string
	FieldType    ast.Expr
	Type         string
	IsArrayField bool
}

type model struct {
	Name         *ast.Ident
	Declarations []declaration
}

type metadata struct {
	Package *ast.Ident
	Models  []model
}

func stringify(node ast.Node) string {
	var typeNameBuf bytes.Buffer
	fset := token.NewFileSet()
	printer.Fprint(&typeNameBuf, fset, node)
	return typeNameBuf.String()
}

func adjustValidationType(validationType string, fieldType ast.Expr) string {
	result := validationType
	if typeDependentValidations[validationType] {
		suffix := stringify(fieldType)
		if suffix != "int" {
			// use string type by default
			suffix = "string"
		}
		result = fmt.Sprintf("%s%s", validationType, strings.Title(suffix))
	}
	return result
}

func prepareDeclarations(fields []*ast.Field, modelName *ast.Ident) []declaration {
	var declarations []declaration
	for _, field := range fields {
		if field.Tag != nil {
			matches := validateRegexp.FindAllStringSubmatch(field.Tag.Value, -1)
			if len(matches) > 0 && len(matches[0]) > 1 {
				for _, expr := range strings.Split(matches[0][1], "|") {
					params := strings.Split(expr, ":")
					vType := adjustValidationType(params[0], field.Type)
					_, isArrayField := field.Type.(*ast.ArrayType)
					declarations = append(declarations, declaration{
						Method:       fmt.Sprintf("_validate%s%s", strings.Title(vType), strings.Title(stringify(field.Names[0]))),
						ModelName:    modelName,
						Field:        field.Names[0],
						Type:         vType,
						FieldType:    field.Type,
						Param:        params[1],
						IsArrayField: isArrayField,
					})
				}
			}
		}
	}
	return declarations
}

func prepareModels(input *ast.File) []model {
	var models []model
	for _, d := range input.Decls {
		decl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range decl.Specs {
			s, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			spec, ok := s.Type.(*ast.StructType)
			if !ok {
				continue
			}
			declarations := prepareDeclarations(spec.Fields.List, s.Name)
			if len(declarations) > 0 {
				models = append(models, model{
					Name:         s.Name,
					Declarations: declarations,
				})
			}
		}
	}
	return models
}

func prepareMetadata(file *ast.File) metadata {
	return metadata{
		Package: file.Name,
		Models:  prepareModels(file),
	}
}

func parseAstFile(file *ast.File) (*bytes.Buffer, error) {
	tmpl, err := template.New("validation").Parse(baseTemplate)
	if err != nil {
		return nil, errors.New("cannot parse template correctly")
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, prepareMetadata(file)); err != nil {
		return nil, errors.New("cannot generate template correctly")
	}
	return &buf, nil
}
