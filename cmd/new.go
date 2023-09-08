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
		initCommon(cmd)

		template, _ := cmd.Flags().GetString("template")
		editor, _ := cmd.Flags().GetString("editor")

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

	newCmd.Flags().StringP("template", "t", "", "template file to use for the new ADR (located in ADR folder)")
	newCmd.Flags().StringP("editor", "e", "", "path to editor executable for opening the ADR")
}
