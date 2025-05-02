//go:build ignore

/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package main

import (
	"flag"
	"github.com/coreos/go-semver/semver"
	"html/template"
	"os"
	"strings"
)

func main() {
	var (
		tag   string
		patch bool
		minor bool
		major bool
	)

	flag.StringVar(&tag, "tag", "", "set version")
	flag.BoolVar(&patch, "patch", false, "add patch version")
	flag.BoolVar(&minor, "minor", false, "add patch version")
	flag.BoolVar(&major, "major", false, "add patch version")
	flag.Parse()

	versionGenByte, err := os.ReadFile("internal/config/version_gen.go")
	die(err)

	versionGen := string(versionGenByte)

	versionSplit := strings.Split(versionGen, `"`)
	versionString := versionSplit[1]
	version, err := semver.NewVersion(versionString)
	die(err)

	if tag != "" {
		tmp, err := semver.NewVersion(tag)
		die(err)
		version = tmp
	} else {
		if patch {
			version.Patch += 1
		}
		if minor {
			version.Minor += 1
			version.Patch = 0
		}
		if major {
			version.Major += 1
			version.Minor = 0
			version.Patch = 0
		}
	}

	textTemplate, err := template.Must(template.New(""), nil).Parse(`/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

// Code generated! DO NOT EDIT!

package config

var tmpVersion = "{{ .Version }}"`)
	die(err)

	f, err := os.Create("internal/config/version_gen.go")
	die(err)
	defer f.Close()

	textTemplate.Execute(f, struct {
		Version string
	}{
		Version: version.String(),
	})
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
