/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package service

import (
	"fmt"
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
	"go.portalnesia.com/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"sync"
)

func (s *addService) addServiceRepository(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)
	structName := fmt.Sprintf("%sRepository", strings.ToLower(s.cfg.Name))

	_, _ = color.New(color.FgBlue).Printf("Generating repository\n")
	pkgImport := []string{
		`"errors"`,
		fmt.Sprintf(`"%s/internal/context"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/request"`, s.cfg.Module),
	}

	decls := make([]ast.Decl, 0)

	// interface
	decls = append(decls, &ast.GenDecl{
		Tok: token.TYPE,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(fmt.Sprintf("%sRepository", serviceName)),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("List")},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{ast.NewIdent("ctx")},
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("context"),
														Sel: ast.NewIdent("Context"),
													},
												},
											},
											{
												Names: []*ast.Ident{ast.NewIdent("query")},
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("request"),
														Sel: ast.NewIdent("Request"),
													},
												},
											},
											{
												Names: []*ast.Ident{ast.NewIdent("options")},
												Type: &ast.Ellipsis{
													Elt: ast.NewIdent("Options"),
												},
											},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{Type: ast.NewIdent("any")},
											{Type: ast.NewIdent("error")},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	// type struct
	decls = append(decls, &ast.GenDecl{
		Tok: token.TYPE,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(structName),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								// Embedded field: hanya Type, tanpa Names
								Type: ast.NewIdent("base"),
							},
						},
					},
				},
			},
		},
	})

	// implemented List
	decls = append(decls, addImplementation(structName, fmt.Sprintf("List%s", serviceName)))

	imports := helper.GenImport(pkgImport...)
	decls = append([]ast.Decl{imports.GenDecl()}, decls...)
	file := &ast.File{
		Name:    ast.NewIdent("repository"),
		Imports: imports.ImportSpec(),
		Decls:   decls,
	}

	res <- config2.Builder{
		File:     file,
		Pathname: fmt.Sprintf("internal/repository/%s_repository.go", strings.ToLower(s.cfg.Name)),
	}
}

func addImplementation(structName, fnName string) ast.Decl {
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  ast.NewIdent(structName),
				},
			},
		},
		Name: ast.NewIdent("List"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("query")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("request"),
								Sel: ast.NewIdent("Request"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("options")},
						Type: &ast.Ellipsis{
							Elt: ast.NewIdent("Options"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: ast.NewIdent("any")},
					{Type: ast.NewIdent("error")},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.BasicLit{
						Kind:  token.STRING,
						Value: `// db := r.getDatabase(ctx, options...)`,
					},
				},
				helper.BodyListNewLines(),
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("errors"),
								Sel: ast.NewIdent("New"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"not implemented"`,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s *addService) addToRepository(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)

	// Parse file routes
	resp := config2.Builder{}
	defer func() {
		res <- resp
	}()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, s.app.Dir("internal/repository/repository.go"), nil, parser.AllErrors)
	if err != nil {
		resp.Err = err
		return
	}

	for _, decl := range file.Decls {
		// struct Registry
		if interfaceDecl, okGenDecl := decl.(*ast.GenDecl); okGenDecl {
			for specs, _ := range interfaceDecl.Specs {
				if typeSpecs, okTypeSpec := interfaceDecl.Specs[specs].(*ast.TypeSpec); okTypeSpec {
					structType, okStruct := typeSpecs.Type.(*ast.StructType)
					if typeSpecs.Name.Name == "Registry" && okStruct {
						structType.Fields.List = append(structType.Fields.List, &ast.Field{
							Names: []*ast.Ident{ast.NewIdent(serviceName)},
							Type:  ast.NewIdent(fmt.Sprintf("%sRepository", serviceName)),
						})
					}
				}
			}
		}

		if funcDecl, okFunc := decl.(*ast.FuncDecl); okFunc {
			if funcDecl.Name.Name == "NewRepository" {
				for _, stmt := range funcDecl.Body.List {
					if returnStmt, okStmt := stmt.(*ast.ReturnStmt); okStmt {
						if compositLit, okCompositLit := returnStmt.Results[0].(*ast.CompositeLit); okCompositLit {
							compositLit.Elts = append(compositLit.Elts, &ast.KeyValueExpr{
								Key: ast.NewIdent(serviceName),
								Value: &ast.CompositeLit{
									Type: ast.NewIdent(fmt.Sprintf("%sRepository", s.cfg.Name)),
									Elts: []ast.Expr{
										ast.NewIdent("bs"),
									},
								},
							})
						}
					}
				}
			}
		}
	}

	resp.File = file
	resp.Pathname = "internal/repository/repository.go"
}
