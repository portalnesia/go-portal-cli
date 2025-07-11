/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package ginit

import (
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"strings"
	"sync"
)

func (c *initType) initConfigLog(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/log.go\n")

	src, _ := c.app.DataEmbed.ReadFile("data/golang/internal/config/log.txt")
	srcStr := string(src)
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME", c.cfg.Module)

	src = []byte(srcStr)
	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/config/log.go",
	}
}
