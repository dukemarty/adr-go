/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/dukemarty/adr-go/data"
	"github.com/dukemarty/adr-go/utils"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage user-specific configuration file",
	Long: `Create, change and/or print a user-specific configuration dotfile.
	
This file can point to a default editor to use for ADR editing, and it may
contain the path to a central ADR store (support for this is not fully implemented
yet).`,
	Args: cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		editor, _ := cmd.Flags().GetString("editor")
		store, _ := cmd.Flags().GetString("store")

		logger := utils.SetupLogger(verbose)

		logger.Printf("Command 'config' called with editor='%s', and central adr store at '%s'.\n", editor, store)

		// Basic function
		config, err := data.LoadUserConfiguration()
		if err != nil {
			logger.Printf("Could not load user configuration: %v\n", err)
			fmt.Println("New user configuration is created.")
			config = *data.NewUserConfiguration(editor, store)
		} else {
			if len(editor) > 0 {
				config.Editor = editor
			}
			if len(store) > 0 {
				config.CentralAdrStore = store
			}
		}
		config.Store()

		fmt.Printf("%v\n", config)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringP("editor", "e", "", "path to editor exucutable to use (by default) to open ADRs")
	configCmd.Flags().StringP("store", "s", "", "path to central ADR store")
}
