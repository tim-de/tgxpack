package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var configPath string = ".tgxpack.ini"

func main() {
	app := &cli.App{
		Name: "tgxpack",
		Usage: "Create and modify tgx archive files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "config",
				Aliases: []string{"c"},
				Usage: "use `FILE` for config",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "pack",
				Aliases: []string{"p"},
				Usage: "pack a directory tree into a tgx file",
				ArgsUsage: "MODDIR",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "output",
						Aliases: []string{"o"},
						Usage: "save mod to `FILE`",
					},
					&cli.StringFlag{
						Name: "dest",
						Aliases: []string{"d"},
						Usage: "save mod in `DIR`",
					},
				},
				Action: packCommand,
			},
			{
				Name: "unpack",
				Aliases: []string{"u"},
				Usage: "unpack a tgxfile into a directory tree",
				ArgsUsage: "MODFILE",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "output",
						Aliases: []string{"o"},
						Usage: "unpack mod to `DIR`",
					},
				},
				Action: unpackCommand,
			},
			{
				Name: "new",
				Aliases: []string{"n"},
				Usage: "initialise a new mod directory",
				ArgsUsage: "NAME",
				Action: newCommand,
			},
			{
				Name: "add",
				Aliases: []string{"a"},
				Usage: "add a file to the current mod directory",
				ArgsUsage: "QUERY",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "source",
						Aliases: []string{"s"},
						Usage: "search within `FILE`",
					},
					&cli.StringFlag{
						Name: "dest",
						Aliases: []string{"d"},
						Usage: "unpack file to modroot `DIR`",
						Value: ".",
						DefaultText: filepath.FromSlash("./"),
					},
					&cli.IntFlag{
						Name: "count",
						Aliases: []string{"c"},
						Usage: "Display `N` results",
						Value: 10,
						DefaultText: "10",
					},
				},
				Action: addCommand,
			},
			{
				Name: "setup",
				Aliases: []string{"s"},
				Usage: "set values in the configuration file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "dest",
						Aliases: []string{"d"},
						Usage: "save .tgx files into `FILE` by default",
					},
					&cli.StringFlag{
						Name: "source",
						Aliases: []string{"s"},
						Usage: "search for subfiles in `FILE` by default",
					},
				},
				Action: setupCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getConfigFromCli(cCtx *cli.Context) (config, error) {
	configPath := cCtx.String("config")
	if configPath == "" {
		configPath = getConfigFilePath()
	}

	return readConfig(configPath)
}

func filenameWithoutExt(path string) string {
	return path[:len(path) - len(filepath.Ext(path))]
}

func packCommand(cCtx *cli.Context) error {
	var targetPath string
	var err error
	if !cCtx.Args().Present() {
		targetPath, err = filepath.Abs(".")
		if err != nil {
			return err
		}
	} else {
		targetPath, err = filepath.Abs(cCtx.Args().First())
		if err != nil {
			return err
		}
	}

	config, err := getConfigFromCli(cCtx)
	if err != nil {
		return err
	}

	outpath := cCtx.String("output")
	if outpath == "" {
		outpath = filepath.Base(targetPath)
	} else {
		outpath = filenameWithoutExt(outpath)
	}
	outpath = fmt.Sprintf("%s.tgx", outpath)

	dest := cCtx.String("dest")
	if dest == "" {
		dest = config.DefaultDest
	}

	outpath = filepath.Join(dest, outpath)
	
	outpath, err = filepath.Abs(outpath)
	if err != nil {
		return err
	}
	return buildMod(targetPath, outpath)
}

func unpackCommand(cCtx *cli.Context) error {
	if !cCtx.Args().Present() {
		fmt.Fprintln(os.Stderr, "Please supply a mod file to unpack. Type 'tgxpack help unpack' for more information")
		return nil
	}
	outpath := cCtx.String("output")
	if outpath == "" {
		outpath= filenameWithoutExt(filepath.Base(cCtx.Args().First()))
	}
	return extractMod(cCtx.Args().First(), outpath)
}

func newCommand(cCtx *cli.Context) error {
	if !cCtx.Args().Present() {
		fmt.Fprintln(os.Stderr, "Please supply a name for the mod. Type 'tgxpack help new' for more information")
		return nil
	}
	
	moddir, err := filepath.Abs(cCtx.Args().First())
	if err != nil {
		return err
	}
	return newModDir(moddir)
}

func addCommand(cCtx *cli.Context) error {
	outpath := cCtx.String("dest")
	source := cCtx.String("source")
	matchnum := cCtx.Int("count")
	if !cCtx.Args().Present() {
		fmt.Fprintln(os.Stderr, "Please supply a query string. Type 'tgxpack help add' for more information")
		return nil
	}

	config, err := getConfigFromCli(cCtx)
	if err != nil {
		return err
	}

	if source == "" {
		source = config.DefaultSource
	}
	
	query := cCtx.Args().First()
	moddir, err := filepath.Abs(outpath)
	if err != nil {
		return err
	}
					
	return addFileToMod(source, moddir, query, matchnum)
}

func setupCommand(cCtx *cli.Context) error {
	source := cCtx.String("source")
	dest := cCtx.String("dest")

	if source == "" && dest == "" {
		return nil
	}

	// This is a hack to get around string escaping on windows. Where
	// a path gets autocompleted to a directory it has a backslash
	// appended. This then leads to the ending doublequote getting
	// left on the end of the string, causing problems when you need
	// to pack something and it's an invalid filepath because there's
	// a doublequote on the end of the directory name
	if len(source) > 0 && source[len(source)-1] == '"' {
		source = source[:len(source)-1]
	}
	// Yeah there's definitely a proper solution but it was complicated
	// so I chose the dumb one
	if len(dest) > 0 && dest[len(dest)-1] == '"' {
		dest = dest[:len(dest)-1]
	}

	configPath := cCtx.String("config")
	if configPath == "" {
		configPath = getConfigFilePath()
	}

	if err := ensureConfigDirExists(); err != nil {
		return err
	}

	configData := config{
		DefaultSource: source,
		DefaultDest: dest,
	}

	return writeConfig(configPath, configData)
}
