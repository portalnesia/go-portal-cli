/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package helper

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"
)

func BodyListNewLines() ast.Stmt {
	return &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  token.STRING,
			Value: ``,
		}, // dummy expression, tidak valid
	}
}

// GetModuleName reads the go.mod file and returns the module name
func GetModuleName(goModPath string) (string, error) {
	file, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			return moduleName, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("module name not found in %s", goModPath)
}

func FirstToLower(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
