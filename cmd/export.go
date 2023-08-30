/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	adrexport "github.com/dukemarty/adr-go/export"
	"github.com/dukemarty/adr-go/logic"
	"github.com/dukemarty/adr-go/utils"
	"github.com/spf13/cobra"
)

// TODO: 'export' command not implemented yet!

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export ADR reporter in HTML, CSV, JSON, Markdown",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	ValidArgs: adrexport.SupportedExporters,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		store, _ := cmd.Flags().GetBool("store")

		logger := utils.SetupLogger(verbose)

		logger.Printf("Command 'export' called with format '%s', store-to-file=%v.", args[0], store)

		dataPath, data := loadData(logger)

		exporter, err := adrexport.CreateExporter(logger, args[0])
		if err != nil {
			logger.Printf("Error when creating exporter: %v\n", err)
			return
		}

		exportData := exporter.Export(logger, data, dataPath)
		if store {
			filename := "export." + args[0]
			err := os.WriteFile(filename, []byte(exportData), 0644)
			if err != nil {
				logger.Printf("Error occurred writing the export file: %v\n", err)
			}
		} else {
			fmt.Println(exportData)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")
	exportCmd.Flags().BoolP("store", "s", false, "store export to file instead of printing to console")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loadData(logger *log.Logger) (string, []logic.AdrStatus) {
	am, err := logic.OpenAdrManager(logger)
	if err != nil {
		logger.Printf("Error opening ADR management: %v\n", err)
		return "", []logic.AdrStatus{}
	}

	allAdrs, err := am.GetListOfAllAdrsStatus(logger)
	if err != nil {
		logger.Printf("Error while loading ADR status': %v\n", err)
		return am.Config.Path, []logic.AdrStatus{}
	}
	logger.Printf("Number of parsed and loaded ADRs: %d\n", len(allAdrs))

	return am.Config.Path, allAdrs
}
