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
	"go.portalnesia.com/portal-cli/pkg/helper"
	"go/ast"
	"go/token"
	"sync"
)

func (c *initType) initConfigDatabase(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/database.go\n")

	pkgImport := []string{
		`"fmt"`,
		`"github.com/dromara/carbon/v2"`,
		`nativeLog "log"`,
		`"os"`,
		`"time"`,
		`"github.com/spf13/viper"`,
		`"gorm.io/gorm"`,
		`"gorm.io/gorm/logger"`,
		`othermysql "github.com/go-sql-driver/mysql"`,
		`"gorm.io/driver/mysql"`,
	}
	decls := make([]ast.Decl, 0)

	decls = append(decls, &ast.FuncDecl{
		Name: ast.NewIdent("getDatabaseLogger"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("silent")},
						Type:  &ast.Ellipsis{Elt: ast.NewIdent("bool")},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("logger.Interface"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("level")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("logger"),
							Sel: ast.NewIdent("Error"),
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.CallExpr{
							Fun:  ast.NewIdent("len"),
							Args: []ast.Expr{ast.NewIdent("silent")},
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
								Lhs: []ast.Expr{ast.NewIdent("level")},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("logger"),
										Sel: ast.NewIdent("Silent"),
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("logger"),
								Sel: ast.NewIdent("New"),
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("nativeLog"),
										Sel: ast.NewIdent("New"),
									},
									Args: []ast.Expr{
										ast.NewIdent("os.Stdout"),
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"\\r\\n\"",
										},
										&ast.SelectorExpr{
											X:   ast.NewIdent("nativeLog"),
											Sel: ast.NewIdent("LstdFlags"),
										},
									},
								},
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("logger"),
										Sel: ast.NewIdent("Config"),
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: ast.NewIdent("SlowThreshold"),
											Value: &ast.SelectorExpr{
												X:   ast.NewIdent("time"),
												Sel: ast.NewIdent("Second"),
											},
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("LogLevel"),
											Value: ast.NewIdent("level"),
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("IgnoreRecordNotFoundError"),
											Value: ast.NewIdent("true"),
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("Colorful"),
											Value: ast.NewIdent("true"),
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

	// MYSQL
	decls = append(decls, &ast.FuncDecl{
		Name: ast.NewIdent("connectMysql"),
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("dbChan")},
						Type: &ast.ChanType{
							Dir: ast.SEND,
							Value: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("gorm"),
									Sel: ast.NewIdent("DB"),
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// mysqlconfig := othermysql.Config{...}
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("mysqlconfig")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("othermysql"),
								Sel: ast.NewIdent("Config"),
							},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key: ast.NewIdent("User"),
									Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("viper"),
											Sel: ast.NewIdent("GetString"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"db.user\"",
											},
										},
									},
								},
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Passwd"),
									Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("viper"),
											Sel: ast.NewIdent("GetString"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"db.password\"",
											},
										},
									},
								},
								&ast.KeyValueExpr{
									Key: ast.NewIdent("DBName"),
									Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("viper"),
											Sel: ast.NewIdent("GetString"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"db.database\"",
											},
										},
									},
								},
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Net"),
									Value: &ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"tcp\"",
									},
								},
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Addr"),
									Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("fmt"),
											Sel: ast.NewIdent("Sprintf"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"%s:%d\"",
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("viper"),
													Sel: ast.NewIdent("GetString"),
												},
												Args: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: "\"db.host\"",
													},
												},
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("viper"),
													Sel: ast.NewIdent("GetInt"),
												},
												Args: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: "\"db.port\"",
													},
												},
											},
										},
									},
								},
								&ast.KeyValueExpr{
									Key:   ast.NewIdent("ParseTime"),
									Value: ast.NewIdent("true"),
								},
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Loc"),
									Value: &ast.SelectorExpr{
										X:   ast.NewIdent("time"),
										Sel: ast.NewIdent("UTC"),
									},
								},
								&ast.KeyValueExpr{
									Key:   ast.NewIdent("AllowNativePasswords"),
									Value: ast.NewIdent("true"),
								},
							},
						},
					},
				},

				// dsn := mysqlconfig.FormatDSN()
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("dsn")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("mysqlconfig"),
								Sel: ast.NewIdent("FormatDSN"),
							},
						},
					},
				},

				// portalnesia := mysql.New(mysql.Config{DSN: dsn})
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("portalnesia")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("mysql"),
								Sel: ast.NewIdent("New"),
							},
							Args: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("mysql"),
										Sel: ast.NewIdent("Config"),
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("DSN"),
											Value: ast.NewIdent("dsn"),
										},
									},
								},
							},
						},
					},
				},

				// database, err := gorm.Open(portalnesia, &gorm.Config{...})
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("database"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("gorm"),
								Sel: ast.NewIdent("Open"),
							},
							Args: []ast.Expr{
								ast.NewIdent("portalnesia"),
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("gorm"),
											Sel: ast.NewIdent("Config"),
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("PrepareStmt"),
												Value: ast.NewIdent("true"),
											},
											&ast.KeyValueExpr{
												Key: ast.NewIdent("Logger"),
												Value: &ast.CallExpr{
													Fun: ast.NewIdent("getDatabaseLogger"),
												},
											},
											&ast.KeyValueExpr{
												Key: ast.NewIdent("NowFunc"),
												Value: &ast.FuncLit{
													Type: &ast.FuncType{
														Params: &ast.FieldList{},
														Results: &ast.FieldList{
															List: []*ast.Field{
																{
																	Type: &ast.SelectorExpr{
																		X:   ast.NewIdent("time"),
																		Sel: ast.NewIdent("Time"),
																	},
																},
															},
														},
													},
													Body: &ast.BlockStmt{
														List: []ast.Stmt{
															&ast.ReturnStmt{
																Results: []ast.Expr{
																	&ast.CallExpr{
																		Fun: &ast.SelectorExpr{
																			X: &ast.CallExpr{
																				Fun: &ast.SelectorExpr{
																					X:   ast.NewIdent("carbon"),
																					Sel: ast.NewIdent("Now"),
																				},
																			},
																			Sel: ast.NewIdent("StdTime"),
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
				},

				// if err != nil { Log.Fatal(err, "system").Msg("Failed to initialize mysql") }
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
												X:   ast.NewIdent("log"),
												Sel: ast.NewIdent("Fatal"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"system\"",
												},
											},
										},
										Sel: ast.NewIdent("Msg"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"Failed to initialize mysql\"",
										},
									},
								},
							},
						},
					},
				},

				// dbChan <- database
				&ast.SendStmt{
					Chan:  ast.NewIdent("dbChan"),
					Value: ast.NewIdent("database"),
				},

				// return
				&ast.ReturnStmt{},
			},
		},
	})

	imports := helper.GenImport(pkgImport...)
	decls = append([]ast.Decl{imports.GenDecl()}, decls...)
	file := &ast.File{
		Package: token.Pos(1),
		Name:    ast.NewIdent("config"),
		Imports: imports.ImportSpec(),
		Decls:   decls,
	}

	f, err := helper.AstToDst(file)
	if err != nil {
		res <- config2.Builder{
			Err: err,
		}
		return
	}
	res <- config2.Builder{
		DstFile:  f,
		Pathname: "internal/config/database.go",
	}
}
