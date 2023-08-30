/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/dukemarty/adr-go/utils"
	"github.com/spf13/cobra"
)

const VERSION = "1.1.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the version number",
	Long: `Output information about the program version:
	version number and, if available, git revision.`,
	Args: cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")

		logger := utils.SetupLogger(verbose)
		logger.Println("Command 'version' called.")

		revisionInfo := ""
		if bi, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range bi.Settings {
				if setting.Key == "vcs.revision" {
					revisionInfo = fmt.Sprintf(" (git revision: %s)", setting.Value)
				}
			}
		}

		fmt.Printf("%s%s\n", VERSION, revisionInfo)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
