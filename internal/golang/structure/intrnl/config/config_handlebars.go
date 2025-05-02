/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package s_config

import (
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"strings"
	"sync"
)

func (c *Config) initConfigHandlebars(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/handlebars.go\n")

	src, _ := c.app.DataEmbed.ReadFile("data/golang/internal/config/handlebars.txt")
	src = []byte(strings.ReplaceAll(string(src), "app_name", c.cfg.Module))

	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/config/handlebars.go",
	}
}
