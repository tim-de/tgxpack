package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	
	"github.com/tim-de/tgxlib"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func newModDir(modpath string) error {

	var modinfo, modchanges *os.File

	err := os.MkdirAll(filepath.Join(modpath, "DATA"), 0700)
	if err != nil {
		return err
	}

	err = os.Chdir(modpath)
	if err != nil {
		return err
	}

	modinfo, err = os.Create(filepath.FromSlash("./DATA/ModInfo.ini"))
	if err != nil {
		return err
	}
	defer modinfo.Close()

	// This is the contents of a blank ModInfo.ini file, with all keys
	// present but no values assigned to them
	modinfodata := []byte{
		0x5b, 0x4d, 0x6f, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x5d, 0x0d, 0x0a, 0x54,
		0x77, 0x6f, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49,
		0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x3d, 0x0d, 0x0a,
		0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x3d, 0x0d, 0x0a, 0x41, 0x75,
		0x74, 0x68, 0x6f, 0x72, 0x3d, 0x0d, 0x0a, 0x41, 0x75, 0x74, 0x68, 0x6f,
		0x72, 0x5f, 0x77, 0x77, 0x77, 0x3d, 0x0d, 0x0a, 0x41, 0x75, 0x74, 0x68,
		0x6f, 0x72, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x3d, 0x0d, 0x0a, 0x4b,
		0x6f, 0x68, 0x61, 0x6e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x3d,
		0x4b, 0x41, 0x47, 0x20, 0x31, 0x2e, 0x33, 0x2e, 0x37, 0x0d, 0x0a, 0x4b,
		0x4d, 0x6f, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x3d, 0x0d,
		0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x3d}

	_, err = modinfo.Write(modinfodata)
	if err != nil {
		return err
	}

	modchanges, err = os.Create(filepath.FromSlash("./DATA/ModChanges.txt"))
	if err != nil {
		return err
	}
	defer modchanges.Close()
	
	return nil
}

func addFileToMod(source_path, mod_dir, query string, matchnum int) error {
	source_file, err := tgxlib.ReadFromFile(source_path)
	if err != nil {
		return err
	}

	subfile_names := make([]string, source_file.FileCount)
	for ix, subfile := range source_file.SubFiles {
		subfile_names[ix] = tgxlib.OsToInternalPath(subfile.FilePath)
	}

	matches := fuzzy.RankFindNormalizedFold(query, subfile_names)
	sort.Sort(matches)

	var pos int = 0
	var count int
	
	var selected_indices []int

Pageloop:
	for {
		if pos >= len(matches) {
			break
		}
		rest := matches[pos:]
		if len(rest) < matchnum {
			count = len(rest)
		} else {
			count = matchnum
		}
		
		fmt.Fprintf(os.Stderr, "Results %d-%d of %d for \"%s\" in \"%s\":\n", pos+1, pos+count, len(matches), query, source_path)
		fmt.Fprintln(os.Stderr, "Type a number or a space-separated list of numbers to select files, or \"n\" to see more results")
		for ix, match := range rest[:count] {
			fmt.Fprintf(os.Stderr, "[%d] -> %s\n", ix, match.Target)
		}

		var choices string
		fmt.Fprintf(os.Stderr, "Choose file(s) to import [0-%d]: ", count-1)

		stdin := bufio.NewReader(os.Stdin)
		choices, err = stdin.ReadString('\n')
		if err != nil {
			return err
		}
		choices = stripNewLine(choices)
		if len(choices) == 0 {
			break
		}
		
		for _, choice := range strings.Split(choices, " ") {
			if choice[0] == 'n' || choice [0] == 'N' {
				pos += count
				continue Pageloop
			}
			if choiceix, err := strconv.Atoi(choice); err == nil {
				selected_indices = append(selected_indices, choiceix+pos)
			}
		}
		break
	}

	for _, choice := range selected_indices {
		file_ix := matches[choice].OriginalIndex
		subfile := source_file.SubFiles[file_ix]
		fmt.Fprintf(os.Stderr, "Unpacking \"%s\" into \"%s\"\n", tgxlib.OsToInternalPath(subfile.FilePath), mod_dir)
		subfile.Dump(mod_dir)
	}
	
	return nil
}
