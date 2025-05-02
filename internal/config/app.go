/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package config

import (
	"embed"
	"fmt"
	"github.com/subosito/gotenv"
	"os"
)

type App struct {
	Production bool
	DataEmbed  embed.FS
}

func New(emb embed.FS) *App {
	_ = gotenv.Load()

	return &App{
		DataEmbed:  emb,
		Production: os.Getenv("ENV") != "development",
	}
}

func (a *App) Close() {

}

func (a *App) Dir(dir string) string {
	if a.Production {
		return fmt.Sprintf("testst/%s", dir)
	}

	return fmt.Sprintf("tmp/%s", dir)
}
