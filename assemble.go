package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/tim-de/tgxlib"
	"github.com/go-ini/ini"
)

func buildMod(modroot, outpath string) error {
	
	modroot = filepath.Clean(modroot)

	rootinfo, err := os.Stat(modroot)
	if err != nil {
		return err
	}
	if !rootinfo.IsDir() {
		return fmt.Errorf("%s is not a directory", modroot)
	}

	fmt.Fprintf(os.Stderr, "Moving to %s\n", modroot)
	err = os.Chdir(modroot)
	if err != nil {
		return err
	}

	fmt.Fprint(os.Stderr, "Loading Version and 2 character identifier\n")
	str_version, short_id, err := getVersionAndId()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "\tVersion: %s\n\tIdentifier: %s\n", str_version, short_id)

	version, err := tgxlib.UnpackVersionFromString(str_version)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Finding files in %s\n", modroot)
	pathlist, err := tgxlib.FindFilesRecursive(".")
	if err != nil {
		return err
	}
	for _, filename := range pathlist {
		fmt.Fprintf(os.Stderr, "\tFound file: %s\n", filename)
	}

	tgxfile, err := tgxlib.FromPathList(version, short_id, pathlist)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Writing file %s\n", outpath)
	err = tgxfile.WriteFile(outpath)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Done!")
	return nil
}

func getVersionAndId() (string, string, error) {
	pathlist, err := filepath.Glob(filepath.FromSlash("./[dD][aA][tT][aA]/[mM][oO][dD][iI][nN][fF][oO].[iI][nN][iI]"))
	if err != nil {
		return "", "", err
	}
	var short_id, version string
	for _, path := range pathlist {
		if strings.ToUpper(path) == filepath.FromSlash("DATA/MODINFO.INI") {
			modinfo, err := ini.Load(path)
			if err != nil {
				return "", "", err
			}
			short_id = modinfo.Section("ModInfo").Key("TwoCharacterIdentifier").String()
			version = modinfo.Section("ModInfo").Key("Version").String()
			break
		}
	}
	if version == "" {
		return "", "", errors.New("Version not found in Data\\ModInfo.ini")
	}
	if short_id == "" {
		return "", "", errors.New("Identifier not found in Data\\ModInfo.ini")
	}

	return version, short_id, nil
}
