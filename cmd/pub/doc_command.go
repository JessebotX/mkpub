package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

type DocCommand struct {
	InputPath string `arg:"" type:"existingfile"`
}

// TODO: still some work left here...
func (d DocCommand) Run(ctx *Context) error {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, d.InputPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	p, err := doc.NewFromFiles(fset, []*ast.File{f}, "./")
	if err != nil {
		return err
	}

	for _, t := range p.Types {
		if t.Name != "Book" {
			continue
		}

		for _, spec := range t.Decl.Specs {
			switch spec.(type) {
			case *ast.TypeSpec:
				typeSpec := spec.(*ast.TypeSpec)

				fmt.Printf("%s\n", typeSpec.Name.Name)

				switch typeSpec.Type.(type) {
				case *ast.StructType:
					structType := typeSpec.Type.(*ast.StructType)
					for _, field := range structType.Fields.List {
						for _, name := range field.Names {
							fmt.Printf("--- Field: %s ---\n", name.Name)
							if field.Doc != nil {
								fmt.Printf("%s\n", strings.TrimSpace(field.Doc.Text()))
							}
						}
					}
				}
			}
		}
	}

	return nil
}
