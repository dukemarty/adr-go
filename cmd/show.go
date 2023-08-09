/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dukemarty/adr-go/documents"
	"github.com/dukemarty/adr-go/utils"
	"github.com/spf13/cobra"
)

// changelogCmd represents the changelog command
var showCmd = &cobra.Command{
	Use:   "show <document>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		verbose, _ := cmd.Flags().GetBool("verbose")
		create, _ := cmd.Flags().GetBool("create")

		logger := utils.SetupLogger(verbose)

		if create {
			logger.Println("Command 'show' called,  'create' flag.")
			for _, df := range documents.Docs {
				err := os.WriteFile(df.Filename, []byte(df.Content), 0644)
				if err != nil {
					logger.Printf("Problem writing %s: %v\n", df.Filename, err)
				} else {
					logger.Printf("Wrote %s.\n", df.Filename)
				}
			}
		} else {
			var doc string
			if len(args) == 0 {
				logger.Println("Command 'show' called, without a document type selected.")
				prompt := &survey.Select{
					Message: "What document shall be shown:",
					Options: []string{"Changelog", "License"},
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

	// Here you will define your flags and configuration settings.
	showCmd.Flags().BoolP("create", "c", false, "create files for all documents")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changelogCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changelogCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
