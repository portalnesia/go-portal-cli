/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package spkg

import (
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"strings"
	"sync"
)

func (p *Pkg) initPkgMigration(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	_, _ = color.New(color.FgBlue).Printf("Generating pkg/migration/migration.go\n")

	src, _ := p.app.DataEmbed.ReadFile("data/golang/pkg/migration.txt")
	srcStr := string(src)
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME", p.cfg.Module)
	src = []byte(srcStr)

	res <- config2.Builder{
		Static:   src,
		Pathname: "pkg/migration/migration.go",
	}

	_, _ = color.New(color.FgBlue).Printf("Generating pkg/migration/README.md\n")

	src, _ = p.app.DataEmbed.ReadFile("data/golang/pkg/migration_readme.txt")

	res <- config2.Builder{
		Static:   src,
		Pathname: "pkg/migration/README.md",
	}
}
