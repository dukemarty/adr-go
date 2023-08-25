/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/dukemarty/adr-go/data"
	"github.com/dukemarty/adr-go/logic"
	"github.com/dukemarty/adr-go/utils"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <adr index> [new status]",
	Short: "Change one ADR status",
	Long: fmt.Sprintf(`Change status of a selected ADR,  may be used interactively.

	If the new status is provided as argument, any value is allowed; but
	the values used by default (and available via the interactive prompt)
	are: %v`, utils.SupportedStatus),
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")

		logger := utils.SetupLogger(verbose)

		adrIdx, err := strconv.Atoi(args[0])

		if err != nil {
			logger.Fatalf("ERROR: provided ADR index must be number, could not be parsed: %s\n", args[0])
		}

		var newStatus string
		if len(args) > 1 {
			logger.Printf("Command 'status' called for ADR #%d with new status %s.\n", adrIdx, args[1])
			newStatus = args[1]
		} else {
			logger.Printf("Command 'status' called for ADR #%d without new status.\n", adrIdx)
			newStatus = utils.GetStatusInteractively(fmt.Sprintf("ADR #%d()", adrIdx))
		}

		adrFile, err := logic.GetAdrFilePathByIndex(adrIdx, logger)
		if err != nil {

		}
		data.AddStatusEntry(logger, adrFile, newStatus)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
