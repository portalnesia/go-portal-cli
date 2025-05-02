/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package golang

import (
	"github.com/spf13/cobra"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
)

var (
	globalConfig config2.GlobalConfig
)

var golangCmd = &cobra.Command{
	Use:   "golang",
	Short: "Golang helper",
	Long:  `Helper for golang programming language`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Init() *cobra.Command {
	golangCmd.AddCommand(initCmd)
	return golangCmd
}

func init() {
	golangCmd.PersistentFlags().BoolVarP(&globalConfig.ServerDirectory, "server", "s", true, "Handler, routes, middleware, etc files in server directory")
}
