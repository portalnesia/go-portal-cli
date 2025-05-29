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
	"strings"
	"sync"
)

func (c *initType) initServerHandlerUtils(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/rest/handler/utils.go\n")

	src, _ := c.app.DataEmbed.ReadFile("data/golang/internal/rest/handler/utils.txt")
	srcStr := string(src)
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME", c.cfg.Module)
	if c.cfg.Redis {
		srcStr = strings.ReplaceAll(srcStr, "{{IF_REDIS}}", `if sess, errSess := app.SessionStore().Get(c); errSess == nil {
		_ = sess.Save()
	}`)
	} else {
		srcStr = strings.ReplaceAll(srcStr, "{{IF_REDIS}}", ``)
	}

	src = []byte(srcStr)

	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/rest/handler/utils.go",
	}
}
