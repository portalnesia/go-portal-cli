/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package scmd

import (
	"fmt"
	"github.com/fatih/color"
	config2 "go.portalnesia.com/portal-cli/internal/golang/config"
	"path"
	"strings"
	"sync"
)

func (s *Cmd) addStatic(file string, wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating %s.go\n", file)

	src, _ := s.app.DataEmbed.ReadFile(fmt.Sprintf("data/golang/%s.txt", file))
	srcStr := string(src)
	srcStr = strings.ReplaceAll(srcStr, "CLI_NAME", path.Base(s.cfg.Module))
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME", s.cfg.Module)
	src = []byte(srcStr)

	res <- config2.Builder{
		Static:   src,
		Pathname: fmt.Sprintf("%s.go", file),
	}
}
