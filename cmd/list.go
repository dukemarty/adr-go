/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dukemarty/adr-go/logic"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var statusColors = map[string]tablewriter.Colors{
	"PROPOSED":   {tablewriter.Normal, tablewriter.FgHiWhiteColor},
	"ACCEPTED":   {tablewriter.Normal, tablewriter.FgHiCyanColor},
	"DONE":       {tablewriter.Normal, tablewriter.FgHiGreenColor},
	"DEPRECATED": {tablewriter.Normal, tablewriter.FgHiRedColor},
	"SUPERSEDED": {tablewriter.Normal, tablewriter.FgHiYellowColor},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ADRs",
	Long: `Print a table of all ADRs (order by index) containing
	their index, name, current status, and timestamp of last status
	change.`,
	Args: cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		initCommon(cmd)

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
			row := []string{fmt.Sprintf(fmtString, adrst.Index), adrst.Title, adrst.LastModified, adrst.LastStatus}
			statusColor := tablewriter.Colors{}
			if val, present := statusColors[strings.ToUpper(adrst.LastStatus)]; present {
				statusColor = val
			}
			colors := []tablewriter.Colors{{}, {}, {}, statusColor}
			tbl.Rich(row, colors)
		}
		tbl.Render()

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
