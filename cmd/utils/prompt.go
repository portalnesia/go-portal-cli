/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package utils

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func PromptInitBool(name string, value *bool) error {
	yellow := color.New(color.FgYellow)

	var tmp string
	_, _ = yellow.Printf("%s? (y/N) ", name)
	_, _ = fmt.Scanln(&tmp)
	if tmp == "" {
		return fmt.Errorf("%s required", name)
	}
	*value = strings.ToLower(tmp) == "y"
	return nil
}

func PromptInitString(name string, value *string, forcePrompt ...bool) error {
	yellow := color.New(color.FgYellow)

	var tmp string
	if value != nil {
		tmp = *value
	}
	if len(forcePrompt) > 0 && forcePrompt[0] && tmp == "" {
		_, _ = yellow.Printf("%s? ", name)
		_, _ = fmt.Scanln(&tmp)
	}
	if tmp == "" {
		return fmt.Errorf("%s required", name)
	}
	*value = tmp
	return nil
}
