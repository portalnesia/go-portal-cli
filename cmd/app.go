/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package cmd

import (
	"context"
	"embed"
	"go.portalnesia.com/portal-cli/cmd/golang"
	"go.portalnesia.com/portal-cli/internal/config"

	"github.com/spf13/cobra"
)

func addCommands(rootCmd *cobra.Command) {
	// mark completion hidden
	rootCmd.AddCommand(completionCommand)

	rootCmd.AddCommand(golang.Init())
}

// Run adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Run(emb embed.FS) {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "portal-cli",
		Short: "CLI helper for Portalnesia",
		Long: `A simple CLI tool to speed up development for Portalnesia projects.

Currently, it focuses on generating boilerplate code structures quickly and consistently.

Designed to be lightweight and modular for future feature expansions.`,
		Version: config.GetVersion().String(),
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	app := config.New(emb)
	defer app.Close()

	ctx := context.WithValue(context.TODO(), "app", app)
	rootCmd.SetContext(ctx)

	addCommands(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
