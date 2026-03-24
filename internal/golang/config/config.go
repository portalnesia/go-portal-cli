/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package config

import (
	"go/ast"
	"regexp"
	"strings"

	"github.com/dave/dst"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type InitConfig struct {
	Module string

	Redis      bool
	Firebase   bool
	Handlebars bool
}

type AddServiceConfig struct {
	Module   string
	Name     string
	Path     string
	PathName string
	Version  string
}

type AddEndpointConfig struct {
	Module          string
	ServiceName     string // ex: user
	ServicePathName string // ex: user
	Name            string // GetUser
	Path            string // endpoint path, include version
	Method          string // GET, POST, PUT, PATCH, DELETE
}

type AddRepositoryConfig struct {
	Module   string
	Name     string
	PathName string
	NoModel  bool
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

func ParseName(name string) (structName, pathName string) {
	// replace symbol to space
	name = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(name, " ")
	name = cases.Title(language.English).String(name)
	structName = strings.ReplaceAll(name, " ", "")
	pathName = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	return
}
