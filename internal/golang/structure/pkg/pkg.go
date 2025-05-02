/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package spkg

import (
	"go.portalnesia.com/portal-cli/internal/config"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"sync"
)

type Pkg struct {
	app *config.App
	cfg *config2.InitConfig
}

func Init(parentWg *sync.WaitGroup, app *config.App, cfg *config2.InitConfig, resp chan []config2.Builder) {
	defer parentWg.Done()

	h := &Pkg{
		app: app,
		cfg: cfg,
	}

	i := 2

	var (
		allFiles = make([]config2.Builder, 0)
		respChan = make(chan config2.Builder, i+1)
		wg       = &sync.WaitGroup{}
	)
	wg.Add(i)

	go h.initPkgHelper(wg, respChan)
	go h.initPkgMigration(wg, respChan)

	wg.Wait()
	close(respChan)

	for res := range respChan {
		allFiles = append(allFiles, res)
	}

	resp <- allFiles
}
