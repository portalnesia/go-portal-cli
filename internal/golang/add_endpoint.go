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

func (g *Golang) AddEndpoint(cfg config2.AddEndpointConfig) error {
	_, _ = color.New(color.FgBlue).Printf("\nPlease wait...\n")

	builder, err := service.AddEndpoint(g.app, cfg)
	if err != nil {
		return err
	}

	if err = g.Build(builder, true); err != nil {
		return err
	}

	return nil
}
