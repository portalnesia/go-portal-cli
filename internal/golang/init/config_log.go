/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package ginit

import (
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go/ast"
	"go/token"
	"sync"
)

// Pertama, buat helper buildLoggerLevelBody
func (c *initType) buildLoggerLevelBody(level string) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: []ast.Stmt{
			// tmp := l.logger.<level>().Str("service", service)
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("tmp"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("l"),
										Sel: ast.NewIdent("logger"),
									},
									Sel: ast.NewIdent(level),
								},
								Args: []ast.Expr{},
							},
							Sel: ast.NewIdent("Str"),
						},
						Args: []ast.Expr{
							ast.NewIdent("service"),
							ast.NewIdent("service"),
						},
					},
				},
			},
			// if len(err) > 0 { tmp = tmp.Err(err[0]) }
			&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X: &ast.CallExpr{
						Fun: ast.NewIdent("len"),
						Args: []ast.Expr{
							ast.NewIdent("err"),
						},
					},
					Op: token.GTR,
					Y: &ast.BasicLit{
						Kind:  token.INT,
						Value: "0",
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("tmp"),
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("tmp"),
										Sel: ast.NewIdent("Err"),
									},
									Args: []ast.Expr{
										&ast.IndexExpr{
											X: ast.NewIdent("err"),
											Index: &ast.BasicLit{
												Kind:  token.INT,
												Value: "0",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			// return tmp
			&ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("tmp"),
				},
			},
		},
	}
}

// Helper untuk fungsi Error dan Fatal (yang beda)
func (c *initType) buildLoggerErrorBody(level string) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("l"),
												Sel: ast.NewIdent("logger"),
											},
											Sel: ast.NewIdent(level),
										},
										Args: []ast.Expr{},
									},
									Sel: ast.NewIdent("Err"),
								},
								Args: []ast.Expr{
									ast.NewIdent("err"),
								},
							},
							Sel: ast.NewIdent("Str"),
						},
						Args: []ast.Expr{
							ast.NewIdent("service"),
							ast.NewIdent("service"),
						},
					},
				},
			},
		},
	}
}

func (c *initType) initConfigLog(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/log.go\n")

	src, _ := c.app.DataEmbed.ReadFile("data/golang/internal/config/log.txt")

	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/config/log.go",
	}
}
