/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/dukemarty/adr-go/logic"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update ADR",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		initCommon(cmd)

		logger.Println("Command 'update' called.")

		logic.UpdateAdrRepository(logger)

		logger.Println("Filenames updated as required.")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
