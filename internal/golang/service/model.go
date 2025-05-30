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
	"go.portalnesia.com/utils"
	"go/ast"
	"go/token"
	"strings"
	"sync"
)

func (s *addRepository) addRepositoryModel(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	_, _ = color.New(color.FgBlue).Printf("Generating model\n")
	serviceName := utils.Ucwords(s.cfg.Name)

	decls := make([]ast.Decl, 0)

	// type User struct{}
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
					Fields: &ast.FieldList{},
				},
			},
		},
	})

	// func (User) TableName() string { return "user" }
	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: ast.NewIdent(serviceName),
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent("TableName"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: ast.NewIdent("string")},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s"`, strings.ToLower(s.cfg.Name)),
						},
					},
				},
			},
		},
	})

	// func (User) GetDefaultOrder() []string { return []string{"created_at", "desc"} }
	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: ast.NewIdent(serviceName),
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent("GetDefaultOrder"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("string"),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.ArrayType{
								Elt: ast.NewIdent("string"),
							},
							Elts: []ast.Expr{
								&ast.BasicLit{Kind: token.STRING, Value: `"created_at"`},
								&ast.BasicLit{Kind: token.STRING, Value: `"desc"`},
							},
						},
					},
				},
			},
		},
	})

	// func (User) GetAvailableOrder() [][]string { return [][]string{{"created_at", "desc"}} }
	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: ast.NewIdent("User"),
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent("GetAvailableOrder"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: &ast.ArrayType{
								Elt: ast.NewIdent("string"),
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.ArrayType{
								Elt: &ast.ArrayType{
									Elt: ast.NewIdent("string"),
								},
							},
							Elts: []ast.Expr{
								&ast.CompositeLit{
									Elts: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"created_at"`},
										&ast.BasicLit{Kind: token.STRING, Value: `"desc"`},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	file := &ast.File{
		Name:  ast.NewIdent("model"),
		Decls: decls,
	}

	res <- config2.Builder{
		File:     file,
		Pathname: fmt.Sprintf("internal/model/%s.go", strings.ToLower(s.cfg.Name)),
	}
}
