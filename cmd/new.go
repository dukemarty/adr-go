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

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <adr title>",
	Short: "Create new ADR",
	Long: `Create a new ADR with a given title. The new ADR is automatically numbered,
and a template file (either standard or a selected template), and then opened in
an editor.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		verbose, _ := cmd.Flags().GetBool("verbose")
		template, _ := cmd.Flags().GetString("template")
		editor, _ := cmd.Flags().GetString("editor")

		logger := utils.SetupLogger(verbose)

		logger.Printf("Command 'new' called, with title '%s', explicit template?=%v ('%s')\n", args[0], len(template) > 0, template)

		am, err := logic.OpenAdrManager(logger)
		if err != nil {
			logger.Fatalf("Error opening ADR management: %v\n", err)
		}

		var adrFile string
		if len(template) > 0 {
			adrFile, err = am.AddAdrFromTemplate(args[0], template, logger)
		} else {
			adrFile, err = am.AddAdr(args[0], logger)
		}
		if err != nil {
			logger.Fatalf("Error when creating new ADR: %v\n", err)
		}
		logger.Printf("Created new ADR as %s\n", adrFile)

		utils.EditFile(adrFile, editor, data.LoadEditor(logger), logger)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.
	newCmd.Flags().StringP("template", "t", "", "template file to use for the new ADR (located in ADR folder)")
	newCmd.Flags().StringP("editor", "e", "", "Path to editor executable for opening the ADR")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
