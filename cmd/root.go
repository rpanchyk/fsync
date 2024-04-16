package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"rpanchyk/fsync/internal/service"
	"strings"

	"github.com/spf13/cobra"
)

var verboseFlag bool

//var deleteFlag  bool

var rootCmd = &cobra.Command{
	Use:   "fsync",
	Short: "fsync is a file transfer program capable of local update. \nUsage: fsync [OPTION]... SRC DEST",
	Long: `fsync is a file transfer program capable of local update. \nUsage: fsync [OPTION]... SRC DEST

Attention! Use this tool on your own risk.`,
	Run: func(cmd *cobra.Command, args []string) {
		if verboseFlag {
			fmt.Println("verbose:", verboseFlag)
			// fmt.Println("delete:", deleteFlag)
			fmt.Println("Non-flag arguments:", args)
			fmt.Println()
		}

		if len(args) != 2 {
			fmt.Println("Invalid args")
			os.Exit(1)
		}

		currDir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			if verboseFlag {
				fmt.Println("Current folder:", currDir)
			}
		}

		src, dst := args[0], args[1]
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
			os.Exit(1)
		}

		syncer := &service.Syncer{VerboseFlag: verboseFlag}
		err = syncer.Copy(srcPath, dstPath)
		if err != nil {
			fmt.Println("Cannot synchronize source", srcPath, "with destination", dstPath)
			os.Exit(1)
		}
		fmt.Println("Sync finished")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "increase verbosity")
}
