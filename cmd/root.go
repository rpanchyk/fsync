package cmd

import (
	"fmt"
	"os"
	"rpanchyk/fsync/internal/service"
	"rpanchyk/fsync/internal/validation"

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
		validator := &validation.ArgsValidator{}
		if err := validator.Validate(args); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if verboseFlag {
			fmt.Println("verbose:", verboseFlag)
			// fmt.Println("delete:", deleteFlag)
			fmt.Println("Non-flag arguments:", args)
			fmt.Println()
		}

		syncer := &service.Syncer{VerboseFlag: verboseFlag}
		if err := syncer.Copy(args[0], args[1]); err != nil {
			fmt.Println("Sync failed")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("Sync finished")
		os.Exit(0)
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
