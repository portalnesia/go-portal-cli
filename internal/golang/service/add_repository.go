/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package service

import (
	"go.portalnesia.com/portal-cli/internal/config"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"sync"
)

type addRepository struct {
	app *config.App
	cfg config2.AddRepositoryConfig
}

func AddRepository(app *config.App, cfg config2.AddRepositoryConfig) ([]config2.Builder, error) {
	c := &addRepository{
		app: app,
		cfg: cfg,
	}

	i := 2
	if !cfg.NoModel {
		i += 1
	}
	var (
		allFiles = make([]config2.Builder, 0)
		respChan = make(chan config2.Builder, i)
		wg       = &sync.WaitGroup{}
	)
	wg.Add(i)

	go c.addServiceRepository(wg, respChan)
	go c.addToRepository(wg, respChan)
	if !cfg.NoModel {
		go c.addRepositoryModel(wg, respChan)
	}

	wg.Wait()
	close(respChan)

	for res := range respChan {
		if res.Err != nil {
			return nil, res.Err
		}
		allFiles = append(allFiles, res)
	}

	return allFiles, nil
}
