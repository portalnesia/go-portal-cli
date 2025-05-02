/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package sserver

import (
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"strings"
	"sync"
)

func (s *Server) initConfigApp(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/server/config/app.go\n")

	src, _ := s.app.DataEmbed.ReadFile("data/golang/internal/server/config/app.txt")
	srcStr := string(src)
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME", s.cfg.Module)
	if s.cfg.Redis {
		srcStr = strings.ReplaceAll(srcStr, "{{IF_REDIS}}", `if sess, errSess := a.App.SessionStore.Get(c); errSess == nil {
		_ = sess.Save()
	}`)
	} else {
		srcStr = strings.ReplaceAll(srcStr, "{{IF_REDIS}}", ``)
	}

	src = []byte(srcStr)

	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/server/config/app.go",
	}
}
