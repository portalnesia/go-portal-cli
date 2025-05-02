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

type GlobalConfig struct {
	ServerDirectory bool
}

type InitConfig struct {
	Global GlobalConfig

	Module string

	Redis      bool
	Firebase   bool
	Handlebars bool
}

type Builder struct {
	Comment  []string
	File     *ast.File
	Pathname string
	Static   []byte
}
