package main

import (
	"log"
	"path/filepath"

	"github.com/tim-de/tgxlib"
)

func extractMod(modpath, outdir string) error {
	modpath = filepath.Clean(modpath)

	var tgxfile tgxlib.TgxFile

	outdir, err := filepath.Abs(outdir)
	if err != nil {
		return err
	}

	tgxfile, err = tgxlib.ReadFromFile(modpath)
	if err != nil {
		return err
	}

	for _, subfile := range tgxfile.SubFiles {
		log.Printf("Dumping subfile %s\n", subfile.FilePath)
		err = subfile.Dump(outdir)
		if err != nil {
			return err
		}
	}

	return nil
}
