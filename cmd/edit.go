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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		editor, _ := cmd.Flags().GetString("editor")

		logger := utils.SetupLogger(verbose)
		logger.Println("Command 'edit' called.")
		logger.Printf("... for ADR with index %s\n", args[0])

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

	// Here you will define your flags and configuration settings.
	editCmd.Flags().StringP("editor", "e", "", "Path to editor executable for opening the ADR")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
