/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/dukemarty/adr-go/utils"
	"github.com/spf13/cobra"
)

// TODO: 'status' command not implemented yet!

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <adr index> [new status]",
	Short: "Change one ADR status",
	Long: `Change status of a selected ADR,  may be used interactively.

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")

		logger := utils.SetupLogger(verbose)

		logger.Println("Command 'status' called.")

		adrIdx, err := strconv.Atoi(args[0])

		if err != nil {
			logger.Fatalf("ERROR: provided ADR index must be number, could not be parsed: %s\n", args[0])
		}

		var newStatus string
		if len(args) > 1 {
			newStatus = args[1]
		} else {
			newStatus = utils.GetStatusInteractively(fmt.Sprintf("ADR #%d()", adrIdx))
		}

		fmt.Printf("status called for ADR #%d with new status='%s'\n", adrIdx, newStatus)
		logger.Fatalln("Command <status>: Not implemented yet!")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
