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
	"sync"
)

func (s *addService) addServiceRoutes(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)
	group := fmt.Sprintf(`"/%s/%s"`, s.cfg.Version, s.cfg.Name)
	if s.cfg.Version == "" {
		group = fmt.Sprintf(`"/%s"`, s.cfg.Name)
	}

	_, _ = color.New(color.FgBlue).Printf("Generating routes\n")
	pkgImport := []string{
		fmt.Sprintf(`"%s/internal/rest/handler"`, s.cfg.Module),
	}
	decls := make([]ast.Decl, 0)

	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("Routes")},
				},
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("routes%s", serviceName)),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// h := handler.NewTest(r.app)
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("h")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("handler"),
								Sel: ast.NewIdent(fmt.Sprintf("New%s", serviceName)),
							},
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("r"),
									Sel: ast.NewIdent("app"),
								},
							},
						},
					},
				},
				// route := r.fiber.Group("/v1/test")
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("route")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("r"),
									Sel: ast.NewIdent("fiber"),
								},
								Sel: ast.NewIdent("Group"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(group),
								},
							},
						},
					},
				},
				helper.BodyListNewLines(),
				// List
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("route"),
							Sel: ast.NewIdent("Get"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"/\"",
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("h"),
								Sel: ast.NewIdent(fmt.Sprintf("List%s", serviceName)),
							},
						},
					},
				},
				// Get
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("route"),
							Sel: ast.NewIdent("Get"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"/:id\"",
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("h"),
								Sel: ast.NewIdent(fmt.Sprintf("Get%s", serviceName)),
							},
						},
					},
				},
				// Post
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("route"),
							Sel: ast.NewIdent("Post"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"/\"",
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("h"),
								Sel: ast.NewIdent(fmt.Sprintf("Create%s", serviceName)),
							},
						},
					},
				},
				// Put
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("route"),
							Sel: ast.NewIdent("Put"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"/:id\"",
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("h"),
								Sel: ast.NewIdent(fmt.Sprintf("Update%s", serviceName)),
							},
						},
					},
				},
				// Patch
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("route"),
							Sel: ast.NewIdent("Patch"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"/:id\"",
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("h"),
								Sel: ast.NewIdent(fmt.Sprintf("Update%s", serviceName)),
							},
						},
					},
				},
				// Delete
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("route"),
							Sel: ast.NewIdent("Delete"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"/:id\"",
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("h"),
								Sel: ast.NewIdent(fmt.Sprintf("Delete%s", serviceName)),
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
		Name:    ast.NewIdent("routes"),
		Imports: imports.ImportSpec(),
		Decls:   decls,
	}

	res <- config2.Builder{
		File:     file,
		Pathname: fmt.Sprintf("internal/rest/routes/%s_route.go", s.cfg.Name),
	}
}

func (s *addService) addToRoutes(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)

	// Parse file routes
	resp := config2.Builder{}
	defer func() {
		res <- resp
	}()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, s.app.Dir("internal/rest/routes/routes.go"), nil, parser.AllErrors)
	if err != nil {
		resp.Err = err
		return
	}
	var found bool

	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name.Name == "initRoutes" {
				_, _ = color.New(color.FgBlue).Printf("Add to routes\n")
				found = true
				// Tambahkan routes
				stmt := &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent(fmt.Sprintf("routes%s", serviceName)),
						},
					},
				}

				// Append ke Body.List
				funcDecl.Body.List = append(funcDecl.Body.List, stmt)
				resp.File = file
				resp.Pathname = "internal/rest/routes/routes.go"
				break
			}
		}
	}
	if !found {
		resp.Err = fmt.Errorf("routes not found")
		return
	}
}

func (s *addEndpoint) addEndpointRoutes(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	_, _ = color.New(color.FgBlue).Printf("Generating routes\n")

	serviceName := utils.Ucwords(s.cfg.ServiceName)

	// Parse file routes
	resp := config2.Builder{}
	defer func() {
		res <- resp
	}()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, s.app.Dir(fmt.Sprintf("internal/rest/routes/%s_route.go", s.cfg.ServiceName)), nil, parser.AllErrors)
	if err != nil {
		resp.Err = err
		return
	}

	decls := file.Decls[len(file.Decls)-1]
	funcDecl, ok := decls.(*ast.FuncDecl)
	if !ok {
		resp.Err = fmt.Errorf("invalid routes: function routes%s not found", serviceName)
		return
	}
	funcDecl.Body.List = append(funcDecl.Body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("route"),
				Sel: ast.NewIdent(s.cfg.Method),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, s.cfg.Path),
				},
				&ast.SelectorExpr{
					X:   ast.NewIdent("h"),
					Sel: ast.NewIdent(s.cfg.Name),
				},
			},
		},
	})
	if s.cfg.Method == "Put" {
		funcDecl.Body.List = append(funcDecl.Body.List, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("route"),
					Sel: ast.NewIdent("Patch"),
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s"`, s.cfg.Path),
					},
					&ast.SelectorExpr{
						X:   ast.NewIdent("h"),
						Sel: ast.NewIdent(s.cfg.Name),
					},
				},
			},
		})
	}

	file.Decls[len(file.Decls)-1] = funcDecl

	resp = config2.Builder{
		File:     file,
		Pathname: fmt.Sprintf("internal/rest/routes/%s_route.go", s.cfg.ServiceName),
	}
}
