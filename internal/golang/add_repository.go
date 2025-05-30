/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package b_golang

import (
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/internal/golang/service"
)

func (g *Golang) AddRepository(cfg config2.AddRepositoryConfig) error {
	_, _ = color.New(color.FgBlue).Printf("\nPlease wait...\n")

	builder, err := service.AddRepository(g.app, cfg)
	if err != nil {
		return err
	}

	if err = g.Build(builder); err != nil {
		return err
	}

	return nil
}
