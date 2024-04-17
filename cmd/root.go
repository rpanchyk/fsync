package cmd

import (
	"fmt"
	"os"

	"github.com/rpanchyk/fsync/internal/checksum"
	"github.com/rpanchyk/fsync/internal/syncer"
	"github.com/spf13/cobra"
)

var verboseFlag bool

//var deleteFlag  bool

var rootCmd = &cobra.Command{
	Use:   "fsync [flags] SRC DEST",
	Short: "fsync is a file transfer program capable of local update.",
	Long: `fsync is a file transfer program capable of local update.
Attention! Use this tool on your own risk! Author is not responsible of synced files.`,
	Args: cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: remove
		// validator := &validation.ArgsValidator{}
		// if err := validator.Validate(args); err != nil {
		// 	fmt.Println(err.Error())
		// 	os.Exit(1)
		// }

		if verboseFlag {
			fmt.Println("verbose:", verboseFlag)
			// fmt.Println("delete:", deleteFlag)
			fmt.Println("Non-flag arguments:", args)
			fmt.Println()
		}

		syncer := &syncer.Syncer{
			VerboseFlag:      verboseFlag,
			Source:           args[0],
			Destination:      args[1],
			ChecksumVerifier: checksum.NewVerifier(checksum.MD5),
		}
		if err := syncer.Sync(); err != nil {
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
