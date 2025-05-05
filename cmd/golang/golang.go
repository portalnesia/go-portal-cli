/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package golang

import (
	"github.com/spf13/cobra"
)

var golangCmd = &cobra.Command{
	Use:   "golang",
	Short: "Golang helper",
	Long:  `Helper for golang programming language`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Init() *cobra.Command {
	golangCmd.AddCommand(initCmd)
	golangCmd.AddCommand(addServiceCmd)
	golangCmd.AddCommand(addEndpointCmd)
	return golangCmd
}
