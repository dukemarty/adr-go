/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/dukemarty/adr-go/data"
	"github.com/dukemarty/adr-go/logic"
	"github.com/dukemarty/adr-go/utils"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit <adr index>",
	Short: "Open ADR in editor",
	Long: `Open the selected ADR in an editor, which can either be provided
	on command line, or the default editor defined in the project configuration
	is used.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		initCommon(cmd)

		editor, _ := cmd.Flags().GetString("editor")

		logger.Printf("Command 'edit' called for ADR with index %s\n", args[0])

		adrFile, err := logic.GetAdrFilePathByIndexString(args[0], logger)
		if err != nil {
			logger.Fatalf("Error while trying to get ADR file for index %s: %v", args[0], err)
		}
		logger.Printf("Found file to edit: %s\n", adrFile)

		utils.EditFile(adrFile, editor, data.LoadEditor(logger), logger)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringP("editor", "e", "", "Path to editor executable for opening the ADR")
}
