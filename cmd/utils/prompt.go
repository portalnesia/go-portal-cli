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

// PromptInitString
//
// forcePrompt: true, will force when value is empty, otherwise when value is empty, it will skip
// optional: true, will skip if empty
func PromptInitString(name string, value *string, forcePromptAndOptional ...bool) error {
	yellow := color.New(color.FgYellow)

	var tmp string
	if value != nil {
		tmp = *value
	}
	if len(forcePromptAndOptional) > 0 && forcePromptAndOptional[0] && tmp == "" {
		_, _ = yellow.Printf("%s? ", name)
		_, _ = fmt.Scanln(&tmp)
	}

	if (len(forcePromptAndOptional) <= 1 || (len(forcePromptAndOptional) > 1 && !forcePromptAndOptional[1])) && tmp == "" {
		return fmt.Errorf("%s required", name)
	}
	*value = tmp
	return nil
}
