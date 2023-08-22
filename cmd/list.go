/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dukemarty/adr-go/logic"
	"github.com/dukemarty/adr-go/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ADRs",
	Long: `Print a table of all ADRs (order by index) containing
	their index, name, current status, and timestamp of last status
	change.`,
	Args: cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")

		logger := utils.SetupLogger(verbose)

		logger.Println("Command 'list' called.")

		am, err := logic.OpenAdrManager(logger)
		if err != nil {
			logger.Fatalf("Error opening ADR management: %v\n", err)
		}

		allAdrs, err := am.GetListOfAllAdrsStatus(logger)
		if err != nil {
			logger.Printf("Error while loading ADR status': %v\n", err)
		}
		logger.Printf("Number of parsed and loaded ADRs: %d\n", len(allAdrs))

		tbl := tablewriter.NewWriter(os.Stdout)
		tbl.SetAutoWrapText(false)
		tbl.SetHeader([]string{"Index", "Decision", "Last modified date", "Last status"})
		// to format index with fitting number of leading zeros, used data from config
		fmtString := fmt.Sprintf("%%0%dd", am.Config.Digits)
		for _, adrst := range allAdrs {
			tbl.Append([]string{fmt.Sprintf(fmtString, adrst.Index), adrst.Title, adrst.LastModified, adrst.LastStatus})
		}
		tbl.Render()

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
