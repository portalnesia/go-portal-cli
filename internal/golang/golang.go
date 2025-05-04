/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package b_golang

import (
	"bytes"
	"fmt"
	"go.portalnesia.com/portal-cli/internal/config"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
	"go/format"
	"go/printer"
	"go/token"
	"os"
	"path"
)

type Golang struct {
	app *config.App
}

func New(app *config.App) *Golang {
	return &Golang{
		app: app,
	}
}

func (g *Golang) Build(builder []config2.Builder) error {
	for _, b := range builder {
		if b.Err != nil {
			return b.Err
		}

		pathname := g.app.Dir(b.Pathname)

		var (
			src []byte
			err error
		)

		if b.File != nil {
			fset := token.NewFileSet()
			// Render AST ke source
			var buf bytes.Buffer

			buf.WriteString(helper.GenCopyright(b.Comment...))

			printerConfig := &printer.Config{
				Mode:     printer.UseSpaces | printer.TabIndent,
				Tabwidth: 4,
			}

			if err = printerConfig.Fprint(&buf, fset, b.File); err != nil {
				return fmt.Errorf("failed to print file: %s", err.Error())
			}
			src, err = format.Source(buf.Bytes())
			if err != nil {
				return fmt.Errorf("failed to format source: %s", err.Error())
			}
		} else {
			src = b.Static
		}

		dir := path.Dir(pathname)
		// Simpan ke file
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %s", err.Error())
		}

		if err = os.WriteFile(pathname, src, 0644); err != nil {
			return fmt.Errorf("failed to write file: %s", err.Error())
		}
	}
	return nil
}
