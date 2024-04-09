package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	verboseFlag bool
	// deleteFlag  bool
)

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\nfsync is a file transfer program capable of local update. \nUsage: fsync [OPTION]... SRC DEST\n\n")
		flag.PrintDefaults()
	}

	flag.BoolVar(&verboseFlag, "verbose", false, "increase verbosity")
	// flag.BoolVar(&deleteFlag, "delete", false, "delete extraneous files from dest dirs")
	flag.Parse()

	if verboseFlag {
		fmt.Println("verbose:", verboseFlag)
		// fmt.Println("delete:", deleteFlag)
		fmt.Println("Non-flag arguments:", flag.Args())
		fmt.Println()
	}

	if len(flag.Args()) != 2 {
		fmt.Println("Invalid args")
		flag.Usage()
		return
	}

	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		if verboseFlag {
			fmt.Println("Current folder:", currDir)
		}
	}

	src, dst := flag.Args()[0], flag.Args()[1]
	if verboseFlag {
		fmt.Println("Input source:", src)
		fmt.Println("Input destination:", dst)
	}

	srcPath := src
	if !filepath.IsAbs(srcPath) {
		srcPath = filepath.Join(currDir, srcPath)
	}
	dstPath := dst
	if !filepath.IsAbs(dstPath) {
		dstPath = filepath.Join(currDir, dstPath)
	}
	if verboseFlag {
		fmt.Println("Normalized source:", srcPath)
		fmt.Println("Normalized destination:", dstPath)
	}
	if strings.HasPrefix(dstPath, srcPath) {
		fmt.Println("Cannot synchronize because destination", dstPath, "is sub-folder of source", srcPath)
		return
	}

	err = copy(srcPath, dstPath, verboseFlag)
	if err != nil {
		fmt.Println("Cannot synchronize source", srcPath, "with destination", dstPath)
		return
	}
	fmt.Println("Sync finished")
}

func copy(src, dst string, verboseFlag bool) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return errors.New(fmt.Sprint("Source doesn't exist:", src))
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err = os.MkdirAll(dst, os.ModeDir); err != nil {
			return errors.New(fmt.Sprint("Cannot create destination folder:", dst))
		} else {
			if verboseFlag {
				fmt.Println("Created destination folder:", dst)
			}
		}
	}

	srcIsDir, err := isDirectory(src)
	if err != nil {
		return errors.New(fmt.Sprint("Cannot analyze source:", src))
	}

	if srcIsDir {
		entries, err := os.ReadDir(src)
		if err != nil {
			return errors.New(fmt.Sprint("Cannot get entries of source folder:", src))
		}

		for _, entry := range entries {
			entryInfo, err := entry.Info()
			if err != nil {
				return errors.New(fmt.Sprint("Cannot get entry info:", entry))
			}

			srcPath := filepath.Join(src, entryInfo.Name())

			dstPath := dst
			if entryInfo.IsDir() {
				dstPath = filepath.Join(dst, entryInfo.Name())
			}

			err = copy(srcPath, dstPath, verboseFlag)
			if err != nil {
				return err
			}
		}
	} else {
		dstPath := filepath.Join(dst, filepath.Base(src))
		nBytes, err := copyFile(src, dstPath)
		if err != nil {
			return errors.New(fmt.Sprint("Cannot copy file:", src, "to", dstPath))
		}
		if verboseFlag {
			fmt.Println("Copied file", src, "of", nBytes, "bytes")
		}
	}

	return nil
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func copyFile(src, dst string) (int64, error) {
	sourceFileInfo, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileInfo.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
