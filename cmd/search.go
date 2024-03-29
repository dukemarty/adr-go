/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/dukemarty/adr-go/logic"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <keywords>+",
	Short: "Search ADRs by keywords",
	Long: `Filter all ADRs by the presence of the provided keywords in their content,
	and print the table of all found ADRs and their status.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		initCommon(cmd)

		caseSensitive, _ := cmd.Flags().GetBool("casesensitive")

		if caseSensitive {
			logger.Printf("Command 'search' called case-sensitive with these keywords: %v\n", args)
		} else {
			logger.Printf("Command 'search' called case-insensitive with these keywords: %v\n", args)
		}

		foundFiles, err := logic.GetAdrFilenamesFiltered(args, caseSensitive, logger)
		if err != nil {
			logger.Printf("Error occured while filtering files: %v\n", err)
			return
		}

		statuss, err := logic.GetStatusFromListOfAdrFiles(foundFiles, logger)
		if err != nil {
			logger.Printf("Error loading status' from ADR files: %v\n", err)
			return
		}

		tbl := tablewriter.NewWriter(os.Stdout)
		tbl.SetAutoWrapText(false)
		tbl.SetHeader([]string{"Filename", "Last status"})
		for _, adrst := range statuss {
			tbl.Append([]string{adrst.Filename, adrst.LastStatus})
		}
		tbl.Render()

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().BoolP("casesensitive", "c", false, "flag to activate case-sensitive search")
}
