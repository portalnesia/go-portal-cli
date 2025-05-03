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
	"go.portalnesia.com/portal-cli/internal/config"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/utils"
	"path"
	"strings"
	"sync"
)

func (c *initType) addStatic(file string, wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating %s.go\n", file)

	src, _ := c.app.DataEmbed.ReadFile(fmt.Sprintf("data/golang/%s.txt", file))
	srcStr := string(src)
	srcStr = strings.ReplaceAll(srcStr, "CLI_NAME", path.Base(c.cfg.Module))
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME_UCWORDS", utils.Ucwords(c.cfg.Module))
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME", c.cfg.Module)
	src = []byte(srcStr)

	res <- config2.Builder{
		Static:   src,
		Pathname: fmt.Sprintf("%s.go", file),
	}
}

func (c *initType) copyStatic(app *config.App, src, dst string, wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating %s\n", dst)

	byt, _ := app.DataEmbed.ReadFile(fmt.Sprintf("data/%s", src))

	res <- config2.Builder{
		Static:   byt,
		Pathname: fmt.Sprintf("%s", dst),
	}
}
