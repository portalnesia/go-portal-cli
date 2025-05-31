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

func (s *addService) addServiceUsecase(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)
	ins := strings.ToLower(s.cfg.Name)[0:1]

	_, _ = color.New(color.FgBlue).Printf("Generating service\n")
	pkgImport := []string{
		`"errors"`,
		fmt.Sprintf(`context2 "%s/internal/context"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/config"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/request"`, s.cfg.Module),
	}
	decls := make([]dst.Decl, 0)

	// type Test struct { app *config.App }
	decls = append(decls, &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(serviceName),
				Type: &dst.StructType{
					Fields: &dst.FieldList{
						List: []*dst.Field{
							{
								Names: []*dst.Ident{dst.NewIdent("app")},
								Type: &dst.SelectorExpr{
									X:   dst.NewIdent("config"),
									Sel: dst.NewIdent("App"),
								},
							},
						},
					},
				},
			},
		},
	})

	// func NewTest(app *config.App) *Test { return &Test{app} }
	decls = append(decls, &dst.FuncDecl{
		Name: dst.NewIdent(fmt.Sprintf("New%s", serviceName)),
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("app")},
						Type: &dst.SelectorExpr{
							X:   dst.NewIdent("config"),
							Sel: dst.NewIdent("App"),
						},
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: &dst.StarExpr{X: dst.NewIdent(serviceName)}},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.UnaryExpr{
							Op: token.AND,
							X: &dst.CompositeLit{
								Type: dst.NewIdent(serviceName),
								Elts: []dst.Expr{dst.NewIdent("app")},
							},
						},
					},
				},
			},
		},
	})

	// func (t *Test) GetTest() (any, error) { return nil, errors.New("method not implemented") }
	decls = append(decls, &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(ins)},
					Type:  &dst.StarExpr{X: dst.NewIdent(serviceName)},
				},
			},
		},
		Name: dst.NewIdent(fmt.Sprintf("Get%s", serviceName)),
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("ctx")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("context2"),
								Sel: dst.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("query")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("request"),
								Sel: dst.NewIdent("Request"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("id")},
						Type:  dst.NewIdent("string"),
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: dst.NewIdent("any")},
					{Type: dst.NewIdent("error")},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						dst.NewIdent("nil"),
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent("errors"),
								Sel: dst.NewIdent("New"),
							},
							Args: []dst.Expr{
								&dst.BasicLit{
									Kind:  token.STRING,
									Value: "\"method not implemented\"",
								},
							},
						},
					},
				},
			},
		},
	})

	// func (t *Test) ListTest() (any, error) { return nil, errors.New("method not implemented") }
	decls = append(decls, &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(ins)},
					Type:  &dst.StarExpr{X: dst.NewIdent(serviceName)},
				},
			},
		},
		Name: dst.NewIdent(fmt.Sprintf("List%s", serviceName)),
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("ctx")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("context2"),
								Sel: dst.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("query")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("request"),
								Sel: dst.NewIdent("Request"),
							},
						},
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: dst.NewIdent("any")},
					{Type: dst.NewIdent("error")},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						dst.NewIdent("nil"),
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent("errors"),
								Sel: dst.NewIdent("New"),
							},
							Args: []dst.Expr{
								&dst.BasicLit{
									Kind:  token.STRING,
									Value: "\"method not implemented\"",
								},
							},
						},
					},
				},
			},
		},
	})

	// func (t *Test) CreateTest() (any, error) { return nil, errors.New("method not implemented") }
	decls = append(decls, &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(ins)},
					Type:  &dst.StarExpr{X: dst.NewIdent(serviceName)},
				},
			},
		},
		Name: dst.NewIdent(fmt.Sprintf("Create%s", serviceName)),
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("ctx")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("context2"),
								Sel: dst.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("query")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("request"),
								Sel: dst.NewIdent("Request"),
							},
						},
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: dst.NewIdent("any")},
					{Type: dst.NewIdent("error")},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						dst.NewIdent("nil"),
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent("errors"),
								Sel: dst.NewIdent("New"),
							},
							Args: []dst.Expr{
								&dst.BasicLit{
									Kind:  token.STRING,
									Value: "\"method not implemented\"",
								},
							},
						},
					},
				},
			},
		},
	})

	// func (t *Test) UpdateTest() (any, error) { return nil, errors.New("method not implemented") }
	decls = append(decls, &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(ins)},
					Type:  &dst.StarExpr{X: dst.NewIdent(serviceName)},
				},
			},
		},
		Name: dst.NewIdent(fmt.Sprintf("Update%s", serviceName)),
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("ctx")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("context2"),
								Sel: dst.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("query")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("request"),
								Sel: dst.NewIdent("Request"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("id")},
						Type:  dst.NewIdent("string"),
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: dst.NewIdent("any")},
					{Type: dst.NewIdent("error")},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						dst.NewIdent("nil"),
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent("errors"),
								Sel: dst.NewIdent("New"),
							},
							Args: []dst.Expr{
								&dst.BasicLit{
									Kind:  token.STRING,
									Value: "\"method not implemented\"",
								},
							},
						},
					},
				},
			},
		},
	})

	// func (t *Test) DeleteTest() error { return errors.New("method not implemented") }
	decls = append(decls, &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(ins)},
					Type:  &dst.StarExpr{X: dst.NewIdent(serviceName)},
				},
			},
		},
		Name: dst.NewIdent(fmt.Sprintf("Delete%s", serviceName)),
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("ctx")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("context2"),
								Sel: dst.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("query")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("request"),
								Sel: dst.NewIdent("Request"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("id")},
						Type:  dst.NewIdent("string"),
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: dst.NewIdent("error")},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent("errors"),
								Sel: dst.NewIdent("New"),
							},
							Args: []dst.Expr{
								&dst.BasicLit{
									Kind:  token.STRING,
									Value: "\"method not implemented\"",
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
		Name:    dst.NewIdent("service"),
		Imports: imports.ImportSpecDst(),
		Decls:   decls,
	}

	res <- config2.Builder{
		DstFile:  file,
		Pathname: fmt.Sprintf("internal/service/%s_service.go", s.cfg.Name),
	}
}

func (s *addEndpoint) addEndpointUsecase(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	_, _ = color.New(color.FgBlue).Printf("Generating service\n")

	serviceName := utils.Ucwords(s.cfg.ServiceName)
	ins := strings.ToLower(s.cfg.ServiceName)[0:1]

	// Parse file routes
	resp := config2.Builder{}
	defer func() {
		res <- resp
	}()

	fset := token.NewFileSet()
	file, err := decorator.ParseFile(fset, s.app.Dir(fmt.Sprintf("internal/service/%s_service.go", s.cfg.ServiceName)), nil, parser.AllErrors)
	if err != nil {
		resp.Err = err
		return
	}

	decl := &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(ins)},
					Type:  &dst.StarExpr{X: dst.NewIdent(serviceName)},
				},
			},
		},
		Name: dst.NewIdent(s.cfg.Name),
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("ctx")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("context2"),
								Sel: dst.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*dst.Ident{dst.NewIdent("query")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("request"),
								Sel: dst.NewIdent("Request"),
							},
						},
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: dst.NewIdent("any")},
					{Type: dst.NewIdent("error")},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						dst.NewIdent("nil"),
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent("errors"),
								Sel: dst.NewIdent("New"),
							},
							Args: []dst.Expr{
								&dst.BasicLit{
									Kind:  token.STRING,
									Value: "\"method not implemented\"",
								},
							},
						},
					},
				},
			},
		},
	}
	decl.Decorations().Before = dst.EmptyLine
	file.Decls = append(file.Decls, decl)

	resp = config2.Builder{
		DstFile:  file,
		Pathname: fmt.Sprintf("internal/service/%s_service.go", s.cfg.ServiceName),
	}
}
