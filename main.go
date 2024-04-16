package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"rpanchyk/fsync/internal/service"
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

	syncer := &service.Syncer{VerboseFlag: verboseFlag}
	err = syncer.Copy(srcPath, dstPath)
	if err != nil {
		fmt.Println("Cannot synchronize source", srcPath, "with destination", dstPath)
		return
	}
	fmt.Println("Sync finished")
}
