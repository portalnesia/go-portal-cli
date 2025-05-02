/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package sserver

import (
	"go.portalnesia.com/portal-cli/internal/config"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"sync"
)

type Server struct {
	app *config.App
	cfg *config2.InitConfig
}

func Init(parentWg *sync.WaitGroup, app *config.App, cfg *config2.InitConfig, resp chan []config2.Builder) {
	defer parentWg.Done()

	s := &Server{
		app: app,
		cfg: cfg,
	}

	i := 6

	var (
		allFiles = make([]config2.Builder, 0)
		respChan = make(chan config2.Builder, i)
		wg       = &sync.WaitGroup{}
	)
	wg.Add(i)

	go s.initConfigApp(wg, respChan)

	files := []string{
		"internal/server/config/response",
		"internal/server/routes/routes",
		"internal/server/middleware/middleware",
		"internal/server/decoder",
		"internal/server/fiber",
	}
	for _, f := range files {
		go s.addStatic(f, wg, respChan)
	}

	wg.Wait()
	close(respChan)

	for res := range respChan {
		allFiles = append(allFiles, res)
	}

	resp <- allFiles
}
