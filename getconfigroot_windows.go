//go:build windows

package main

import (
	"os"
)

func getConfigRoot() string {
	confroot := os.Getenv("LOCALAPPDATA")

	if confroot != "" {
		return confroot
	}

	confroot = os.Getenv("APPDATA")

	if confroot != "" {
		return confroot
	}

	return "."
}
