//go:build !windows

package main

import (
	"os"
	"path/filepath"
)

func getConfigRoot() string {
	confroot := os.Getenv("XDG_CONFIG_HOME")

	if confroot == "" {
		home := os.Getenv("HOME")
		confroot = filepath.Join(home, ".config")
	}

	return confroot
}
