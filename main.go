/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */
package main

import (
	"embed"
	"go.portalnesia.com/portal-cli/cmd"
)

var (
	//go:embed data/*
	embedData embed.FS
)

func main() {
	cmd.Run(embedData)
}
