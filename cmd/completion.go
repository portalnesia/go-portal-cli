/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package cmd

import "github.com/spf13/cobra"

var completionCommand = &cobra.Command{
	Use:    "completion",
	Short:  "Generate the autocompletion script for the specified shell",
	Hidden: true,
}
