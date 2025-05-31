/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package service

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
	"go.portalnesia.com/utils"
	"go/parser"
	"go/token"
	"strings"
	"sync"
)

func (s *addRepository) addServiceRepository(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := strings.ReplaceAll(utils.Ucwords(strings.ReplaceAll(s.cfg.Name, "_", " ")), " ", "")
	lowerName := strings.ReplaceAll(helper.FirstToLower(strings.ReplaceAll(s.cfg.Name, "_", " ")), " ", "")
	structName := fmt.Sprintf("%sRepository", lowerName)

	_, _ = color.New(color.FgBlue).Printf("Generating repository\n")
	pkgImport := []string{
		fmt.Sprintf(`"%s/internal/model"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/request"`, s.cfg.Module),
	}

	decls := make([]dst.Decl, 0)

	// interface
	decls = append(decls, &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(fmt.Sprintf("%sRepository", serviceName)),
				Type: &dst.InterfaceType{
					Methods: &dst.FieldList{
						List: []*dst.Field{
							{
								Type: &dst.IndexListExpr{
									X: &dst.Ident{Name: "CrudRepository"},
									Indices: []dst.Expr{
										&dst.SelectorExpr{
											X:   dst.NewIdent("model"),
											Sel: dst.NewIdent(serviceName),
										},
										dst.NewIdent("string"),
										&dst.StarExpr{
											X: &dst.SelectorExpr{
												X:   dst.NewIdent("request"),
												Sel: dst.NewIdent("Request"),
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
	decls = append(decls, &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(structName),
				Type: &dst.StructType{
					Fields: &dst.FieldList{
						List: []*dst.Field{
							{
								Type: &dst.IndexListExpr{
									X: dst.NewIdent("crudRepository"),
									Indices: []dst.Expr{
										&dst.SelectorExpr{
											X:   dst.NewIdent("model"),
											Sel: dst.NewIdent(serviceName),
										},
										dst.NewIdent("string"),
										&dst.StarExpr{
											X: &dst.SelectorExpr{
												X:   dst.NewIdent("request"),
												Sel: dst.NewIdent("Request"),
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
	decls = append(decls, &dst.FuncDecl{
		Name: dst.NewIdent(fmt.Sprintf("new%sRepository", serviceName)),
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("bs")},
						Type:  dst.NewIdent("base"),
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{
						Type: dst.NewIdent(fmt.Sprintf("%sRepository", serviceName)),
					},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.CompositeLit{
							Type: dst.NewIdent(structName),
							Elts: []dst.Expr{
								&dst.CompositeLit{
									Type: &dst.IndexListExpr{
										X: dst.NewIdent("crudRepository"),
										Indices: []dst.Expr{
											&dst.SelectorExpr{
												X:   dst.NewIdent("model"),
												Sel: dst.NewIdent(serviceName),
											},
											dst.NewIdent("string"),
											&dst.StarExpr{
												X: &dst.SelectorExpr{
													X:   dst.NewIdent("request"),
													Sel: dst.NewIdent("Request"),
												},
											},
										},
									},
									Elts: []dst.Expr{
										dst.NewIdent("bs"),
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
	decls = append([]dst.Decl{imports.GenDeclDst()}, decls...)
	for i := range decls {
		decls[i].Decorations().Before = dst.EmptyLine
	}
	file := &dst.File{
		Name:    dst.NewIdent("repository"),
		Imports: imports.ImportSpecDst(),
		Decls:   decls,
	}

	res <- config2.Builder{
		DstFile:  file,
		Pathname: fmt.Sprintf("internal/repository/%s_repository.go", strings.ToLower(s.cfg.Name)),
	}
}

func (s *addRepository) addToRepository(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := strings.ReplaceAll(utils.Ucwords(strings.ReplaceAll(s.cfg.Name, "_", " ")), " ", "")

	// Parse file routes
	resp := config2.Builder{}
	defer func() {
		res <- resp
	}()

	fset := token.NewFileSet()
	file, err := decorator.ParseFile(fset, s.app.Dir("internal/repository/repository.go"), nil, parser.AllErrors)
	if err != nil {
		resp.Err = err
		return
	}

	for _, decl := range file.Decls {
		// struct Registry
		if interfaceDecl, okGenDecl := decl.(*dst.GenDecl); okGenDecl {
			for specs, _ := range interfaceDecl.Specs {
				if typeSpecs, okTypeSpec := interfaceDecl.Specs[specs].(*dst.TypeSpec); okTypeSpec {
					structType, okStruct := typeSpecs.Type.(*dst.StructType)
					if typeSpecs.Name.Name == "Registry" && okStruct {
						structType.Fields.List = append(structType.Fields.List, &dst.Field{
							Names: []*dst.Ident{dst.NewIdent(serviceName)},
							Type:  dst.NewIdent(fmt.Sprintf("%sRepository", serviceName)),
						})
					}
				}
			}
		}

		if funcDecl, okFunc := decl.(*dst.FuncDecl); okFunc {
			if funcDecl.Name.Name == "NewRepository" {
				for _, stmt := range funcDecl.Body.List {
					if returnStmt, okStmt := stmt.(*dst.ReturnStmt); okStmt {
						if compositLit, okCompositLit := returnStmt.Results[0].(*dst.CompositeLit); okCompositLit {
							compositLit.Elts = append(compositLit.Elts, &dst.KeyValueExpr{
								Key: dst.NewIdent(serviceName),
								Value: &dst.CallExpr{
									Fun:  dst.NewIdent(fmt.Sprintf("new%sRepository", serviceName)),
									Args: []dst.Expr{dst.NewIdent("bs")},
								},
							})
						}
					}
				}
			}
		}
	}

	resp.DstFile = file
	resp.Pathname = "internal/repository/repository.go"
}
