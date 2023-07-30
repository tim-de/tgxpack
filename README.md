# tgxpack: A mod packer for Kohan

This is a simple cli application for creating and modifying .tgx mod files
for Kohan games. It functions by creating a directory structure based on 
an existing mod file (or creating a blank one) which can have any desired
changes made before the mod is then packed into a .tgx archive which can
be read by the Kohan game engine.

The following commands are provided for doing this:
 * setup
 * unpack
 * pack
 * new
 * add
 
### Setup
This command is not directly involved in the process of creating mods
but streamlines the other commands a lot by setting defaults for some
of the command line flags which would otherwise be compulsory.

The 'source' flag sets the default file for the 'add' command to search
in, and is usually going to be 'KAG_INSTALL_DIRECTORY\Kohan_AG.tgw'

The 'dest' flag is the default directory in which mods should be saved,
so setting this to wherever Kohan is installed is the easiest thing to
do.

### Unpack
Extracts the embedded directory structure from a tgx/tgw file. The path/name
of the resulting directory can be specified with the 'output' command line
flag.

### Pack
Packs all child directories and files of the given root directory into a
single tgx file. The destination directory is set with the 'dest' flag,
or is loaded from the config file or set to the current working directory
if both those are unset

The 'output' flag sets the name of the file to write. If unspecified then
the name of the directory is used.

If no path is given as an argument then the current working directory is
used, so if at the root of a mod directory then 'tgxpack pack' will pack
the mod into the default location.

### New
Creates a new subdirectory in the current working directory containing a
blank ModInfo.ini (containing all necessary key names, but with only the
Kohan version key set) and an empty ModChanges.txt file. Takes the directory
name as argument.

### Add
Takes a query string as argument and searches in the file specified by the
'source' option or in the DefaultSource value in the config file. The top
results to the search are printed, and the user can specify which file(s)
to unpack. The specified files are unpacked into the directory rooted at
'dest', or the current working directory if 'dest' is not given.
