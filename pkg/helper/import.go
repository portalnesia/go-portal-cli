/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package helper

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"go/ast"
	"go/token"
)

type Imports struct {
	imports   []*ast.ImportSpec
	importDst []*dst.ImportSpec
}

func (i *Imports) ImportSpec() []*ast.ImportSpec {
	return i.imports
}

func (i *Imports) ImportSpecDst() []*dst.ImportSpec {
	return i.importDst
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

func (i *Imports) GenDeclDst() *dst.GenDecl {
	if i.imports == nil || len(i.imports) <= 0 {
		return nil
	}

	var tmp []dst.Spec
	for _, im := range i.importDst {
		tmp = append(tmp, im)
	}

	return &dst.GenDecl{
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
	importDst := make([]*dst.ImportSpec, 0)

	for _, p := range pkg {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{Kind: token.STRING, Value: p},
		})
		importDst = append(importDst, &dst.ImportSpec{
			Path: &dst.BasicLit{Kind: token.STRING, Value: p},
		})
	}

	return Imports{
		imports:   imports,
		importDst: importDst,
	}
}

func fixManualASTPackagePos(fset *token.FileSet, file *ast.File) {
	if file.Pos() == token.NoPos {
		// Daftarkan posisi token buatan untuk file kamu
		tf := fset.AddFile("manual.go", -1, 1000) // nama, base, size
		file.Package = token.Pos(tf.Base())       // Set posisi Package agar bukan NoPos
	}
}

func AstToDst(file *ast.File) (*dst.File, error) {
	fset := token.NewFileSet()
	tokFile := fset.AddFile("manual.go", fset.Base(), 1000)
	pos := tokFile.Pos(0)
	file.Package = pos
	dec := decorator.NewDecorator(fset)
	return dec.DecorateFile(file)
}
