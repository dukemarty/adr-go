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

var flagNewStatus = data.AdrStatus("")

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <adr index>",
	Short: "Change one ADR status",
	Long: `Change status of a selected ADR,  may be used interactively.

	The new status can either be provided using the -s/--status flag, or an interactive
	prompt is shown for the user to select the new status.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")

		logger := utils.SetupLogger(verbose)

		adrIdx, err := strconv.Atoi(args[0])

		if err != nil {
			logger.Fatalf("ERROR: provided ADR index must be number, could not be parsed: %s\n", args[0])
		}

		var newStatus string
		// if len(args) > 1 {
		if len(flagNewStatus) > 1 {
			logger.Printf("Command 'status' called for ADR #%d with new status %s.\n", adrIdx, args[1])
			// newStatus = args[1]
			newStatus = flagNewStatus.String()
		} else {
			logger.Printf("Command 'status' called for ADR #%d without new status.\n", adrIdx)
			newStatus = logic.GetStatusInteractively(fmt.Sprintf("ADR #%d()", adrIdx))
		}

		adrFile, err := logic.GetAdrFilePathByIndex(adrIdx, logger)
		if err != nil {

		}
		data.AddStatusEntry(logger, adrFile, newStatus)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().VarP(&flagNewStatus, "status", "s", fmt.Sprintf(`new status to assign, allowed: %v`, data.SupportedStatus))
}
