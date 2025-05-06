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
)

var (
	initWithAll bool
	initUseFlag bool
	initConfig  config2.InitConfig
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init structure directory",
	Long:  `Golang helper for init structure directory`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cfg config2.InitConfig
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

		if initUseFlag || initWithAll {
			cfg = initConfig
			if initWithAll {
				cfg.Redis = true
				cfg.Firebase = true
				cfg.Handlebars = true
			}
		} else {
			if err := utils.PromptInitBool("Add redis", &cfg.Redis); err != nil {
				_, _ = color.New(color.FgRed).Println("Error:", err)
				return
			}
			if err := utils.PromptInitBool("Add firebase", &cfg.Firebase); err != nil {
				_, _ = color.New(color.FgRed).Println("Error:", err)
				return
			}
			if err := utils.PromptInitBool("Add handlebars", &cfg.Handlebars); err != nil {
				_, _ = color.New(color.FgRed).Println("Error:", err)
				return
			}
			if err := utils.PromptInitBool("Override existing files", &cfg.Override); err != nil {
				_, _ = color.New(color.FgRed).Println("Error:", err)
				return
			}
		}

		golang := bgolang.New(app)
		if err := golang.Init(cfg); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
		}
	},
}

func init() {
	initCmd.Flags().BoolVarP(&initWithAll, "all", "a", false, "Add all library")
	initCmd.Flags().BoolVarP(&initUseFlag, "flag", "f", false, "Use flag instead of prompt")

	// prompt
	initCmd.Flags().StringVar(&initConfig.Module, "module", "", "Module name")
	initCmd.Flags().BoolVar(&initConfig.Redis, "redis", false, "Add redis")
	initCmd.Flags().BoolVar(&initConfig.Firebase, "firebase", false, "Add firebase")
	initCmd.Flags().BoolVar(&initConfig.Handlebars, "handlebars", false, "Add handlebars")
	initCmd.Flags().BoolVarP(&initConfig.Override, "override", "o", false, "Force override existing files")
}
