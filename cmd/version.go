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

const VERSION = "0.1.0-beta.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the version number",
	Long: `Output information about the program version:
	version number and, if available, git revision.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
