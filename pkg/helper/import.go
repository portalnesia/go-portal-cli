/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package helper

import (
	"go/ast"
	"go/token"
)

type Imports struct {
	imports []*ast.ImportSpec
}

func (i *Imports) ImportSpec() []*ast.ImportSpec {
	return i.imports
}

func (i *Imports) GenDecl() *ast.GenDecl {
	if i.imports == nil || len(i.imports) <= 0 {
		return nil
	}

	var tmp []ast.Spec
	for _, im := range i.imports {
		tmp = append(tmp, im)
	}

	return &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: tmp,
	}
}

func GenImport(
	pkg ...string,
) Imports {
	if len(pkg) == 0 {
		return Imports{}
	}

	imports := make([]*ast.ImportSpec, 0)

	for _, p := range pkg {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{Kind: token.STRING, Value: p},
		})
	}

	return Imports{
		imports: imports,
	}
}
