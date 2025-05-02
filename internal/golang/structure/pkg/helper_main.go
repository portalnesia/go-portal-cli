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
	"sync"
)

func (p *Pkg) initPkgHelper(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()

	_, _ = color.New(color.FgBlue).Printf("Generating pkg/helper/main.go\n")

	src, _ := p.app.DataEmbed.ReadFile("data/golang/pkg/helper.txt")

	res <- config2.Builder{
		Static:   src,
		Pathname: "pkg/helper/main.go",
	}
}
