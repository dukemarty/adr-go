/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dukemarty/adr-go/documents"
	"github.com/spf13/cobra"
)

var availableDocuments = []string{"Changelog", "License"}

// changelogCmd represents the changelog command
var showCmd = &cobra.Command{
	Use:   "show <document>",
	Short: "Show or store an embedded document.",
	Long: fmt.Sprintf(`adr-go contains a number of documents embedded, e.g. its license.
	Using this command, those documents can be printed to stdout, or all of the documents
	can be stored to a file each.  
	
	The currently available documents are: %v

	So if the <document> argument is provided, it must be one of those values.

	By providing the -c/--create flag, the selected document is stored with its
	standard file name.

	If no document type is provided, _with_ the create flag all files are stored.
	Without the flag, an interactive prompt is show to the user for selecting
	which document he wants to see.`, availableDocuments),
	ValidArgs: availableDocuments,
	Args:      cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),

	Run: func(cmd *cobra.Command, args []string) {
		initCommon(cmd)

		create, _ := cmd.Flags().GetBool("create")

		if create {
			logger.Println("Command 'show' called with 'create' flag.")
			if len(args) == 0 {
				storeAllDocuments(logger)
			} else {
				storeSingleDocument(logger, args[0])
			}
		} else {
			var doc string
			if len(args) == 0 {
				logger.Println("Command 'show' called, without a document type selected.")
				prompt := &survey.Select{
					Message: "What document shall be shown:",
					Options: availableDocuments,
				}
				survey.AskOne(prompt, &doc)
			} else {
				logger.Printf("Command 'show' called, with document type '%s'.\n", args[0])
				doc = args[0]
			}

			switch strings.ToUpper(doc) {
			case "CHANGELOG":
				fmt.Println(documents.Changelog)
			case "LICENSE":
				fmt.Println(documents.License)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().BoolP("create", "c", false, "create files for all documents")
}

func storeAllDocuments(logger *log.Logger) {
	for _, df := range documents.Docs {
		err := os.WriteFile(df.Filename, []byte(df.Content), 0644)
		if err != nil {
			logger.Printf("Problem writing %s: %v\n", df.Filename, err)
		} else {
			logger.Printf("Wrote %s.\n", df.Filename)
		}
	}
}

func storeSingleDocument(logger *log.Logger, docType string) {
	var filename string
	var content string
	switch strings.ToUpper(docType) {
	case "CHANGELOG":
		filename = "CHANGELOG.md"
		content = documents.Changelog
	case "LICENSE":
		filename = "LICENSE"
		content = documents.License
	}
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		logger.Printf("Problem writing %s: %v\n", filename, err)
	} else {
		logger.Printf("Wrote %s.\n", filename)
	}
}
