/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package ginit

import (
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"sync"
)

func (c *initType) initConfigFirebase(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/firebase.go\n")

	src, _ := c.app.DataEmbed.ReadFile("data/golang/internal/config/firebase.txt")

	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/config/firebase.go",
	}
}
