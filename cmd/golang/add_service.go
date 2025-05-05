/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package golang

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.portalnesia.com/portal-cli/cmd/utils"
	"go.portalnesia.com/portal-cli/internal/config"
	bgolang "go.portalnesia.com/portal-cli/internal/golang"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
	"strings"
)

var (
	newServiceUseFlag bool
	newServiceConfig  config2.AddServiceConfig
)

var newServiceCmd = &cobra.Command{
	Use:   "add-service",
	Short: "Add new service",
	Long:  `Add new service and CRUD routes, handler, and usecase`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cfg config2.AddServiceConfig
			err error
		)

		if newServiceUseFlag {
			cfg = newServiceConfig
		}
		if err = utils.PromptInitString("Service name", &cfg.Name, !newServiceUseFlag, false); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		if err = utils.PromptInitString("Endpoint path", &cfg.Path, !newServiceUseFlag, true); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		if err = utils.PromptInitString("Endpoint version", &cfg.Version, !newServiceUseFlag, true); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		cfg.Name = strings.ToLower(cfg.Name)

		if cfg.Path == "" {
			cfg.Path = cfg.Name
		}
		cfg.Path = strings.ToLower(cfg.Path)

		if cfg.Version != "" {
			cfg.Version = strings.ToLower(cfg.Version)
			if !strings.HasPrefix(cfg.Version, "v") {
				cfg.Version = "v" + cfg.Version
			}
		}

		appCtx := cmd.Context().Value("app")
		if appCtx == nil {
			panic("app is nil")
		}
		app, ok := appCtx.(*config.App)
		if !ok {
			panic("invalid app")
		}

		cfg.Module, err = helper.GetModuleName(app.Dir("go.mod"))
		if err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}

		golang := bgolang.New(app)
		if err := golang.NewService(cfg); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
		}
	},
}

func init() {
	newServiceCmd.Flags().BoolVarP(&newServiceUseFlag, "flag", "f", false, "Use flag instead of prompt")

	// prompt
	newServiceCmd.Flags().StringVarP(&newServiceConfig.Name, "name", "n", "", "Service name")
	newServiceCmd.Flags().StringVarP(&newServiceConfig.Path, "path", "p", "", "Endpoint path, example: /users. Default use service name")
	newServiceCmd.Flags().StringVarP(&newServiceConfig.Version, "version", "v", "", "Endpoint version, example: /v1. Default without version")
}
