/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package ginit

import (
	"go.portalnesia.com/portal-cli/internal/config"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"sync"
)

type initType struct {
	app *config.App
	cfg *config2.InitConfig
}

func Init(parentWg *sync.WaitGroup, app *config.App, cfg *config2.InitConfig, resp chan []config2.Builder) {
	defer parentWg.Done()

	c := &initType{
		app: app,
		cfg: cfg,
	}

	i := 29
	if cfg.Redis {
		i += 1
	}
	if cfg.Firebase {
		i += 1
	}
	if cfg.Handlebars {
		i += 1
	}

	var (
		allFiles = make([]config2.Builder, 0)
		respChan = make(chan config2.Builder, i)
		wg       = &sync.WaitGroup{}
	)
	wg.Add(i)

	if cfg.Redis {
		go c.initConfigRedis(wg, respChan)
	}
	if cfg.Firebase {
		go c.initConfigFirebase(wg, respChan)
	}
	if cfg.Handlebars {
		go c.initConfigHandlebars(wg, respChan)
	}

	go c.initConfigVersionGen(wg, respChan)
	go c.initConfigVersion(wg, respChan)
	go c.initMainVerionGen(wg, respChan)
	go c.initConfigDatabase(wg, respChan)
	go c.initConfigLog(wg, respChan)
	go c.initConfigApp(wg, respChan)
	go c.initConfigValidator(wg, respChan)

	go c.initServerConfigApp(wg, respChan)
	go c.initCmdStart(wg, respChan)

	files := []string{
		"internal/server/config/response",
		"internal/server/routes/routes",
		"internal/server/middleware/middleware",
		"internal/server/decoder",
		"internal/server/fiber",

		"internal/cerror/exception",
		"internal/cerror/notfound",
		"internal/cerror/server",
		"internal/cerror/parameter",
		"internal/context/context",
		"internal/request/request",
		"pkg/helper/main",
		"pkg/migration/migration",
		"cmd/completion",
		"cmd/app",
		"main",
	}
	for _, f := range files {
		go c.addStatic(f, wg, respChan)
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
		{
			"golang/pkg/migration/README.txt",
			"pkg/migration/README.md",
		},
	}

	for _, f := range copyFiles {
		go c.copyStatic(app, f[0], f[1], wg, respChan)
	}

	wg.Wait()
	close(respChan)

	for res := range respChan {
		allFiles = append(allFiles, res)
	}

	resp <- allFiles
}
