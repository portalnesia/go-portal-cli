/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package service

import (
	"fmt"
	"go/token"
	"strings"
	"sync"

	"github.com/dave/dst"
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
	"go.portalnesia.com/utils"
)

func (s *addRepository) addRepositoryModel(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	_, _ = color.New(color.FgBlue).Printf("Generating model\n")
	serviceName := strings.ReplaceAll(utils.Ucwords(strings.ReplaceAll(s.cfg.Name, "_", " ")), " ", "")

	pkgImport := []string{
		`"github.com/uptrace/bun"`,
	}

	decls := make([]dst.Decl, 0)

	// type User struct{}
	decls = append(decls, &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(serviceName),
				Type: &dst.StructType{
					Fields: &dst.FieldList{
						List: []*dst.Field{
							{
								// Field ini tidak punya Name karena merupakan Embedding (Anonymous Field)
								Type: &dst.SelectorExpr{
									X:   dst.NewIdent("bun"),
									Sel: dst.NewIdent("BaseModel"),
								},
								// Tag harus dibungkus dengan backtick di dalam string literal
								Tag: &dst.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf("`bun:\"table:%s\"`", strings.ToLower(s.cfg.Name)),
								},
							},
						},
					},
				},
			},
		},
	})

	// func (User) GetDefaultOrder() []string { return []string{"created_at", "desc"} }
	decls = append(decls, &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Type: dst.NewIdent(serviceName),
				},
			},
		},
		Name: dst.NewIdent("GetDefaultOrder"),
		Type: &dst.FuncType{
			Results: &dst.FieldList{
				List: []*dst.Field{
					{
						Type: &dst.ArrayType{
							Len: &dst.BasicLit{Kind: token.INT, Value: "2"}, Elt: dst.NewIdent("string"),
						},
					},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.CompositeLit{
							Type: &dst.ArrayType{
								Len: &dst.BasicLit{Kind: token.INT, Value: "2"}, Elt: dst.NewIdent("string"),
							},
							Elts: []dst.Expr{
								&dst.BasicLit{Kind: token.STRING, Value: `"created_at"`},
								&dst.BasicLit{Kind: token.STRING, Value: `"desc"`},
							},
						},
					},
				},
			},
		},
	})

	// func (User) GetAvailableOrder() [][]string { return [][]string{{"created_at", "desc"}} }
	decls = append(decls, &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Type: dst.NewIdent(serviceName),
				},
			},
		},
		Name: dst.NewIdent("GetAvailableOrder"),
		Type: &dst.FuncType{
			Results: &dst.FieldList{
				List: []*dst.Field{{Type: &dst.ArrayType{Elt: dst.NewIdent("string")}}},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.CompositeLit{
							Type: &dst.ArrayType{Elt: dst.NewIdent("string")},
							Elts: []dst.Expr{
								&dst.BasicLit{Kind: token.STRING, Value: `"created_at"`},
							},
						},
					},
				},
			},
		},
	})

	for i := range decls {
		decls[i].Decorations().Before = dst.EmptyLine
	}

	imports := helper.GenImport(pkgImport...)
	decls = append([]dst.Decl{imports.GenDeclDst()}, decls...)
	file := &dst.File{
		Name:  dst.NewIdent("model"),
		Decls: decls,
	}

	res <- config2.Builder{
		DstFile:  file,
		Pathname: fmt.Sprintf("internal/model/%s.go", strings.ToLower(s.cfg.Name)),
	}
}
