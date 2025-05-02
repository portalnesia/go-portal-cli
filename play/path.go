/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package main

import (
	"fmt"
	"path"
)

func main() {
	pathname := "internal/config/config.go"

	fmt.Println(path.Dir(pathname))
}
