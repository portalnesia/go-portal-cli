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
	"sync"
)

func (c *Config) initConfigVersionGen(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/version_gen.go\n")

	src, _ := c.app.DataEmbed.ReadFile("data/golang/internal/config/version_gen.txt")

	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/config/version_gen.go",
	}
}

func (c *Config) initConfigVersion(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/version.go\n")

	src, _ := c.app.DataEmbed.ReadFile("data/golang/internal/config/version.txt")

	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/config/version.go",
	}
}

func (c *Config) initMainVerionGen(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating version_gen.go\n")
	src, _ := c.app.DataEmbed.ReadFile("data/golang/version_gen.txt")

	res <- config2.Builder{
		Static:   src,
		Pathname: "version_gen.go",
	}
}
