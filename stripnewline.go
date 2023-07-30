//go:build !windows

package main

import "strings"

func stripNewLine(str string) string {
	return strings.TrimRight(str, "\n")
}
