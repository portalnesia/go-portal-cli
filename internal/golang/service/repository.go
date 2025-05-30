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

func (s *addRepository) addServiceRepository(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)
	lowerName := strings.ToLower(s.cfg.Name)
	structName := fmt.Sprintf("%sRepository", lowerName)

	_, _ = color.New(color.FgBlue).Printf("Generating repository\n")
	pkgImport := []string{
		fmt.Sprintf(`"%s/internal/model"`, s.cfg.Module),
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
								Type: &ast.IndexListExpr{
									X: &ast.Ident{Name: "CrudRepository"},
									Indices: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("model"),
											Sel: ast.NewIdent(serviceName),
										},
										ast.NewIdent("string"),
										&ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("request"),
												Sel: ast.NewIdent("Request"),
											},
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
								Type: &ast.IndexListExpr{
									X: ast.NewIdent("crudRepository"),
									Indices: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("model"),
											Sel: ast.NewIdent(serviceName),
										},
										ast.NewIdent("string"),
										&ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("request"),
												Sel: ast.NewIdent("Request"),
											},
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

	// implemented List
	decls = append(decls, &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("new%sRepository", serviceName)),
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("bs")},
						Type:  ast.NewIdent("base"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("%sRepository", serviceName)),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CompositeLit{
							Type: ast.NewIdent(structName),
							Elts: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.IndexListExpr{
										X: ast.NewIdent("crudRepository"),
										Indices: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("model"),
												Sel: ast.NewIdent(serviceName),
											},
											ast.NewIdent("string"),
											&ast.StarExpr{
												X: &ast.SelectorExpr{
													X:   ast.NewIdent("request"),
													Sel: ast.NewIdent("Request"),
												},
											},
										},
									},
									Elts: []ast.Expr{
										ast.NewIdent("bs"),
									},
								},
							},
						},
					},
				},
			},
		},
	})

	imports := helper.GenImport(pkgImport...)
	decls = append([]ast.Decl{imports.GenDecl()}, decls...)
	file := &ast.File{
		Name:    ast.NewIdent("repository"),
		Imports: imports.ImportSpec(),
		Decls:   decls,
	}

	res <- config2.Builder{
		File:     file,
		Pathname: fmt.Sprintf("internal/repository/%s_repository.go", lowerName),
	}
}

func (s *addRepository) addToRepository(wg *sync.WaitGroup, res chan<- config2.Builder) {
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
								Value: &ast.CallExpr{
									Fun:  ast.NewIdent(fmt.Sprintf("new%sRepository", serviceName)),
									Args: []ast.Expr{ast.NewIdent("bs")},
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
