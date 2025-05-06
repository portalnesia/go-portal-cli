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
	addServiceUseFlag bool
	addServiceConfig  config2.AddServiceConfig
)

var addServiceCmd = &cobra.Command{
	Use:   "add-service",
	Short: "Add new service",
	Long:  `Add new service and CRUD routes, handler, and usecase`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cfg config2.AddServiceConfig
			err error
		)

		appCtx := cmd.Context().Value("app")
		if appCtx == nil {
			panic("app is nil")
		}
		app, ok := appCtx.(*config.App)
		if !ok {
			panic("invalid app")
		}

		addServiceConfig.Module, err = helper.GetModuleName(app.Dir("go.mod"))
		if err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		cfg.Module = addServiceConfig.Module

		if addServiceUseFlag {
			cfg = addServiceConfig
		}
		if err = utils.PromptInitString("Service name", &cfg.Name, true, false); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		if err = utils.PromptInitString("Endpoint path", &cfg.Path, !addServiceUseFlag, true); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		if err = utils.PromptInitString("Endpoint version", &cfg.Version, !addServiceUseFlag, true); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		if !addServiceUseFlag {
			if err := utils.PromptInitBool("Override existing files", &cfg.Override); err != nil {
				_, _ = color.New(color.FgRed).Println("Error:", err)
				return
			}
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

		golang := bgolang.New(app)
		if err := golang.AddService(cfg); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
		}
	},
}

func init() {
	addServiceCmd.Flags().BoolVarP(&addServiceUseFlag, "flag", "f", false, "Use flag instead of prompt")

	// prompt
	addServiceCmd.Flags().StringVarP(&addServiceConfig.Name, "name", "n", "", "Service name")
	addServiceCmd.Flags().StringVarP(&addServiceConfig.Path, "path", "p", "", "Endpoint path, example: /users. Default use service name")
	addServiceCmd.Flags().StringVarP(&addServiceConfig.Version, "version", "v", "", "Endpoint version, example: /v1. Default without version")
	addServiceCmd.Flags().BoolVarP(&initConfig.Override, "override", "o", false, "Force override existing files")
}
