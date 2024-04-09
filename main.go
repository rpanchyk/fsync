package main

import (
	"flag"
	"fmt"
)

var (
	deleteFlag bool
)

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\nfsync is a file transfer program capable of local update. \nUsage: fsync [OPTION]... SRC [SRC]... DEST\n\n")
		flag.PrintDefaults()
	}

	flag.BoolVar(&deleteFlag, "delete", false, "delete extraneous files from dest dirs")
	flag.Parse()

	fmt.Println("deleteFlag value is: ", deleteFlag)
	fmt.Println("Non-flag arguments:", flag.Args())
}
