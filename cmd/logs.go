/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dukemarty/adr-go/data"
	"github.com/dukemarty/adr-go/logic"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs <adr index>",
	Short: "List one ADR status logs",
	Long: `This command lists in a table the different status the selected
	ADR has had and the timestamp when the status was reached.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		initCommon(cmd)

		logger.Printf("Command 'logs' called for ADR #%s.\n", args[0])

		adrFile, err := logic.GetAdrFilePathByIndexString(args[0], logger)
		if err != nil {
			logger.Fatalf("Error while trying to get ADR file for index %s: %v", args[0], err)
		}

		status, err := data.ReadStatusEntries(logger, adrFile)
		if err != nil {
			logger.Printf("Error reading status entries: %v\n", err)
			return
		}

		fmt.Printf("ADR #%s: %s\n", args[0], adrFile)
		tbl := tablewriter.NewWriter(os.Stdout)
		tbl.SetHeader([]string{"Date of Change", "Status"})
		for _, st := range status {
			tbl.Append([]string{st.Date, st.Status})
		}

		tbl.Render()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
