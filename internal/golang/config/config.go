/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package config

import (
	"go/ast"
)

type InitConfig struct {
	Module string

	Redis      bool
	Firebase   bool
	Handlebars bool
}

type NewServiceConfig struct {
	Module  string
	Name    string
	Path    string
	Version string
}

type Builder struct {
	Comment  []string
	File     *ast.File
	Pathname string
	Static   []byte
	Err      error
}
