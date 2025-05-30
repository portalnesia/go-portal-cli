/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package b_golang

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	ginit "go.portalnesia.com/portal-cli/internal/golang/init"
	"os"
	"sync"
)

func (g *Golang) Init(cfg config2.InitConfig) error {
	_, _ = color.New(color.FgBlue).Printf("\nPlease wait...\n")

	var (
		i       = 2
		resChan = make(chan []config2.Builder, i)
		wg      = &sync.WaitGroup{}
	)
	wg.Add(i)

	go g.generateConfig(wg, cfg, resChan)
	go ginit.Init(wg, g.app, &cfg, resChan)

	wg.Wait()
	close(resChan)

	dirs := []string{
		g.app.Dir("internal/service"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return errors.Join(fmt.Errorf("failed to create directory %s", d), err)
		}
	}

	for builder := range resChan {
		if err := g.Build(builder); err != nil {
			return errors.Join(fmt.Errorf("failed to build %s", builder[0].Pathname), err)
		}
	}

	return nil
}

func (g *Golang) generateConfig(wg *sync.WaitGroup, cfg config2.InitConfig, resp chan []config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generate config.json\n")

	data := map[string]any{
		"ports": []int{4000},
		"secret": map[string]any{
			"crypto": "",
		},
		"db": map[string]any{
			"host":     "localhost",
			"port":     3306,
			"user":     "",
			"password": "",
			"database": "",
		},
		"link": map[string]any{
			"api":    "",
			"web":    "",
			"static": "",
		},
	}
	if cfg.Redis {
		data["redis"] = map[string]any{
			"host":     "localhost",
			"port":     6379,
			"user":     "",
			"password": "",
			"database": 0,
		}
	}

	byt, _ := json.MarshalIndent(data, "", "    ")

	files := []config2.Builder{
		{
			Pathname: "config.json",
			Static:   byt,
		},
	}

	resp <- files
}
