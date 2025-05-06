/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package utils

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

func PromptInitBool(name string, value *bool) error {
	yellow := color.New(color.FgYellow)

	var tmp string
	for tmp == "" {
		_, _ = yellow.Printf("%s? (y/N) ", name)
		_, _ = fmt.Scanln(&tmp)
		if tmp == "" {
			fmt.Print("\033[1A")
			fmt.Print("\033[2K")
		}
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

	// force to prompt question
	if len(forcePromptAndOptional) > 0 && forcePromptAndOptional[0] && tmp == "" {
		scanner := bufio.NewScanner(os.Stdin)
		// required
		if len(forcePromptAndOptional) <= 1 || (len(forcePromptAndOptional) > 1 && !forcePromptAndOptional[1]) {
			for tmp == "" {
				_, _ = yellow.Printf("%s? ", name)
				if scanner.Scan() {
					tmp = scanner.Text()
				}
				if tmp == "" {
					fmt.Print("\033[1A")
					fmt.Print("\033[2K")
				}
			}
		} else { // optional
			_, _ = yellow.Printf("%s? ", name)
			if scanner.Scan() {
				tmp = scanner.Text()
			}
		}
	}
	*value = tmp
	return nil
}
