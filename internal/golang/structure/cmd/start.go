/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package scmd

import (
	"fmt"
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
	"go/ast"
	"go/token"
	"sync"
)

func (s *Cmd) initCmdStart(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating cmd/start.go\n")

	pkgImport := []string{
		`"fmt"`,
		`"github.com/spf13/cobra"`,
		`"github.com/spf13/viper"`,
		`"os"`,
		`"os/signal"`,
		fmt.Sprintf(`"%s/internal/config"`, s.cfg.Module),
		fmt.Sprintf(`"%s/internal/server"`, s.cfg.Module),
		`"syscall"`,
	}
	decls := make([]ast.Decl, 0)

	var runFunc []ast.Stmt

	runFunc = append(runFunc,
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("app")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{&ast.StarExpr{X: ast.NewIdent("appConfig")}},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("app"),
					Sel: ast.NewIdent("DB"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent("true")},
		},
	)

	if s.cfg.Redis {
		runFunc = append(runFunc, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("app"),
					Sel: ast.NewIdent("Redis"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent("true")},
		})
	}
	if s.cfg.Firebase {
		runFunc = append(runFunc, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("app"),
					Sel: ast.NewIdent("Firebase"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent("true")},
		})
	}

	runFunc = append(runFunc,
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("apps")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("config"),
						Sel: ast.NewIdent("New"),
					},
					Args: []ast.Expr{ast.NewIdent("app")},
				},
			},
		},
		&ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("apps"),
					Sel: ast.NewIdent("Close"),
				},
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("fiberApp")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("server"),
						Sel: ast.NewIdent("New"),
					},
					Args: []ast.Expr{ast.NewIdent("apps")},
				},
			},
		},
		&ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("fiberApp"),
					Sel: ast.NewIdent("Close"),
				},
			},
		},
		helper.BodyListNewLines(),
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("ports")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("viper"),
						Sel: ast.NewIdent("GetIntSlice"),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"ports"`,
						},
					},
				},
			},
		},
		&ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Value: ast.NewIdent("port"),
			Tok:   token.DEFINE,
			X:     ast.NewIdent("ports"),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.GoStmt{
						Call: &ast.CallExpr{
							Fun: &ast.FuncLit{
								Type: &ast.FuncType{
									Params: &ast.FieldList{},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("apps"),
																Sel: ast.NewIdent("Log"),
															},
															Sel: ast.NewIdent("Info"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"system\"",
															},
														},
													},
													Sel: ast.NewIdent("Msgf"),
												},
												Args: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: "\"Starting server on port %d\"",
													},
													ast.NewIdent("port"),
												},
											},
										},
										&ast.AssignStmt{
											Lhs: []ast.Expr{ast.NewIdent("err")},
											Tok: token.DEFINE,
											Rhs: []ast.Expr{
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("fiberApp.Fiber"),
														Sel: ast.NewIdent("Listen"),
													},
													Args: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X:   ast.NewIdent("fmt"),
																Sel: ast.NewIdent("Sprintf"),
															},
															Args: []ast.Expr{
																&ast.BasicLit{
																	Kind:  token.STRING,
																	Value: "\"127.0.0.1:%d\"",
																},
																ast.NewIdent("port"),
															},
														},
													},
												},
											},
										},
										&ast.IfStmt{
											Cond: &ast.BinaryExpr{
												X:  ast.NewIdent("err"),
												Op: token.NEQ,
												Y:  ast.NewIdent("nil"),
											},
											Body: &ast.BlockStmt{
												List: []ast.Stmt{
													&ast.ExprStmt{
														X: &ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.SelectorExpr{
																			X:   ast.NewIdent("apps"),
																			Sel: ast.NewIdent("Log"),
																		},
																		Sel: ast.NewIdent("Error"),
																	},
																	Args: []ast.Expr{
																		ast.NewIdent("err"),
																		&ast.BasicLit{
																			Kind:  token.STRING,
																			Value: "\"system\"",
																		},
																	},
																},
																Sel: ast.NewIdent("Msgf"),
															},
															Args: []ast.Expr{
																&ast.BasicLit{
																	Kind:  token.STRING,
																	Value: "\"API Server is error when running on port %d\"",
																},
																ast.NewIdent("port"),
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
		helper.BodyListNewLines(),
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("signKill")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("make"),
					Args: []ast.Expr{
						&ast.ChanType{
							Dir: ast.SEND | ast.RECV,
							Value: &ast.SelectorExpr{
								X:   ast.NewIdent("os"),
								Sel: ast.NewIdent("Signal"),
							},
						},
						&ast.BasicLit{
							Kind:  token.INT,
							Value: "1",
						},
					},
				},
			},
		},
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("signal"),
					Sel: ast.NewIdent("Notify"),
				},
				Args: []ast.Expr{
					ast.NewIdent("signKill"),
					&ast.SelectorExpr{
						X:   ast.NewIdent("os"),
						Sel: ast.NewIdent("Interrupt"),
					},
					&ast.SelectorExpr{
						X:   ast.NewIdent("syscall"),
						Sel: ast.NewIdent("SIGINT"),
					},
					&ast.SelectorExpr{
						X:   ast.NewIdent("syscall"),
						Sel: ast.NewIdent("SIGTERM"),
					},
				},
			},
		},
		helper.BodyListNewLines(),
		&ast.ExprStmt{
			X: &ast.UnaryExpr{
				Op: token.ARROW,
				X:  ast.NewIdent("signKill"),
			},
		},
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("fmt"),
					Sel: ast.NewIdent("Print"),
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: `"\n\n=========================================\n\n"`,
					},
				},
			},
		},
	)

	decls = append(decls, &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("startCmd")},
				Values: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X: &ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("cobra"),
								Sel: ast.NewIdent("Command"),
							},
							Elts: []ast.Expr{
								// Use: "start"
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Use"),
									Value: &ast.BasicLit{
										Kind:  token.STRING,
										Value: `"start"`,
									},
								},
								// Short: "A brief description..."
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Short"),
									Value: &ast.BasicLit{
										Kind:  token.STRING,
										Value: `"A brief description of your command"`,
									},
								},
								// Long: multi-line string
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Long"),
									Value: &ast.BasicLit{
										Kind:  token.STRING,
										Value: "`A longer description that spans multiple lines and likely contains examples\nand usage of using your command. For example:\n\nCobra is a CLI library for Go that empowers applications.\nThis application is a tool to generate the needed files\nto quickly create a Cobra application.`",
									},
								},
								// Run: func(cmd *cobra.Command, args []string) { ... }
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Run"),
									Value: &ast.FuncLit{
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Names: []*ast.Ident{ast.NewIdent("cmd")},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("cobra"),
																Sel: ast.NewIdent("Command"),
															},
														},
													},
													{
														Names: []*ast.Ident{ast.NewIdent("args")},
														Type:  &ast.ArrayType{Elt: ast.NewIdent("string")},
													},
												},
											},
										},
										Body: &ast.BlockStmt{
											List: runFunc,
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
		Name:    ast.NewIdent("cmd"),
		Imports: imports.ImportSpec(),
		Decls:   decls,
	}

	res <- config2.Builder{
		File:     file,
		Pathname: "cmd/start.go",
	}
}
