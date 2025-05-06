/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package golang

import (
	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.portalnesia.com/portal-cli/cmd/utils"
	"go.portalnesia.com/portal-cli/internal/config"
	bgolang "go.portalnesia.com/portal-cli/internal/golang"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
	util "go.portalnesia.com/utils"
	"strings"
)

var (
	addEndpointUseFlag bool
	addEndpointConfig  config2.AddEndpointConfig
)

var addEndpointCmd = &cobra.Command{
	Use:   "add-endpoint",
	Short: "Add new endpoint to a service",
	Long:  `Add new endpoint to existing routes, handler, and usecase`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cfg config2.AddEndpointConfig
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

		if addEndpointUseFlag {
			cfg = addEndpointConfig
		}
		if err = utils.PromptInitString("Service name", &cfg.ServiceName, true, false); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		if err = utils.PromptInitString("Method name", &cfg.Name, true, false); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		if err = utils.PromptInitString("Endpoint path", &cfg.Path, true, false); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		if err = utils.PromptInitString("HTTP method", &cfg.Method, true, false); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}

		method := strings.ToLower(cfg.Method)
		if !lo.Contains([]string{
			"get",
			"post",
			"put",
			"patch",
			"delete",
		}, method) {
			_, _ = color.New(color.FgRed).Println("Error: Invalid HTTP method")
			return
		}
		if method == "patch" {
			cfg.Method = "put"
		}
		cfg.Method = util.FirstToUpper(strings.ToLower(cfg.Method))
		cfg.ServiceName = strings.ToLower(cfg.ServiceName)
		cfg.Name = strings.Join(
			lo.Map[string, string](strings.Split(cfg.Name, " "), func(item string, index int) string {
				return util.FirstToUpper(item)
			}),
			"",
		)
		cfg.Path = strings.ToLower(cfg.Path)

		golang := bgolang.New(app)
		if err := golang.AddEndpoint(cfg); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
		}
	},
}

func init() {
	addEndpointCmd.Flags().BoolVarP(&addEndpointUseFlag, "flag", "f", false, "Use flag instead of prompt")

	// prompt
	addEndpointCmd.Flags().StringVarP(&addEndpointConfig.ServiceName, "service", "s", "", "Service name; example: user")
	addEndpointCmd.Flags().StringVarP(&addEndpointConfig.Name, "name", "n", "", "Method name; example: FollowUser")
	addEndpointCmd.Flags().StringVarP(&addEndpointConfig.Path, "path", "p", "", "Endpoint path indlude version, example: /:id/follow")
	addEndpointCmd.Flags().StringVarP(&addEndpointConfig.Method, "method", "m", "", "HTTP method, example: GET, POST, PUT, PATCH, DELETE")
}
