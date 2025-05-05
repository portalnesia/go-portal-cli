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

func (s *addService) addServiceUsecase(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)
	ins := strings.ToLower(s.cfg.Name)[0:1]

	_, _ = color.New(color.FgBlue).Printf("Generating usecase\n")
	pkgImport := []string{
		`"errors"`,
		`"gorm.io/gorm"`,
		fmt.Sprintf(`context2 "%s/internal/context"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/request"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/server/config"`, s.cfg.Module),
	}
	decls := make([]ast.Decl, 0)

	// type Test struct { app *config.App }
	decls = append(decls, &ast.GenDecl{
		Tok: token.TYPE,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(serviceName),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("app")},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("config"),
										Sel: ast.NewIdent("App"),
									},
								},
							},
						},
					},
				},
			},
		},
	})

	// func NewTest(app *config.App) *Test { return &Test{app} }
	decls = append(decls, &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("New%s", serviceName)),
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("app")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("config"),
								Sel: ast.NewIdent("App"),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.StarExpr{X: ast.NewIdent(serviceName)}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent(serviceName),
								Elts: []ast.Expr{ast.NewIdent("app")},
							},
						},
					},
				},
			},
		},
	})

	// func (t *Test) GetTest() (any, error) { return nil, errors.New("method not implemented") }
	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("Get%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("context2"),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("trxDb")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("gorm"),
								Sel: ast.NewIdent("DB"),
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
						Names: []*ast.Ident{ast.NewIdent("id")},
						Type:  ast.NewIdent("string"),
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
	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("List%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("context2"),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("trxDb")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("gorm"),
								Sel: ast.NewIdent("DB"),
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
	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("Create%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("context2"),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("trxDb")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("gorm"),
								Sel: ast.NewIdent("DB"),
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
	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("Update%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("context2"),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("trxDb")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("gorm"),
								Sel: ast.NewIdent("DB"),
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
						Names: []*ast.Ident{ast.NewIdent("id")},
						Type:  ast.NewIdent("string"),
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
	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("Delete%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("context2"),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("trxDb")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("gorm"),
								Sel: ast.NewIdent("DB"),
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
						Names: []*ast.Ident{ast.NewIdent("id")},
						Type:  ast.NewIdent("string"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: ast.NewIdent("error")},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("errors"),
								Sel: ast.NewIdent("New"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
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
	decls = append([]ast.Decl{imports.GenDecl()}, decls...)
	file := &ast.File{
		Name:    ast.NewIdent("usecase"),
		Imports: imports.ImportSpec(),
		Decls:   decls,
	}

	res <- config2.Builder{
		File:     file,
		Pathname: fmt.Sprintf("internal/server/usecase/%s.go", s.cfg.Name),
	}
}

func (s *addEndpoint) addEndpointUsecase(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	_, _ = color.New(color.FgBlue).Printf("Generating usecase\n")

	serviceName := utils.Ucwords(s.cfg.ServiceName)
	ins := strings.ToLower(s.cfg.ServiceName)[0:1]

	// Parse file routes
	resp := config2.Builder{}
	defer func() {
		res <- resp
	}()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, s.app.Dir(fmt.Sprintf("internal/server/usecase/%s.go", s.cfg.ServiceName)), nil, parser.AllErrors)
	if err != nil {
		resp.Err = err
		return
	}

	file.Decls = append(file.Decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent(s.cfg.Name),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("context2"),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("trxDb")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("gorm"),
								Sel: ast.NewIdent("DB"),
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
									Value: "\"method not implemented\"",
								},
							},
						},
					},
				},
			},
		},
	})

	resp = config2.Builder{
		File:     file,
		Pathname: fmt.Sprintf("internal/server/usecase/%s.go", s.cfg.ServiceName),
	}
}
