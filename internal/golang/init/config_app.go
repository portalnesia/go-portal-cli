/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package ginit

import (
	"fmt"
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
	"go/ast"
	"go/token"
	"sync"
)

func (c *initType) initConfigApp(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/app.go\n")
	pkgImport := []string{
		`"embed"`,
		`"github.com/rs/zerolog"`,
		`"github.com/spf13/viper"`,
		`"github.com/subosito/gotenv"`,
		`pncrypto "go.portalnesia.com/crypto"`,
		`"gorm.io/gorm"`,
		`"os"`,
		fmt.Sprintf(`"%s/internal/repository"`, c.cfg.Module),
		fmt.Sprintf(`iface "%s/internal/interface"`, c.cfg.Module),
		`"strings"`,
		`"time"`,
	}
	decls := make([]ast.Decl, 0)

	// TYPE
	decls = c.appInitType(decls)

	// NEW
	decls = c.appInitNew(decls)

	// CLOSE
	decls = c.appInitClose(decls)

	if c.cfg.Redis {
		pkgImport = append(pkgImport,
			`"github.com/gofiber/fiber/v2"`,
			`fiberredis "github.com/gofiber/storage/redis/v3"`,
			`"github.com/gofiber/fiber/v2/middleware/session"`,
			`"github.com/redis/go-redis/v9"`,
		)
	}

	imports := helper.GenImport(pkgImport...)
	decls = append([]ast.Decl{imports.GenDecl()}, decls...)
	file := &ast.File{
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
		Pathname: "internal/config/app.go",
	}
}

func (c *initType) appInitNew(decls []ast.Decl) []ast.Decl {
	// VAR
	specs := []ast.Spec{
		&ast.ValueSpec{
			Names: []*ast.Ident{ast.NewIdent("chDb")},
			Values: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("make"),
					Args: []ast.Expr{
						&ast.ChanType{
							Dir: ast.SEND | ast.RECV,
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
		&ast.ValueSpec{
			Names: []*ast.Ident{ast.NewIdent("err")},
			Type:  ast.NewIdent("error"),
		},
	}
	if c.cfg.Redis {
		specs = append(specs,
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("storage")},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("fiberredis"),
						Sel: ast.NewIdent("Storage"),
					},
				},
			},
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("redisClient")},
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("redis"),
					Sel: ast.NewIdent("UniversalClient"),
				},
			},
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("sessionStore")},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("session"),
						Sel: ast.NewIdent("Store"),
					},
				},
			},
		)
	}

	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok:   token.VAR,
					Specs: specs,
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("embedFs")},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("config"),
						Sel: ast.NewIdent("Embed"),
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "env",
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.Ident{
							Name: "envImpl",
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "build",
								},
								Value: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "config",
									},
									Sel: &ast.Ident{
										Name: "Build",
									},
								},
							},
						},
					},
				},
			},
			helper.BodyListNewLines(),
		},
	}

	// initHandlebars()
	if c.cfg.Handlebars {
		body.List = append(body.List, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: ast.NewIdent("initHandlebars"),
			},
		}, helper.BodyListNewLines())
	}

	// logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Logger()
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `// Init log`,
		}, // dummy expression, tidak valid
	}, &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("logger")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("zerolog"),
													Sel: ast.NewIdent("New"),
												},
												Args: []ast.Expr{
													&ast.CompositeLit{
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("zerolog"),
															Sel: ast.NewIdent("ConsoleWriter"),
														},
														Elts: []ast.Expr{
															&ast.KeyValueExpr{
																Key: ast.NewIdent("Out"),
																Value: &ast.SelectorExpr{
																	X:   ast.NewIdent("os"),
																	Sel: ast.NewIdent("Stderr"),
																},
															},
															&ast.KeyValueExpr{
																Key: ast.NewIdent("TimeFormat"),
																Value: &ast.SelectorExpr{
																	X:   ast.NewIdent("time"),
																	Sel: ast.NewIdent("RFC3339"),
																},
															},
														},
													},
												},
											},
											Sel: ast.NewIdent("With"),
										},
									},
									Sel: ast.NewIdent("Timestamp"),
								},
							},
							Sel: ast.NewIdent("Caller"),
						},
					},
					Sel: ast.NewIdent("Logger"),
				},
			},
		},
	})
	// if !env.IsServer() { logger = logger.Level(zerolog.InfoLevel) }
	body.List = append(body.List, &ast.IfStmt{
		Cond: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("env"),
				Sel: ast.NewIdent("IsServer"),
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("logger"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("logger"),
								Sel: ast.NewIdent("Level"),
							},
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("zerolog"),
									Sel: ast.NewIdent("InfoLevel"),
								},
							},
						},
					},
				},
			},
		},
	})
	// log = Logger{logger}
	body.List = append(body.List, &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("log")},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CompositeLit{
				Type: ast.NewIdent("Logger"),
				Elts: []ast.Expr{
					ast.NewIdent("logger"),
				},
			},
		},
	})
	// log.Info("system").Str("version", version.String()).Msg("Initializing application configuration...")
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("log"),
								Sel: ast.NewIdent("Info"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"system"`,
								},
							},
						},
						Sel: ast.NewIdent("Str"),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"version"`,
						},
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("version"),
								Sel: ast.NewIdent("String"),
							},
						},
					},
				},
				Sel: ast.NewIdent("Msg"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"Initializing application configuration..."`,
				},
			},
		},
	}, helper.BodyListNewLines())

	// _ = gotenv.Load()
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `// Init config`,
		}, // dummy expression, tidak valid
	}, &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("_")},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("gotenv"),
					Sel: ast.NewIdent("Load"),
				},
			},
		},
	})
	// replacer := strings.NewReplacer(".", "_")
	body.List = append(body.List, &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("replacer")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("strings"),
					Sel: ast.NewIdent("NewReplacer"),
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: `"."`,
					},
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: `"_"`,
					},
				},
			},
		},
	})
	// viper.SetEnvKeyReplacer(replacer)
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("viper"),
				Sel: ast.NewIdent("SetEnvKeyReplacer"),
			},
			Args: []ast.Expr{ast.NewIdent("replacer")},
		},
	})
	// viper.SetEnvPrefix("config")
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("viper"),
				Sel: ast.NewIdent("SetEnvPrefix"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, c.cfg.Module),
				},
			},
		},
	})
	// viper.SetConfigName("config")
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("viper"),
				Sel: ast.NewIdent("SetConfigName"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"config"`,
				},
			},
		},
	})
	// viper.SetConfigType("json")
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("viper"),
				Sel: ast.NewIdent("SetConfigType"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"json"`,
				},
			},
		},
	})
	// // viper.AddConfigPath("docker-data/json")
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  token.COMMENT,
			Value: "// viper.AddConfigPath(\"docker-data/json\")",
		},
	})
	// viper.AddConfigPath(".")
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("viper"),
				Sel: ast.NewIdent("AddConfigPath"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"."`,
				},
			},
		},
	})
	// viper.SetDefault("env", "development")
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("viper"),
				Sel: ast.NewIdent("SetDefault"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"env"`,
				},
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"development"`,
				},
			},
		},
	})
	// viper.AutomaticEnv()
	body.List = append(body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("viper"),
				Sel: ast.NewIdent("AutomaticEnv"),
			},
		},
	})
	// err = viper.ReadInConfig()
	body.List = append(body.List, &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("err")},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("viper"),
					Sel: ast.NewIdent("ReadInConfig"),
				},
			},
		},
	})
	// if err != nil { log.Fatal(err, "viper").Msg("Failed to load config file") }
	body.List = append(body.List, &ast.IfStmt{
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
										Value: `"viper"`,
									},
								},
							},
							Sel: ast.NewIdent("Msg"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"Failed to load config file"`,
							},
						},
					},
				},
			},
		},
	}, helper.BodyListNewLines())

	// CRYPTO
	body.List = append(body.List,
		// // Init crypto
		&ast.ExprStmt{
			X: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `// Init crypto`,
			}, // dummy expression, tidak valid
		},
		// crypto := pncrypto.New(viper.GetString("secret.crypto"))
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("crypto")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("pncrypto"),
						Sel: ast.NewIdent("New"),
					},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("viper"),
								Sel: ast.NewIdent("GetString"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"secret.crypto"`,
								},
							},
						},
					},
				},
			},
		},
		helper.BodyListNewLines(),
	)

	// DATABASE
	body.List = append(body.List,
		// // Init database
		&ast.ExprStmt{
			X: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `// Init database`,
			}, // dummy expression, tidak valid
		},
		&ast.IfStmt{
			Cond: &ast.SelectorExpr{
				X:   ast.NewIdent("config"),
				Sel: ast.NewIdent("DB"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					// log.Info("system").Msg("Initializing database...")
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("log"),
										Sel: ast.NewIdent("Info"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"system"`,
										},
									},
								},
								Sel: ast.NewIdent("Msg"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"Initializing database..."`,
								},
							},
						},
					},
					// go connectMysql(chDb)
					&ast.GoStmt{
						Call: &ast.CallExpr{
							Fun:  ast.NewIdent("connectMysql"),
							Args: []ast.Expr{ast.NewIdent("chDb")},
						},
					},
				},
			},
		},
		helper.BodyListNewLines(),
	)

	if c.cfg.Redis {
		body.List = append(body.List,
			&ast.ExprStmt{
				X: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `// Init redis`,
				}, // dummy expression, tidak valid
			},
			&ast.IfStmt{
				Cond: &ast.SelectorExpr{
					X:   ast.NewIdent("config"),
					Sel: ast.NewIdent("Redis"),
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						// log.Info("system").Msg("Initializing redis...")
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("log"),
											Sel: ast.NewIdent("Info"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `"system"`,
											},
										},
									},
									Sel: ast.NewIdent("Msg"),
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"Initializing redis..."`,
									},
								},
							},
						},
						// redisName := "APP_NAME"
						&ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("redisName")},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(`"%s"`, c.cfg.Module),
								},
							},
						},
						// if isProduction { redisName += "-" + Version.String() } else { redisName += "-local" }
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("env"),
									Sel: ast.NewIdent("IsProduction"),
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											ast.NewIdent("redisName"),
										},
										Tok: token.ADD_ASSIGN,
										Rhs: []ast.Expr{
											&ast.BinaryExpr{
												X: &ast.BasicLit{
													Kind:  token.STRING,
													Value: `"-"`,
												},
												Op: token.ADD,
												Y: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("version"),
														Sel: ast.NewIdent("String"),
													},
												},
											},
										},
									},
								},
							},
							Else: &ast.IfStmt{
								Cond: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("env"),
										Sel: ast.NewIdent("IsLocal"),
									},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.AssignStmt{
											Lhs: []ast.Expr{
												ast.NewIdent("redisName"),
											},
											Tok: token.ADD_ASSIGN,
											Rhs: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: `"-local"`,
												},
											},
										},
									},
								},
								Else: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.AssignStmt{
											Lhs: []ast.Expr{
												ast.NewIdent("redisName"),
											},
											Tok: token.ADD_ASSIGN,
											Rhs: []ast.Expr{
												&ast.BinaryExpr{
													X: &ast.BasicLit{
														Kind:  token.STRING,
														Value: `"-"`,
													},
													Op: token.ADD,
													Y: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("env"),
															Sel: ast.NewIdent("EnvShortString"),
														},
													},
												},
											},
										},
									},
								},
							},
						},
						// storage = initFiberStorage(redisName)
						&ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("storage")},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun:  ast.NewIdent("initFiberStorage"),
									Args: []ast.Expr{ast.NewIdent("redisName")},
								},
							},
						},
						// redisClient = storage.Conn()
						&ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("redisClient")},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("storage"),
										Sel: ast.NewIdent("Conn"),
									},
								},
							},
						},
						// sessionStore = initSession(storage)
						&ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("sessionStore")},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun:  ast.NewIdent("initSession"),
									Args: []ast.Expr{ast.NewIdent("storage")},
								},
							},
						},
					},
				},
			},
			helper.BodyListNewLines(),
		)
	}

	appBodyList := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: "apps",
			},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: &ast.Ident{
						Name: "app",
					},
					Elts: []ast.Expr{
						&ast.KeyValueExpr{
							Key: &ast.Ident{
								Name: "config",
							},
							Value: &ast.Ident{
								Name: "config",
							},
						},
						&ast.KeyValueExpr{
							Key: &ast.Ident{
								Name: "env",
							},
							Value: &ast.Ident{
								Name: "env",
							},
						},
						&ast.KeyValueExpr{
							Key: &ast.Ident{
								Name: "embed",
							},
							Value: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "config",
								},
								Sel: &ast.Ident{
									Name: "Embed",
								},
							},
						},
						&ast.KeyValueExpr{
							Key: &ast.Ident{
								Name: "crypto",
							},
							Value: &ast.Ident{
								Name: "crypto",
							},
						},
						&ast.KeyValueExpr{
							Key: &ast.Ident{
								Name: "log",
							},
							Value: &ast.UnaryExpr{
								Op: token.AND,
								X: &ast.Ident{
									Name: "log",
								},
							},
						},
					},
				},
			},
		},
	}
	if c.cfg.Redis {
		appBodyList.Rhs[0].(*ast.UnaryExpr).X.(*ast.CompositeLit).Elts = append(appBodyList.Rhs[0].(*ast.UnaryExpr).X.(*ast.CompositeLit).Elts, &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "redis",
			},
			Value: &ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: &ast.Ident{
						Name: "redisImpl",
					},
					Elts: []ast.Expr{
						&ast.Ident{
							Name: "redisClient",
						},
					},
				},
			},
		}, &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "fiberStorage",
			},
			Value: &ast.Ident{
				Name: "storage",
			},
		}, &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "sessionStore",
			},
			Value: &ast.Ident{
				Name: "sessionStore",
			},
		})
	}
	body.List = append(body.List, appBodyList, helper.BodyListNewLines())

	body.List = append(body.List, &ast.IfStmt{
		Cond: &ast.SelectorExpr{
			X:   ast.NewIdent("config"),
			Sel: ast.NewIdent("DB"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("apps"),
							Sel: ast.NewIdent("db"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.ARROW,
							X:  ast.NewIdent("chDb"),
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("apps"),
							Sel: ast.NewIdent("repo"),
						},
					},
					Tok: token.ASSIGN, // =
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("repository"),
								Sel: ast.NewIdent("NewRepository"),
							},
							Args: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "repository"},
										Sel: &ast.Ident{Name: "NewRepositoryConfig"},
									},
									Elts: []ast.Expr{
										// DB: apps.db
										&ast.KeyValueExpr{
											Key: &ast.Ident{Name: "DB"},
											Value: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "apps"},
												Sel: &ast.Ident{Name: "db"},
											},
										},
										// Env: env
										&ast.KeyValueExpr{
											Key:   &ast.Ident{Name: "Env"},
											Value: &ast.Ident{Name: "env"},
										},
										// Log: &log
										&ast.KeyValueExpr{
											Key: &ast.Ident{Name: "Log"},
											Value: &ast.UnaryExpr{
												Op: token.AND,
												X:  &ast.Ident{Name: "log"},
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
	}, helper.BodyListNewLines())

	body.List = append(body.List, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.Ident{
				Name: "apps",
			},
		},
	})

	decls = append(decls, &ast.FuncDecl{
		Name: ast.NewIdent("New"),
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("config")},
						Type:  ast.NewIdent("Config"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("App"),
					},
				},
			},
		},
		Body: body,
	})

	return decls
}

func (c *initType) appInitClose(decls []ast.Decl) []ast.Decl {
	bodyList := &ast.BlockStmt{
		List: []ast.Stmt{
			// a.Log.Info("system").Msg("Closing application...")
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("a"),
									Sel: ast.NewIdent("log"),
								},
								Sel: ast.NewIdent("Info"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"system"`,
								},
							},
						},
						Sel: ast.NewIdent("Msg"),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"Closing application..."`,
						},
					},
				},
			},
			// var err error
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{
								{Name: "err"},
							},
							Type: &ast.Ident{Name: "error"},
						},
					},
				},
			},
			helper.BodyListNewLines(),
		},
	}

	// redis
	if c.cfg.Redis {
		bodyList.List = append(bodyList.List, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("a"),
					Sel: ast.NewIdent("fiberStorage"),
				},
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					// a.Log.Info("system").Msg("Closing redis...")
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("a"),
											Sel: ast.NewIdent("log"),
										},
										Sel: ast.NewIdent("Info"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"system"`,
										},
									},
								},
								Sel: ast.NewIdent("Msg"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"Closing redis..."`,
								},
							},
						},
					},
					// if err = a.FiberStorage.Close(); err != nil {
					&ast.IfStmt{
						Init: &ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("err")},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("a"),
											Sel: ast.NewIdent("fiberStorage"),
										},
										Sel: ast.NewIdent("Close"),
									},
								},
							},
						},
						Cond: &ast.BinaryExpr{
							X:  ast.NewIdent("err"),
							Op: token.NEQ,
							Y:  ast.NewIdent("nil"),
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								// a.Log.Error(err, "redis").Msg("Error closing redis connection")
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("a"),
														Sel: ast.NewIdent("log"),
													},
													Sel: ast.NewIdent("Error"),
												},
												Args: []ast.Expr{
													ast.NewIdent("err"),
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"system"`,
													},
												},
											},
											Sel: ast.NewIdent("Msg"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `"Error closing redis connection"`,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, helper.BodyListNewLines())
	}

	// DB
	bodyList.List = append(bodyList.List, &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("a"),
				Sel: ast.NewIdent("db"),
			},
			Op: token.NEQ,
			Y:  ast.NewIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// a.Log.Info("system").Msg("Closing postgresql...")
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("a"),
										Sel: ast.NewIdent("log"),
									},
									Sel: ast.NewIdent("Info"),
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"system"`,
									},
								},
							},
							Sel: ast.NewIdent("Msg"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"Closing database..."`,
							},
						},
					},
				},
				// if sqlDB, _ := a.DB.DB(); sqlDB != nil {
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("sqlDB"),
							ast.NewIdent("_"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("a"),
										Sel: ast.NewIdent("db"),
									},
									Sel: ast.NewIdent("DB"),
								},
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("sqlDB"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							// if err = sqlDB.Close(); err != nil {
							&ast.IfStmt{
								Init: &ast.AssignStmt{
									Lhs: []ast.Expr{ast.NewIdent("err")},
									Tok: token.ASSIGN,
									Rhs: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("sqlDB"),
												Sel: ast.NewIdent("Close"),
											},
										},
									},
								},
								Cond: &ast.BinaryExpr{
									X:  ast.NewIdent("err"),
									Op: token.NEQ,
									Y:  ast.NewIdent("nil"),
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										// a.Log.Error(err, "postgre").Msg("Error closing postgre connection")
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("a"),
																Sel: ast.NewIdent("log"),
															},
															Sel: ast.NewIdent("Error"),
														},
														Args: []ast.Expr{
															ast.NewIdent("err"),
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: `"system"`,
															},
														},
													},
													Sel: ast.NewIdent("Msg"),
												},
												Args: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"Error closing database connection"`,
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

	decls = append(decls, &ast.FuncDecl{
		// Nama fungsi
		Name: &ast.Ident{Name: "Close"},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		// Receiver (parameter receiver *App)
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "a"},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "app"},
					},
				},
			},
		},
		// Tipe fungsi (tidak ada parameter dan hasil)
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: &ast.FieldList{},
		},
		// Body dari fungsi
		Body: bodyList,
	})

	return decls
}

func (c *initType) appInitType(decls []ast.Decl) []ast.Decl {
	// embed =
	decls = append(decls, &ast.GenDecl{
		Tok: token.TYPE,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Embed"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("Migration")},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("embed"),
									Sel: ast.NewIdent("FS"),
								},
							},
							{
								Names: []*ast.Ident{ast.NewIdent("Public")},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("embed"),
									Sel: ast.NewIdent("FS"),
								},
							},
							{
								Names: []*ast.Ident{ast.NewIdent("Data")},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("embed"),
									Sel: ast.NewIdent("FS"),
								},
							},
						},
					},
				},
			},
		},
	})
	// var log,embed
	decls = append(decls, &ast.GenDecl{
		Tok: token.VAR,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("env")},
				Type: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "iface"},
					Sel: &ast.Ident{Name: "Env"},
				},
			},
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("log")},
				Type:  ast.NewIdent("Logger"),
			},
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("embedFs")},
				Type:  ast.NewIdent("Embed"),
			},
		},
	})

	// type Config
	typeConfigListField := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent("DB")},
			Type:  ast.NewIdent("bool"),
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Build")},
			Type:  ast.NewIdent("string"),
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Embed")},
			Type:  ast.NewIdent("Embed"),
		},
	}
	if c.cfg.Redis {
		typeConfigListField = append(typeConfigListField, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent("Redis")},
			Type:  ast.NewIdent("bool"),
		})
	}
	decls = append(decls, &ast.GenDecl{
		Tok: token.TYPE,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Config"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: typeConfigListField,
					},
				},
			},
		},
	})

	//interface App
	interfaceAppListField := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent("Env")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X: ast.NewIdent("iface"), Sel: ast.NewIdent("Env"),
							},
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Embed")},
			Type: &ast.FuncType{
				Params:  &ast.FieldList{},
				Results: &ast.FieldList{List: []*ast.Field{{Type: ast.NewIdent("Embed")}}},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Crypto")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X: ast.NewIdent("pncrypto"), Sel: ast.NewIdent("Crypto"),
							},
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Log")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X: ast.NewIdent("iface"), Sel: ast.NewIdent("Logger"),
							},
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("DB")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X: ast.NewIdent("gorm"), Sel: ast.NewIdent("DB"),
								},
							},
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Repository")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X: ast.NewIdent("repository"), Sel: ast.NewIdent("Registry"),
							},
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Close")},
			Type: &ast.FuncType{
				Params:  &ast.FieldList{},
				Results: nil,
			},
		},
	}
	if c.cfg.Redis {
		interfaceAppListField = append(interfaceAppListField,
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("Redis")},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X: ast.NewIdent("iface"), Sel: ast.NewIdent("RedisInterface"),
								},
							},
						},
					},
				},
			},
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("FiberStorage")},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X: ast.NewIdent("fiber"), Sel: ast.NewIdent("Storage"),
								},
							},
						},
					},
				},
			},
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("SessionStore")},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: ast.NewIdent("session"), Sel: ast.NewIdent("Store"),
									},
								},
							},
						},
					},
				},
			},
		)
	}
	decls = append(decls, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("App"),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: interfaceAppListField,
					},
				},
			},
		},
	})

	// type app
	typeAppListField := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent("config")},
			Type:  ast.NewIdent("Config"),
		},
		{
			Names: []*ast.Ident{ast.NewIdent("embed")},
			Type:  ast.NewIdent("Embed"),
		},
		{
			Names: []*ast.Ident{ast.NewIdent("env")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("iface"),
				Sel: ast.NewIdent("Env"),
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("crypto")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("pncrypto"),
				Sel: ast.NewIdent("Crypto"),
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("log")},
			Type:  &ast.StarExpr{X: ast.NewIdent("Logger")},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("db")},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("gorm"),
					Sel: ast.NewIdent("DB"),
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("repo")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("repository"),
				Sel: ast.NewIdent("Registry"),
			},
		},
	}
	if c.cfg.Redis {
		typeAppListField = append(typeAppListField,
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("redis")},
				Type:  &ast.StarExpr{X: ast.NewIdent("redisImpl")},
			},
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("sessionStore")},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("session"),
						Sel: ast.NewIdent("Store"),
					},
				},
			},
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("fiberStorage")},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("fiberredis"),
						Sel: ast.NewIdent("Storage"),
					},
				},
			},
		)
	}

	decls = append(decls, &ast.GenDecl{
		Tok: token.TYPE,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//"}, // Komentar kosong, yang nanti diformat jadi newline
			},
		},
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("app"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: typeAppListField,
					},
				},
			},
		},
	})
	return decls
}
