/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package config

import (
	"github.com/coreos/go-semver/semver"
)

var version = semver.New(tmpVersion)

func GetVersion() *semver.Version {
	return version
}
