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

func (c *initType) initConfigRedis(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/redis.go\n")

	src, _ := c.app.DataEmbed.ReadFile("data/golang/internal/config/redis.txt")
	srcStr := string(src)
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME", c.cfg.Module)
	src = []byte(srcStr)

	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/config/redis.go",
	}
}

func (c *initType) initGetter(wg *sync.WaitGroup, res chan<- config2.Builder) {
	defer wg.Done()
	_, _ = color.New(color.FgBlue).Printf("Generating internal/config/getter.go\n")

	src, _ := c.app.DataEmbed.ReadFile("data/golang/internal/config/getter.txt")
	srcStr := string(src)
	srcStr = strings.ReplaceAll(srcStr, "APP_NAME", c.cfg.Module)
	if c.cfg.Redis {
		srcStr = strings.ReplaceAll(srcStr, "{{IMPORT_IF_REDIS}}", `"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/session"
    "github.com/redis/go-redis/v9"`)

		srcStr = strings.ReplaceAll(srcStr, "{{IF_REDIS}}", `func (a *app) Redis() iface.RedisInterface {
	return a.redis
}

func (a *app) FiberStorage() fiber.Storage {
	return a.fiberStorage
}

func (a *app) SessionStore() *session.Store {
	return a.sessionStore
}`)
	} else {
		srcStr = strings.ReplaceAll(srcStr, "{{IMPORT_IF_REDIS}}", ``)
		srcStr = strings.ReplaceAll(srcStr, "{{IF_REDIS}}", ``)
	}

	src = []byte(srcStr)

	res <- config2.Builder{
		Static:   src,
		Pathname: "internal/config/getter.go",
	}
}
