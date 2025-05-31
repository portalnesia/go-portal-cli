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
	decls := make([]dst.Decl, 0)

	decls = append(decls, &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent("r")},
					Type:  &dst.StarExpr{X: dst.NewIdent("Routes")},
				},
			},
		},
		Name: dst.NewIdent(fmt.Sprintf("routes%s", serviceName)),
		Type: &dst.FuncType{
			Params:  &dst.FieldList{},
			Results: nil,
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				// h := handler.NewTest(r.app)
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("h")},
					Tok: token.DEFINE,
					Rhs: []dst.Expr{
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent("handler"),
								Sel: dst.NewIdent(fmt.Sprintf("New%s", serviceName)),
							},
							Args: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent("r"),
									Sel: dst.NewIdent("app"),
								},
							},
						},
					},
				},
				// route := r.fiber.Group("/v1/test")
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("route")},
					Tok: token.DEFINE,
					Rhs: []dst.Expr{
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X: &dst.SelectorExpr{
									X:   dst.NewIdent("r"),
									Sel: dst.NewIdent("fiber"),
								},
								Sel: dst.NewIdent("Group"),
							},
							Args: []dst.Expr{
								&dst.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(group),
								},
							},
						},
					},
				},
				helper.BodyListNewLinesDst(),
				// List
				&dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("route"),
							Sel: dst.NewIdent("Get"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: "\"/\"",
							},
							&dst.SelectorExpr{
								X:   dst.NewIdent("h"),
								Sel: dst.NewIdent(fmt.Sprintf("List%s", serviceName)),
							},
						},
					},
				},
				// Get
				&dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("route"),
							Sel: dst.NewIdent("Get"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: "\"/:id\"",
							},
							&dst.SelectorExpr{
								X:   dst.NewIdent("h"),
								Sel: dst.NewIdent(fmt.Sprintf("Get%s", serviceName)),
							},
						},
					},
				},
				// Post
				&dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("route"),
							Sel: dst.NewIdent("Post"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: "\"/\"",
							},
							&dst.SelectorExpr{
								X:   dst.NewIdent("h"),
								Sel: dst.NewIdent(fmt.Sprintf("Create%s", serviceName)),
							},
						},
					},
				},
				// Put
				&dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("route"),
							Sel: dst.NewIdent("Put"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: "\"/:id\"",
							},
							&dst.SelectorExpr{
								X:   dst.NewIdent("h"),
								Sel: dst.NewIdent(fmt.Sprintf("Update%s", serviceName)),
							},
						},
					},
				},
				// Patch
				&dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("route"),
							Sel: dst.NewIdent("Patch"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: "\"/:id\"",
							},
							&dst.SelectorExpr{
								X:   dst.NewIdent("h"),
								Sel: dst.NewIdent(fmt.Sprintf("Update%s", serviceName)),
							},
						},
					},
				},
				// Delete
				&dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("route"),
							Sel: dst.NewIdent("Delete"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: "\"/:id\"",
							},
							&dst.SelectorExpr{
								X:   dst.NewIdent("h"),
								Sel: dst.NewIdent(fmt.Sprintf("Delete%s", serviceName)),
							},
						},
					},
				},
			},
		},
	})

	imports := helper.GenImport(pkgImport...)
	decls = append([]dst.Decl{imports.GenDeclDst()}, decls...)
	file := &dst.File{
		Name:    dst.NewIdent("routes"),
		Imports: imports.ImportSpecDst(),
		Decls:   decls,
	}

	res <- config2.Builder{
		DstFile:  file,
		Pathname: fmt.Sprintf("internal/rest/routes/%s_route.go", s.cfg.Name),
	}
}

func (s *addService) addToRoutesDst(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)

	// Parse file routes
	resp := config2.Builder{}
	defer func() {
		res <- resp
	}()

	fset := token.NewFileSet()
	file, err := decorator.ParseFile(fset, s.app.Dir("internal/rest/routes/routes.go"), nil, parser.AllErrors)
	if err != nil {
		resp.Err = err
		return
	}
	var found bool

	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*dst.FuncDecl); ok {
			if funcDecl.Name.Name == "initRoutes" {
				_, _ = color.New(color.FgBlue).Printf("Add to routes\n")
				found = true
				// Tambahkan routes
				stmt := &dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("r"),
							Sel: dst.NewIdent(fmt.Sprintf("routes%s", serviceName)),
						},
					},
				}

				// Append ke Body.List
				funcDecl.Body.List = append(funcDecl.Body.List, stmt)
				resp.DstFile = file
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
	file, err := decorator.ParseFile(fset, s.app.Dir(fmt.Sprintf("internal/rest/routes/%s_route.go", s.cfg.ServiceName)), nil, parser.AllErrors)
	if err != nil {
		resp.Err = err
		return
	}

	decls := file.Decls[len(file.Decls)-1]
	funcDecl, ok := decls.(*dst.FuncDecl)
	if !ok {
		resp.Err = fmt.Errorf("invalid routes: function routes%s not found", serviceName)
		return
	}
	funcDecl.Body.List = append(funcDecl.Body.List, &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun: &dst.SelectorExpr{
				X:   dst.NewIdent("route"),
				Sel: dst.NewIdent(s.cfg.Method),
			},
			Args: []dst.Expr{
				&dst.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, s.cfg.Path),
				},
				&dst.SelectorExpr{
					X:   dst.NewIdent("h"),
					Sel: dst.NewIdent(s.cfg.Name),
				},
			},
		},
	})
	if s.cfg.Method == "Put" {
		funcDecl.Body.List = append(funcDecl.Body.List, &dst.ExprStmt{
			X: &dst.CallExpr{
				Fun: &dst.SelectorExpr{
					X:   dst.NewIdent("route"),
					Sel: dst.NewIdent("Patch"),
				},
				Args: []dst.Expr{
					&dst.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s"`, s.cfg.Path),
					},
					&dst.SelectorExpr{
						X:   dst.NewIdent("h"),
						Sel: dst.NewIdent(s.cfg.Name),
					},
				},
			},
		})
	}

	file.Decls[len(file.Decls)-1] = funcDecl

	resp = config2.Builder{
		DstFile:  file,
		Pathname: fmt.Sprintf("internal/rest/routes/%s_route.go", s.cfg.ServiceName),
	}
}
