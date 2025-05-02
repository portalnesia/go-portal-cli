/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package sstatic

import (
	"fmt"
	"github.com/fatih/color"
	"go.portalnesia.com/portal-cli/internal/config"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"strings"
	"sync"
)

func Init(parentWg *sync.WaitGroup, app *config.App, cfg *config2.InitConfig, resp chan []config2.Builder) {
	defer parentWg.Done()

	i := 10

	var (
		allFiles = make([]config2.Builder, 0)
		respChan = make(chan config2.Builder, i)
		wg       = &sync.WaitGroup{}
	)
	wg.Add(i)

	files := []string{
		"internal/cerror/exception",
		"internal/cerror/notfound",
		"internal/cerror/server",
		"internal/cerror/parameter",
		"internal/context/context",
		"internal/request/request",
		"main",
	}

	for _, f := range files {
		go addStatic(app, cfg, f, wg, respChan)
	}

	copyFiles := [][]string{
		{
			"favicon.ico",
			"public/favicon.ico",
		},
		{
			"DELETE.txt",
			"data/DELETE.txt",
		},
		{
			"DELETE.txt",
			"migrations/DELETE.txt",
		},
	}

	for _, f := range copyFiles {
		go copyStatic(app, f[0], f[1], wg, respChan)
	}

	wg.Wait()
	close(respChan)

	for res := range respChan {
		allFiles = append(allFiles, res)
	}

	resp <- allFiles
}

func addStatic(app *config.App, cfg *config2.InitConfig, file string, wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating %s.go\n", file)

	src, _ := app.DataEmbed.ReadFile(fmt.Sprintf("data/golang/%s.txt", file))
	srcStr := string(src)
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME", cfg.Module)
	src = []byte(srcStr)

	res <- config2.Builder{
		Static:   src,
		Pathname: fmt.Sprintf("%s.go", file),
	}
}

func copyStatic(app *config.App, src, dst string, wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating %s\n", dst)

	byt, _ := app.DataEmbed.ReadFile(fmt.Sprintf("data/%s", src))

	res <- config2.Builder{
		Static:   byt,
		Pathname: fmt.Sprintf("%s", dst),
	}
}
