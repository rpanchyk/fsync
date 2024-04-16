# File Sync Utility

The `fsync` is a file transfer program capable of local update.

## Usage

```shell
Usage:
  fsync [flags] SRC DEST

Flags:
  -h, --help      help for fsync
  -v, --verbose   increase verbosity
```

where
- `SRC` is a source (file or folder) to synchronize.
- `DEST` is a destination folder, where source file or folder should be placed.

## Development

### Update dependencies

```shell
go mod tidy && go mod vendor
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
