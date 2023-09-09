/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	adrexport "github.com/dukemarty/adr-go/export"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export <format>",
	Short: "Export ADRs in another (single file) format",
	Long: fmt.Sprintf(`Exports the ADRs in the selected format. Allowed formats
	are: %v

	In JSON format, only a list of the most important information is returned,
	the other formats contain the complete ADRs.

	The exports are printed on the console, to store directly into a file use
	the -s/--store flag.`, adrexport.SupportedExporters),
	ValidArgs: adrexport.SupportedExporters,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		initCommon(cmd)

		store, _ := cmd.Flags().GetBool("store")

		logger.Printf("Command 'export' called with format '%s', store-to-file=%v.", args[0], store)

		dataPath, data := loadAdrData(logger)

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

	exportCmd.Flags().BoolP("store", "s", false, "store export to file instead of printing to console")
}
