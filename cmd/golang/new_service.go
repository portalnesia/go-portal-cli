/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package golang

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.portalnesia.com/portal-cli/cmd/utils"
	"go.portalnesia.com/portal-cli/internal/config"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
)

var (
	newServiceName    string
	newServicePath    string
	newServiceVersion string
	newServiceUseFlag bool
	newServiceConfig  config2.NewServiceConfig
)

var newServiceCmd = &cobra.Command{
	Use:   "new-service",
	Short: "Create new service",
	Long:  `Create new service and add CRUD routes, handler, and usecase`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cfg config2.NewServiceConfig
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
		if cfg.Path == "" {
			cfg.Path = cfg.Name
		}

		fmt.Printf("%+v", cfg)
	},
}

func init() {
	newServiceCmd.Flags().BoolVarP(&newServiceUseFlag, "flag", "f", false, "Use flag instead of prompt")

	// prompt
	newServiceCmd.Flags().StringVarP(&newServiceConfig.Name, "name", "n", "", "Service name")
	newServiceCmd.Flags().StringVarP(&newServiceConfig.Path, "path", "p", "", "Endpoint path, example: /users. Default use service name")
	newServiceCmd.Flags().StringVarP(&newServiceConfig.Version, "version", "v", "", "Endpoint version, example: /v1. Default without version")
}
