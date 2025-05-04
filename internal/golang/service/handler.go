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
	"go/token"
	"strings"
	"sync"
)

func (s *newService) newServiceHandler(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)
	ins := strings.ToLower(s.cfg.Name)[0:1]

	_, _ = color.New(color.FgBlue).Printf("Generating handler\n")
	pkgImport := []string{
		`"github.com/gofiber/fiber/v2"`,
		`"gorm.io/gorm"`,
		fmt.Sprintf(`"%s/internal/context"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/request"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/server/config"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/server/usecase"`, s.cfg.Module),
	}
	decls := make([]ast.Decl, 0)

	// Struct type
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
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("config"),
										Sel: ast.NewIdent("App"),
									},
								},
							},
							{
								Names: []*ast.Ident{ast.NewIdent("u")},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("usecase"),
										Sel: ast.NewIdent(serviceName),
									},
								},
							},
						},
					},
				},
			},
		},
	})

	// Constructor function
	decls = append(decls, &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("New%s", serviceName)),
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
					{
						Type: &ast.StarExpr{
							X: ast.NewIdent(serviceName),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// u := usecase.NewTest(app)
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("u")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("usecase"),
								Sel: ast.NewIdent(fmt.Sprintf("New%s", serviceName)),
							},
							Args: []ast.Expr{
								ast.NewIdent("app"),
							},
						},
					},
				},
				// return &Test{app, u}
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent(serviceName),
								Elts: []ast.Expr{
									ast.NewIdent("app"),
									ast.NewIdent("u"),
								},
							},
						},
					},
				},
			},
		},
	})

	// Get
	decls = append(decls, &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("Get%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("c")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("fiber"),
								Sel: ast.NewIdent("Ctx"),
							},
						},
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
								X:   ast.NewIdent(ins),
								Sel: ast.NewIdent("NewService"),
							},
							Args: []ast.Expr{
								ast.NewIdent("c"),
								&ast.FuncLit{
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
													Names: []*ast.Ident{ast.NewIdent("trxDb")},
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("gorm"),
															Sel: ast.NewIdent("DB"),
														},
													},
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
											// id := c.Params("id")
											&ast.AssignStmt{
												Lhs: []ast.Expr{ast.NewIdent("id")},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("c"),
															Sel: ast.NewIdent("Params"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"id\"",
															},
														},
													},
												},
											},
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("data"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent(ins),
																Sel: ast.NewIdent("u"),
															},
															Sel: ast.NewIdent(fmt.Sprintf("Get%s", serviceName)),
														},
														Args: []ast.Expr{
															ast.NewIdent("ctx"),
															ast.NewIdent("trxDb"),
															ast.NewIdent("query"),
															ast.NewIdent("id"),
														},
													},
												},
											},
											// if err != nil { return err }
											&ast.IfStmt{
												Cond: &ast.BinaryExpr{
													X:  ast.NewIdent("err"),
													Op: token.NEQ,
													Y:  ast.NewIdent("nil"),
												},
												Body: &ast.BlockStmt{
													List: []ast.Stmt{
														&ast.ReturnStmt{
															Results: []ast.Expr{ast.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, data)
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent(ins),
															Sel: ast.NewIdent("Response"),
														},
														Args: []ast.Expr{
															ast.NewIdent("c"),
															ast.NewIdent("data"),
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
				},
			},
		},
	})

	// List
	decls = append(decls, &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("List%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("c")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("fiber"),
								Sel: ast.NewIdent("Ctx"),
							},
						},
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
								X:   ast.NewIdent(ins),
								Sel: ast.NewIdent("NewService"),
							},
							Args: []ast.Expr{
								ast.NewIdent("c"),
								&ast.FuncLit{
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
													Names: []*ast.Ident{ast.NewIdent("trxDb")},
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("gorm"),
															Sel: ast.NewIdent("DB"),
														},
													},
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
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("data"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent(ins),
																Sel: ast.NewIdent("u"),
															},
															Sel: ast.NewIdent(fmt.Sprintf("List%s", serviceName)),
														},
														Args: []ast.Expr{
															ast.NewIdent("ctx"),
															ast.NewIdent("trxDb"),
															ast.NewIdent("query"),
														},
													},
												},
											},
											// if err != nil { return err }
											&ast.IfStmt{
												Cond: &ast.BinaryExpr{
													X:  ast.NewIdent("err"),
													Op: token.NEQ,
													Y:  ast.NewIdent("nil"),
												},
												Body: &ast.BlockStmt{
													List: []ast.Stmt{
														&ast.ReturnStmt{
															Results: []ast.Expr{ast.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, data)
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent(ins),
															Sel: ast.NewIdent("Response"),
														},
														Args: []ast.Expr{
															ast.NewIdent("c"),
															ast.NewIdent("data"),
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
				},
			},
		},
	})

	// Create
	decls = append(decls, &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("Create%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("c")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("fiber"),
								Sel: ast.NewIdent("Ctx"),
							},
						},
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
								X:   ast.NewIdent(ins),
								Sel: ast.NewIdent("NewService"),
							},
							Args: []ast.Expr{
								ast.NewIdent("c"),
								&ast.FuncLit{
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
													Names: []*ast.Ident{ast.NewIdent("trxDb")},
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("gorm"),
															Sel: ast.NewIdent("DB"),
														},
													},
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
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("data"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent(ins),
																Sel: ast.NewIdent("u"),
															},
															Sel: ast.NewIdent(fmt.Sprintf("Create%s", serviceName)),
														},
														Args: []ast.Expr{
															ast.NewIdent("ctx"),
															ast.NewIdent("trxDb"),
															ast.NewIdent("query"),
														},
													},
												},
											},
											// if err != nil { return err }
											&ast.IfStmt{
												Cond: &ast.BinaryExpr{
													X:  ast.NewIdent("err"),
													Op: token.NEQ,
													Y:  ast.NewIdent("nil"),
												},
												Body: &ast.BlockStmt{
													List: []ast.Stmt{
														&ast.ReturnStmt{
															Results: []ast.Expr{ast.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, data)
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent(ins),
															Sel: ast.NewIdent("Response"),
														},
														Args: []ast.Expr{
															ast.NewIdent("c"),
															ast.NewIdent("data"),
															&ast.BasicLit{
																Kind:  token.INT,
																Value: "201",
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
					},
				},
			},
		},
	})

	// Update
	decls = append(decls, &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("Update%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("c")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("fiber"),
								Sel: ast.NewIdent("Ctx"),
							},
						},
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
								X:   ast.NewIdent(ins),
								Sel: ast.NewIdent("NewService"),
							},
							Args: []ast.Expr{
								ast.NewIdent("c"),
								&ast.FuncLit{
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
													Names: []*ast.Ident{ast.NewIdent("trxDb")},
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("gorm"),
															Sel: ast.NewIdent("DB"),
														},
													},
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
											// id := c.Params("id")
											&ast.AssignStmt{
												Lhs: []ast.Expr{ast.NewIdent("id")},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("c"),
															Sel: ast.NewIdent("Params"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"id\"",
															},
														},
													},
												},
											},
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("data"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent(ins),
																Sel: ast.NewIdent("u"),
															},
															Sel: ast.NewIdent(fmt.Sprintf("Update%s", serviceName)),
														},
														Args: []ast.Expr{
															ast.NewIdent("ctx"),
															ast.NewIdent("trxDb"),
															ast.NewIdent("query"),
															ast.NewIdent("id"),
														},
													},
												},
											},
											// if err != nil { return err }
											&ast.IfStmt{
												Cond: &ast.BinaryExpr{
													X:  ast.NewIdent("err"),
													Op: token.NEQ,
													Y:  ast.NewIdent("nil"),
												},
												Body: &ast.BlockStmt{
													List: []ast.Stmt{
														&ast.ReturnStmt{
															Results: []ast.Expr{ast.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, data)
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent(ins),
															Sel: ast.NewIdent("Response"),
														},
														Args: []ast.Expr{
															ast.NewIdent("c"),
															ast.NewIdent("data"),
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
				},
			},
		},
	})

	// Delete
	decls = append(decls, &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(ins)},
					Type:  &ast.StarExpr{X: ast.NewIdent(serviceName)},
				},
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("Delete%s", serviceName)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("c")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("fiber"),
								Sel: ast.NewIdent("Ctx"),
							},
						},
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
								X:   ast.NewIdent(ins),
								Sel: ast.NewIdent("NewService"),
							},
							Args: []ast.Expr{
								ast.NewIdent("c"),
								&ast.FuncLit{
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
													Names: []*ast.Ident{ast.NewIdent("trxDb")},
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("gorm"),
															Sel: ast.NewIdent("DB"),
														},
													},
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
											// id := c.Params("id")
											&ast.AssignStmt{
												Lhs: []ast.Expr{ast.NewIdent("id")},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("c"),
															Sel: ast.NewIdent("Params"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"id\"",
															},
														},
													},
												},
											},
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent(ins),
																Sel: ast.NewIdent("u"),
															},
															Sel: ast.NewIdent(fmt.Sprintf("Delete%s", serviceName)),
														},
														Args: []ast.Expr{
															ast.NewIdent("ctx"),
															ast.NewIdent("trxDb"),
															ast.NewIdent("query"),
															ast.NewIdent("id"),
														},
													},
												},
											},
											// if err != nil { return err }
											&ast.IfStmt{
												Cond: &ast.BinaryExpr{
													X:  ast.NewIdent("err"),
													Op: token.NEQ,
													Y:  ast.NewIdent("nil"),
												},
												Body: &ast.BlockStmt{
													List: []ast.Stmt{
														&ast.ReturnStmt{
															Results: []ast.Expr{ast.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, false, 204)
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent(ins),
															Sel: ast.NewIdent("Response"),
														},
														Args: []ast.Expr{
															ast.NewIdent("c"),
															ast.NewIdent("false"),
															&ast.BasicLit{
																Kind:  token.INT,
																Value: "204",
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
					},
				},
			},
		},
	})

	imports := helper.GenImport(pkgImport...)
	decls = append([]ast.Decl{imports.GenDecl()}, decls...)
	file := &ast.File{
		Name:    ast.NewIdent("handler"),
		Imports: imports.ImportSpec(),
		Decls:   decls,
	}

	res <- config2.Builder{
		File:     file,
		Pathname: fmt.Sprintf("internal/server/handler/%s.go", s.cfg.Name),
	}
}
