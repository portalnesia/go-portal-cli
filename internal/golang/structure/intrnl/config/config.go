/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package s_config

import (
	"go.portalnesia.com/portal-cli/internal/config"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"sync"
)

type Config struct {
	app *config.App
	cfg *config2.InitConfig
}

func Init(parentWg *sync.WaitGroup, app *config.App, cfg *config2.InitConfig, resp chan []config2.Builder) {
	defer parentWg.Done()

	c := &Config{
		app: app,
		cfg: cfg,
	}

	i := 7
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

	wg.Wait()
	close(respChan)

	for res := range respChan {
		allFiles = append(allFiles, res)
	}

	resp <- allFiles
}
