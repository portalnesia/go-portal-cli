/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package ginit

import (
	"go/ast"
	"go/token"
	"sync"

	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
)

func (c *initType) initConfigDatabase(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/database.go\n")

	pkgImport := []string{
		`"context"`,
		`"database/sql"`,
		`"fmt"`,
		`"time"`,

		`sqlDriver "github.com/go-sql-driver/mysql"`,
		`"github.com/spf13/viper"`,
		`"github.com/uptrace/bun"`,
		`"github.com/uptrace/bun/dialect/mysqldialect"`,
	}
	decls := make([]ast.Decl, 0)

	decls = append(decls, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("SqlLogger"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{},
				},
			},
		},
	})
	decls = append(decls, &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("h")},
					Type:  &ast.StarExpr{X: ast.NewIdent("SqlLogger")},
				},
			},
		},
		Name: ast.NewIdent("BeforeQuery"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type:  &ast.SelectorExpr{X: ast.NewIdent("context"), Sel: ast.NewIdent("Context")},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("event")}, Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: ast.NewIdent("bun"), Sel: ast.NewIdent("QueryEvent"),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.SelectorExpr{X: ast.NewIdent("context"), Sel: ast.NewIdent("Context")}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{Results: []ast.Expr{ast.NewIdent("ctx")}},
			},
		},
	})
	decls = append(decls, &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("h")},
					Type:  &ast.StarExpr{X: ast.NewIdent("SqlLogger")},
				},
			},
		},
		Name: ast.NewIdent("AfterQuery"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type:  &ast.SelectorExpr{X: ast.NewIdent("context"), Sel: ast.NewIdent("Context")},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("event")}, Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: ast.NewIdent("bun"), Sel: ast.NewIdent("QueryEvent"),
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// query := event.Query
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("query")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{&ast.SelectorExpr{X: ast.NewIdent("event"), Sel: ast.NewIdent("Query")}},
				},
				// dur := time.Since(event.StartTime)
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("dur")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: ast.NewIdent("time"), Sel: ast.NewIdent("Since")},
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X: ast.NewIdent("event"), Sel: ast.NewIdent("StartTime"),
								},
							},
						},
					},
				},
				// if event.Err != nil { ... } else if ...
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.SelectorExpr{X: ast.NewIdent("event"), Sel: ast.NewIdent("Err")},
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
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: ast.NewIdent("log"), Sel: ast.NewIdent("Error"),
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: ast.NewIdent("event"), Sel: ast.NewIdent("Err"),
														}, helper.StrLit("db"),
													},
												},
												Sel: ast.NewIdent("Dur"),
											},
											Args: []ast.Expr{helper.StrLit("duration"), ast.NewIdent("dur")},
										},
										Sel: ast.NewIdent("Msg"),
									},
									Args: []ast.Expr{ast.NewIdent("query")},
								},
							},
						},
					},
					Else: &ast.IfStmt{
						Cond: &ast.CallExpr{
							Fun:  &ast.SelectorExpr{X: ast.NewIdent("viper"), Sel: ast.NewIdent("GetBool")},
							Args: []ast.Expr{helper.StrLit("config.log_db")},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: ast.NewIdent("log"), Sel: ast.NewIdent("Info"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: `"db"`,
															},
														},
													},
													Sel: ast.NewIdent("Dur"),
												},
												Args: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"duration"`,
													}, ast.NewIdent("dur"),
												},
											},
											Sel: ast.NewIdent("Msg"),
										},
										Args: []ast.Expr{ast.NewIdent("query")},
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
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("dbChan")},
						Type: &ast.ChanType{
							Dir:   ast.SEND | ast.RECV,
							Value: &ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("bun"), Sel: ast.NewIdent("DB")}},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// location, err := time.LoadLocation("Asia/Jakarta")
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("location"), ast.NewIdent("err")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun:  &ast.SelectorExpr{X: ast.NewIdent("time"), Sel: ast.NewIdent("LoadLocation")},
							Args: []ast.Expr{helper.StrLit("Asia/Jakarta")},
						},
					},
				},
				// if err != nil { panic(err) }
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: ast.NewIdent("err"), Op: token.NEQ, Y: ast.NewIdent("nil")},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: ast.NewIdent("panic"), Args: []ast.Expr{ast.NewIdent("err")},
								},
							},
						},
					},
				},
				// mysqlconfig := sqlDriver.Config{ ... }
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("mysqlconfig")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{X: ast.NewIdent("sqlDriver"), Sel: ast.NewIdent("Config")},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key: ast.NewIdent("User"), Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: ast.NewIdent("viper"), Sel: ast.NewIdent("GetString"),
										}, Args: []ast.Expr{helper.StrLit("db.mysql.user")},
									},
								},
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Passwd"), Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: ast.NewIdent("viper"), Sel: ast.NewIdent("GetString"),
										}, Args: []ast.Expr{helper.StrLit("db.mysql.password")},
									},
								},
								&ast.KeyValueExpr{
									Key: ast.NewIdent("DBName"), Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: ast.NewIdent("viper"), Sel: ast.NewIdent("GetString"),
										}, Args: []ast.Expr{helper.StrLit("db.mysql.database")},
									},
								},
								&ast.KeyValueExpr{Key: ast.NewIdent("Net"), Value: helper.StrLit("tcp")},
								&ast.KeyValueExpr{
									Key: ast.NewIdent("Addr"),
									Value: &ast.CallExpr{
										Fun: &ast.SelectorExpr{X: ast.NewIdent("fmt"), Sel: ast.NewIdent("Sprintf")},
										Args: []ast.Expr{
											helper.StrLit("%s:%d"),
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: ast.NewIdent("viper"), Sel: ast.NewIdent("GetString"),
												}, Args: []ast.Expr{helper.StrLit("db.mysql.host")},
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: ast.NewIdent("viper"), Sel: ast.NewIdent("GetInt"),
												}, Args: []ast.Expr{helper.StrLit("db.mysql.port")},
											},
										},
									},
								},
								&ast.KeyValueExpr{Key: ast.NewIdent("ParseTime"), Value: ast.NewIdent("true")},
								&ast.KeyValueExpr{Key: ast.NewIdent("Loc"), Value: ast.NewIdent("location")},
								&ast.KeyValueExpr{
									Key: ast.NewIdent("AllowNativePasswords"), Value: ast.NewIdent("true"),
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
								X: ast.NewIdent("mysqlconfig"), Sel: ast.NewIdent("FormatDSN"),
							},
						},
					},
				},
				// sqldb, err := sql.Open("mysql", dsn)
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("sqldb"), ast.NewIdent("err")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun:  &ast.SelectorExpr{X: ast.NewIdent("sql"), Sel: ast.NewIdent("Open")},
							Args: []ast.Expr{helper.StrLit("mysql"), ast.NewIdent("dsn")},
						},
					},
				},
				// if err != nil { log.Fatal(err, "system").Msg("Failed to open mysql") }
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: ast.NewIdent("err"), Op: token.NEQ, Y: ast.NewIdent("nil")},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: ast.NewIdent("log"), Sel: ast.NewIdent("Fatal"),
											}, Args: []ast.Expr{ast.NewIdent("err"), helper.StrLit("system")},
										},
										Sel: ast.NewIdent("Msg"),
									},
									Args: []ast.Expr{helper.StrLit("Failed to open mysql")},
								},
							},
						},
					},
				},
				// db := bun.NewDB(sqldb, mysqldialect.New())
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("db")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: ast.NewIdent("bun"), Sel: ast.NewIdent("NewDB")},
							Args: []ast.Expr{
								ast.NewIdent("sqldb"),
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: ast.NewIdent("mysqldialect"), Sel: ast.NewIdent("New"),
									},
								},
							},
						},
					},
				},
				// db.AddQueryHook(&SqlLogger{})
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{X: ast.NewIdent("db"), Sel: ast.NewIdent("AddQueryHook")},
						Args: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND, X: &ast.CompositeLit{Type: ast.NewIdent("SqlLogger")},
							},
						},
					},
				},
				// db.SetMaxIdleConns(10) ... dst
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: ast.NewIdent("db"), Sel: ast.NewIdent("SetMaxIdleConns"),
						}, Args: []ast.Expr{helper.IntLit(10)},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: ast.NewIdent("db"), Sel: ast.NewIdent("SetMaxOpenConns"),
						}, Args: []ast.Expr{helper.IntLit(100)},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun:  &ast.SelectorExpr{X: ast.NewIdent("db"), Sel: ast.NewIdent("SetConnMaxLifetime")},
						Args: []ast.Expr{&ast.SelectorExpr{X: ast.NewIdent("time"), Sel: ast.NewIdent("Hour")}},
					},
				},
				// if err = db.Ping(); err != nil { ... }
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{ast.NewIdent("err")}, Tok: token.ASSIGN, Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: ast.NewIdent("db"), Sel: ast.NewIdent("Ping"),
								},
							},
						},
					},
					Cond: &ast.BinaryExpr{X: ast.NewIdent("err"), Op: token.NEQ, Y: ast.NewIdent("nil")},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: ast.NewIdent("log"), Sel: ast.NewIdent("Fatal"),
											}, Args: []ast.Expr{ast.NewIdent("err"), helper.StrLit("system")},
										},
										Sel: ast.NewIdent("Msg"),
									},
									Args: []ast.Expr{helper.StrLit("Failed to ping mysql")},
								},
							},
						},
					},
				},
				// dbChan <- db
				&ast.SendStmt{Chan: ast.NewIdent("dbChan"), Value: ast.NewIdent("db")},
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
