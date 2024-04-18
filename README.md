# File Sync Utility

The `fsync` is a file transfer program capable of local update.

## Usage

```help
Usage:
  fsync [flags] SRC DEST

Flags:
  -d, --delete    delete extraneous files from dest dirs
  -h, --help      help for fsync
  -v, --verbose   increase verbosity
```

where

- `SRC` is a source (file or folder) to synchronize.
- `DEST` is a destination folder, where source file or folder should be placed.

## Why

There is a great tool named `rsync`, but it might be not available (for example, Windows doesn't have it).
For sure, [cwRsync](https://itefix.net/cwrsync) can be installed and happy used.
But I'd like to have owns (even if it is much humbler featured).

## How it works

1. Validate `SRC` and `DEST` input arguments.
1. Check file in `DEST` folder:
   1. If file already exists - compare checksums of files and skip sync if equal.
   1. If file not exists:
      1. Copy `SRC` file to `DEST` temporary file.
      1. Compare checksums of files and fail if they differ.
      1. Rename `DEST` temporary file to `DEST` file.

## Development

### Update dependencies

```shell
go mod tidy
go mod vendor
```

### Build

To make a cross-build, please see available platforms:

```shell
go tool dist list
```

For example, for linux run this command to create a binary file for `linux/amd64` architecture:

```shell
GOOS=linux GOARCH=amd64 go build
```

For batch build use [Makefile](Makefile) and run:

```shell
make build
```

It will create `builds` directory with archived binary files according to preconfigured set of platforms.

# Disclaimer

The software is provided "as is", without warranty of any kind, express or
implied, including but not limited to the warranties of merchantability,
fitness for a particular purpose and noninfringement. in no event shall the
authors or copyright holders be liable for any claim, damages or other
liability, whether in an action of contract, tort or otherwise, arising from,
out of or in connection with the software or the use or other dealings in the
software.

# Contribution

If you have any ideas or inspiration for contributing the project,
please create an [issue](https://github.com/rpanchyk/fsync/issues/new)
or [pull request](https://github.com/rpanchyk/fsync/pulls).
