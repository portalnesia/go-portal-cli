/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package config

import (
	"github.com/dave/dst"
	"go/ast"
)

type InitConfig struct {
	Module string

	Redis      bool
	Firebase   bool
	Handlebars bool
}

type AddServiceConfig struct {
	Module  string
	Name    string
	Path    string
	Version string
}

type AddEndpointConfig struct {
	Module      string
	ServiceName string // ex: user
	Name        string // GetUser
	Path        string // endpoint path, include version
	Method      string // GET, POST, PUT, PATCH, DELETE
}

type AddRepositoryConfig struct {
	Module  string
	Name    string
	NoModel bool
}

type Builder struct {
	Comment        []string
	File           *ast.File
	Pathname       string
	Static         []byte
	Err            error
	DstFile        *dst.File
	WithoutComment bool
}
