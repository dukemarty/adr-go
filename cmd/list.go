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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		logger.Printf("Number of parsed ADRs: %d\n", len(allAdrs))

		tbl := tablewriter.NewWriter(os.Stdout)
		tbl.SetAutoWrapText(false)
		tbl.SetHeader([]string{"Decision", "Last modified date", "Last status"})
		for _, adrst := range allAdrs {
			tbl.Append([]string{fmt.Sprintf("%d %s", adrst.Index, adrst.Title), adrst.LastModified, adrst.LastStatus})
		}
		tbl.Render()

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
