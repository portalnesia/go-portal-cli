/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package golang

import (
	"github.com/fatih/color"
	"go.portalnesia.com/portal-cli/cmd/utils"
	"go.portalnesia.com/portal-cli/internal/config"
	bgolang "go.portalnesia.com/portal-cli/internal/golang"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"go.portalnesia.com/portal-cli/pkg/helper"
	"strings"

	"github.com/spf13/cobra"
)

var (
	addRepositoryUseFlag bool
	addRepositoryConfig  config2.AddRepositoryConfig
)

// addRepoCmd represents the addRepo command
var addRepositoryCmd = &cobra.Command{
	Use:   "add-repository",
	Short: "Add new repository",
	Long:  `Add new repository and interface`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cfg config2.AddRepositoryConfig
			err error
		)

		if addRepositoryUseFlag {
			cfg = addRepositoryConfig
		}

		if err = utils.PromptInitString("Name", &cfg.Name, !addEndpointUseFlag, false); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
			return
		}
		cfg.Name = strings.ToLower(cfg.Name)

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
		if err := golang.AddRepository(cfg); err != nil {
			_, _ = color.New(color.FgRed).Println("Error:", err)
		}
	},
}

func init() {
	addRepositoryCmd.Flags().BoolVarP(&addRepositoryUseFlag, "flag", "f", false, "Use flag instead of prompt")

	// prompt
	addRepositoryCmd.Flags().StringVarP(&addRepositoryConfig.Name, "name", "n", "", "Method name; example: FollowUser")
}
