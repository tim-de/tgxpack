package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
)

type config struct{
	DefaultSource string
	DefaultDest string
}

func getConfigFilePath() string {
	configRoot := getConfigRoot()
	return filepath.Join(configRoot, "tgxpack", "tgxpack.ini")
}

func ensureConfigDirExists() error {
	configPath := getConfigFilePath()
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0700)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func readConfig(filepath string) (config, error) {
	configIni, err := ini.Load(filepath)
	if err != nil {
		return config{}, err
	}

	var configData config

	configData.DefaultDest = configIni.Section("Config").Key("DefaultDest").String()
	configData.DefaultSource = configIni.Section("Config").Key("DefaultSource").String()

	return configData, nil
}

func writeConfig(filepath string, configData config) error {
	configIni, err := ini.Load(filepath)
	if os.IsNotExist(err) {
		configIni = ini.Empty()
	} else if err != nil {
		return err
	}

	if configData.DefaultDest != "" {
		configIni.Section("Config").Key("DefaultDest").SetValue(configData.DefaultDest)
	}
	if configData.DefaultSource != "" {
		configIni.Section("Config").Key("DefaultSource").SetValue(configData.DefaultSource)
	}

	reterr := configIni.SaveTo(filepath)
	fmt.Fprintf(os.Stderr, "Wrote config to %s\n", filepath)
	return reterr
}
