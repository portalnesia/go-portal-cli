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

func (s *addService) addServiceHandler(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	serviceName := utils.Ucwords(s.cfg.Name)
	ins := strings.ToLower(s.cfg.Name)[0:1]

	_, _ = color.New(color.FgBlue).Printf("Generating handler\n")
	pkgImport := []string{
		`"github.com/gofiber/fiber/v2"`,
		fmt.Sprintf(`"%s/internal/context"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/dto"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/config"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/service"`, s.cfg.Module),
	}
	decls := make([]dst.Decl, 0)

	// Struct type
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
							{
								Names: []*dst.Ident{dst.NewIdent("s")},
								Type: &dst.StarExpr{
									X: &dst.SelectorExpr{
										X:   dst.NewIdent("service"),
										Sel: dst.NewIdent(serviceName),
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
					{
						Type: &dst.StarExpr{
							X: dst.NewIdent(serviceName),
						},
					},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				// u := usecase.NewTest(app)
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("s")},
					Tok: token.DEFINE,
					Rhs: []dst.Expr{
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent("service"),
								Sel: dst.NewIdent(fmt.Sprintf("New%s", serviceName)),
							},
							Args: []dst.Expr{
								dst.NewIdent("app"),
							},
						},
					},
				},
				// return &Test{app, u}
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.UnaryExpr{
							Op: token.AND,
							X: &dst.CompositeLit{
								Type: dst.NewIdent(serviceName),
								Elts: []dst.Expr{
									dst.NewIdent("app"),
									dst.NewIdent("s"),
								},
							},
						},
					},
				},
			},
		},
	})

	// Get
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
						Names: []*dst.Ident{dst.NewIdent("c")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("fiber"),
								Sel: dst.NewIdent("Ctx"),
							},
						},
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
							Fun: dst.NewIdent("newHandler"),
							Args: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent(ins),
									Sel: dst.NewIdent("app"),
								},
								dst.NewIdent("c"),
								&dst.FuncLit{
									Type: &dst.FuncType{
										Params: &dst.FieldList{
											List: []*dst.Field{
												{
													Names: []*dst.Ident{dst.NewIdent("ctx")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("context"),
															Sel: dst.NewIdent("Context"),
														},
													},
												},
												{
													Names: []*dst.Ident{dst.NewIdent("query")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("dto"),
															Sel: dst.NewIdent("Request"),
														},
													},
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
											// id := c.Params("id")
											&dst.AssignStmt{
												Lhs: []dst.Expr{dst.NewIdent("id")},
												Tok: token.DEFINE,
												Rhs: []dst.Expr{
													&dst.CallExpr{
														Fun: &dst.SelectorExpr{
															X:   dst.NewIdent("c"),
															Sel: dst.NewIdent("Params"),
														},
														Args: []dst.Expr{
															&dst.BasicLit{
																Kind:  token.STRING,
																Value: "\"id\"",
															},
														},
													},
												},
											},
											// data, err := t.u.GetTest(ctx, query)
											&dst.AssignStmt{
												Lhs: []dst.Expr{
													dst.NewIdent("data"),
													dst.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []dst.Expr{
													&dst.CallExpr{
														Fun: &dst.SelectorExpr{
															X: &dst.SelectorExpr{
																X:   dst.NewIdent(ins),
																Sel: dst.NewIdent("s"),
															},
															Sel: dst.NewIdent(fmt.Sprintf("Get%s", serviceName)),
														},
														Args: []dst.Expr{
															dst.NewIdent("ctx"),
															dst.NewIdent("query"),
															dst.NewIdent("id"),
														},
													},
												},
											},
											// if err != nil { return err }
											&dst.IfStmt{
												Cond: &dst.BinaryExpr{
													X:  dst.NewIdent("err"),
													Op: token.NEQ,
													Y:  dst.NewIdent("nil"),
												},
												Body: &dst.BlockStmt{
													List: []dst.Stmt{
														&dst.ReturnStmt{
															Results: []dst.Expr{dst.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, data)
											&dst.ReturnStmt{
												Results: []dst.Expr{
													&dst.CallExpr{
														Fun: dst.NewIdent("newResponse"),
														Args: []dst.Expr{
															dst.NewIdent("c"),
															dst.NewIdent("data"),
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
						Names: []*dst.Ident{dst.NewIdent("c")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("fiber"),
								Sel: dst.NewIdent("Ctx"),
							},
						},
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
							Fun: dst.NewIdent("newHandler"),
							Args: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent(ins),
									Sel: dst.NewIdent("app"),
								},
								dst.NewIdent("c"),
								&dst.FuncLit{
									Type: &dst.FuncType{
										Params: &dst.FieldList{
											List: []*dst.Field{
												{
													Names: []*dst.Ident{dst.NewIdent("ctx")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("context"),
															Sel: dst.NewIdent("Context"),
														},
													},
												},
												{
													Names: []*dst.Ident{dst.NewIdent("query")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("dto"),
															Sel: dst.NewIdent("Request"),
														},
													},
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
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&dst.AssignStmt{
												Lhs: []dst.Expr{
													dst.NewIdent("data"),
													dst.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []dst.Expr{
													&dst.CallExpr{
														Fun: &dst.SelectorExpr{
															X: &dst.SelectorExpr{
																X:   dst.NewIdent(ins),
																Sel: dst.NewIdent("s"),
															},
															Sel: dst.NewIdent(fmt.Sprintf("List%s", serviceName)),
														},
														Args: []dst.Expr{
															dst.NewIdent("ctx"),
															dst.NewIdent("query"),
														},
													},
												},
											},
											// if err != nil { return err }
											&dst.IfStmt{
												Cond: &dst.BinaryExpr{
													X:  dst.NewIdent("err"),
													Op: token.NEQ,
													Y:  dst.NewIdent("nil"),
												},
												Body: &dst.BlockStmt{
													List: []dst.Stmt{
														&dst.ReturnStmt{
															Results: []dst.Expr{dst.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, data)
											&dst.ReturnStmt{
												Results: []dst.Expr{
													&dst.CallExpr{
														Fun: dst.NewIdent("newResponse"),
														Args: []dst.Expr{
															dst.NewIdent("c"),
															dst.NewIdent("data"),
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
						Names: []*dst.Ident{dst.NewIdent("c")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("fiber"),
								Sel: dst.NewIdent("Ctx"),
							},
						},
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
							Fun: dst.NewIdent("newHandler"),
							Args: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent(ins),
									Sel: dst.NewIdent("app"),
								},
								dst.NewIdent("c"),
								&dst.FuncLit{
									Type: &dst.FuncType{
										Params: &dst.FieldList{
											List: []*dst.Field{
												{
													Names: []*dst.Ident{dst.NewIdent("ctx")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("context"),
															Sel: dst.NewIdent("Context"),
														},
													},
												},
												{
													Names: []*dst.Ident{dst.NewIdent("query")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("dto"),
															Sel: dst.NewIdent("Request"),
														},
													},
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
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&dst.AssignStmt{
												Lhs: []dst.Expr{
													dst.NewIdent("data"),
													dst.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []dst.Expr{
													&dst.CallExpr{
														Fun: &dst.SelectorExpr{
															X: &dst.SelectorExpr{
																X:   dst.NewIdent(ins),
																Sel: dst.NewIdent("s"),
															},
															Sel: dst.NewIdent(fmt.Sprintf("Create%s", serviceName)),
														},
														Args: []dst.Expr{
															dst.NewIdent("ctx"),
															dst.NewIdent("query"),
														},
													},
												},
											},
											// if err != nil { return err }
											&dst.IfStmt{
												Cond: &dst.BinaryExpr{
													X:  dst.NewIdent("err"),
													Op: token.NEQ,
													Y:  dst.NewIdent("nil"),
												},
												Body: &dst.BlockStmt{
													List: []dst.Stmt{
														&dst.ReturnStmt{
															Results: []dst.Expr{dst.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, data)
											&dst.ReturnStmt{
												Results: []dst.Expr{
													&dst.CallExpr{
														Fun: dst.NewIdent("newResponse"),
														Args: []dst.Expr{
															dst.NewIdent("c"),
															dst.NewIdent("data"),
															&dst.BasicLit{
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
						Names: []*dst.Ident{dst.NewIdent("c")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("fiber"),
								Sel: dst.NewIdent("Ctx"),
							},
						},
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
							Fun: dst.NewIdent("newHandler"),
							Args: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent(ins),
									Sel: dst.NewIdent("app"),
								},
								dst.NewIdent("c"),
								&dst.FuncLit{
									Type: &dst.FuncType{
										Params: &dst.FieldList{
											List: []*dst.Field{
												{
													Names: []*dst.Ident{dst.NewIdent("ctx")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("context"),
															Sel: dst.NewIdent("Context"),
														},
													},
												},
												{
													Names: []*dst.Ident{dst.NewIdent("query")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("dto"),
															Sel: dst.NewIdent("Request"),
														},
													},
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
											// id := c.Params("id")
											&dst.AssignStmt{
												Lhs: []dst.Expr{dst.NewIdent("id")},
												Tok: token.DEFINE,
												Rhs: []dst.Expr{
													&dst.CallExpr{
														Fun: &dst.SelectorExpr{
															X:   dst.NewIdent("c"),
															Sel: dst.NewIdent("Params"),
														},
														Args: []dst.Expr{
															&dst.BasicLit{
																Kind:  token.STRING,
																Value: "\"id\"",
															},
														},
													},
												},
											},
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&dst.AssignStmt{
												Lhs: []dst.Expr{
													dst.NewIdent("data"),
													dst.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []dst.Expr{
													&dst.CallExpr{
														Fun: &dst.SelectorExpr{
															X: &dst.SelectorExpr{
																X:   dst.NewIdent(ins),
																Sel: dst.NewIdent("s"),
															},
															Sel: dst.NewIdent(fmt.Sprintf("Update%s", serviceName)),
														},
														Args: []dst.Expr{
															dst.NewIdent("ctx"),
															dst.NewIdent("query"),
															dst.NewIdent("id"),
														},
													},
												},
											},
											// if err != nil { return err }
											&dst.IfStmt{
												Cond: &dst.BinaryExpr{
													X:  dst.NewIdent("err"),
													Op: token.NEQ,
													Y:  dst.NewIdent("nil"),
												},
												Body: &dst.BlockStmt{
													List: []dst.Stmt{
														&dst.ReturnStmt{
															Results: []dst.Expr{dst.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, data)
											&dst.ReturnStmt{
												Results: []dst.Expr{
													&dst.CallExpr{
														Fun: dst.NewIdent("newResponse"),
														Args: []dst.Expr{
															dst.NewIdent("c"),
															dst.NewIdent("data"),
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
						Names: []*dst.Ident{dst.NewIdent("c")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("fiber"),
								Sel: dst.NewIdent("Ctx"),
							},
						},
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
							Fun: dst.NewIdent("newHandler"),
							Args: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent(ins),
									Sel: dst.NewIdent("app"),
								},
								dst.NewIdent("c"),
								&dst.FuncLit{
									Type: &dst.FuncType{
										Params: &dst.FieldList{
											List: []*dst.Field{
												{
													Names: []*dst.Ident{dst.NewIdent("ctx")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("context"),
															Sel: dst.NewIdent("Context"),
														},
													},
												},
												{
													Names: []*dst.Ident{dst.NewIdent("query")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("dto"),
															Sel: dst.NewIdent("Request"),
														},
													},
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
											// id := c.Params("id")
											&dst.AssignStmt{
												Lhs: []dst.Expr{dst.NewIdent("id")},
												Tok: token.DEFINE,
												Rhs: []dst.Expr{
													&dst.CallExpr{
														Fun: &dst.SelectorExpr{
															X:   dst.NewIdent("c"),
															Sel: dst.NewIdent("Params"),
														},
														Args: []dst.Expr{
															&dst.BasicLit{
																Kind:  token.STRING,
																Value: "\"id\"",
															},
														},
													},
												},
											},
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&dst.AssignStmt{
												Lhs: []dst.Expr{
													dst.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []dst.Expr{
													&dst.CallExpr{
														Fun: &dst.SelectorExpr{
															X: &dst.SelectorExpr{
																X:   dst.NewIdent(ins),
																Sel: dst.NewIdent("s"),
															},
															Sel: dst.NewIdent(fmt.Sprintf("Delete%s", serviceName)),
														},
														Args: []dst.Expr{
															dst.NewIdent("ctx"),
															dst.NewIdent("query"),
															dst.NewIdent("id"),
														},
													},
												},
											},
											// if err != nil { return err }
											&dst.IfStmt{
												Cond: &dst.BinaryExpr{
													X:  dst.NewIdent("err"),
													Op: token.NEQ,
													Y:  dst.NewIdent("nil"),
												},
												Body: &dst.BlockStmt{
													List: []dst.Stmt{
														&dst.ReturnStmt{
															Results: []dst.Expr{dst.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, false, 204)
											&dst.ReturnStmt{
												Results: []dst.Expr{
													&dst.CallExpr{
														Fun: dst.NewIdent("newResponse"),
														Args: []dst.Expr{
															dst.NewIdent("c"),
															dst.NewIdent("false"),
															&dst.BasicLit{
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
	decls = append([]dst.Decl{imports.GenDeclDst()}, decls...)
	for i := range decls {
		decls[i].Decorations().Before = dst.EmptyLine
	}
	file := &dst.File{
		Name:    dst.NewIdent("handler"),
		Imports: imports.ImportSpecDst(),
		Decls:   decls,
	}

	res <- config2.Builder{
		DstFile:  file,
		Pathname: fmt.Sprintf("internal/rest/handler/%s_handler.go", s.cfg.Name),
	}
}

func (s *addEndpoint) addEndpointHandler(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	_, _ = color.New(color.FgBlue).Printf("Generating handler\n")

	serviceName := utils.Ucwords(s.cfg.ServiceName)
	ins := strings.ToLower(s.cfg.ServiceName)[0:1]

	// Parse file routes
	resp := config2.Builder{}
	defer func() {
		res <- resp
	}()

	fset := token.NewFileSet()
	file, err := decorator.ParseFile(fset, s.app.Dir(fmt.Sprintf("internal/rest/handler/%s_handler.go", s.cfg.ServiceName)), nil, parser.AllErrors)
	if err != nil {
		resp.Err = err
		return
	}

	file.Decls = append(file.Decls, &dst.FuncDecl{
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
						Names: []*dst.Ident{dst.NewIdent("c")},
						Type: &dst.StarExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("fiber"),
								Sel: dst.NewIdent("Ctx"),
							},
						},
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
							Fun: dst.NewIdent("newHandler"),
							Args: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent(ins),
									Sel: dst.NewIdent("app"),
								},
								dst.NewIdent("c"),
								&dst.FuncLit{
									Type: &dst.FuncType{
										Params: &dst.FieldList{
											List: []*dst.Field{
												{
													Names: []*dst.Ident{dst.NewIdent("ctx")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("context"),
															Sel: dst.NewIdent("Context"),
														},
													},
												},
												{
													Names: []*dst.Ident{dst.NewIdent("query")},
													Type: &dst.StarExpr{
														X: &dst.SelectorExpr{
															X:   dst.NewIdent("dto"),
															Sel: dst.NewIdent("Request"),
														},
													},
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
											// data, err := t.u.GetTest(ctx, trxDb, query)
											&dst.AssignStmt{
												Lhs: []dst.Expr{
													dst.NewIdent("data"),
													dst.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []dst.Expr{
													&dst.CallExpr{
														Fun: &dst.SelectorExpr{
															X: &dst.SelectorExpr{
																X:   dst.NewIdent(ins),
																Sel: dst.NewIdent("s"),
															},
															Sel: dst.NewIdent(s.cfg.Name),
														},
														Args: []dst.Expr{
															dst.NewIdent("ctx"),
															dst.NewIdent("query"),
														},
													},
												},
											},
											// if err != nil { return err }
											&dst.IfStmt{
												Cond: &dst.BinaryExpr{
													X:  dst.NewIdent("err"),
													Op: token.NEQ,
													Y:  dst.NewIdent("nil"),
												},
												Body: &dst.BlockStmt{
													List: []dst.Stmt{
														&dst.ReturnStmt{
															Results: []dst.Expr{dst.NewIdent("err")},
														},
													},
												},
											},
											// return t.Response(c, data)
											&dst.ReturnStmt{
												Results: []dst.Expr{
													&dst.CallExpr{
														Fun: dst.NewIdent("newResponse"),
														Args: []dst.Expr{
															dst.NewIdent("c"),
															dst.NewIdent("data"),
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

	for i := range file.Decls {
		file.Decls[i].Decorations().Before = dst.EmptyLine
	}

	resp = config2.Builder{
		DstFile:  file,
		Pathname: fmt.Sprintf("internal/rest/handler/%s_handler.go", s.cfg.ServiceName),
	}
}
