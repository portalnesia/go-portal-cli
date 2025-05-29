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

type addService struct {
	app *config.App
	cfg config2.AddServiceConfig
}

func AddService(app *config.App, cfg config2.AddServiceConfig) ([]config2.Builder, error) {
	c := &addService{
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

	go c.addServiceHandler(wg, respChan)
	go c.addServiceUsecase(wg, respChan)
	go c.addServiceRoutes(wg, respChan)
	go c.addServiceRepository(wg, respChan)
	go c.addToRoutes(wg, respChan)
	go c.addToRepository(wg, respChan)

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
