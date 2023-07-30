/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/dukemarty/adr-go/data"
	"github.com/dukemarty/adr-go/logic"
	"github.com/dukemarty/adr-go/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs <adr index>",
	Short: "List one ADR status logs",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")

		logger := utils.SetupLogger(verbose)
		logger.Println("Command 'logs' called.")

		adrFile, err := logic.GetAdrFilePathByIndexString(args[0], logger)
		if err != nil {
			logger.Fatalf("Error while trying to get ADR file for index %s: %v", args[0], err)
		}

		status, err := data.ReadStatusEntries(adrFile)
		if err != nil {
			logger.Printf("ERROR: %v\n", err)
			return
		}

		tbl := tablewriter.NewWriter(os.Stdout)
		tbl.SetHeader([]string{"Date of Change", "Status"})
		for _, st := range status {
			tbl.Append([]string{st.Date, st.Status})
		}

		tbl.Render()

		// log.SetFlags(0)
		// log.SetOutput(ioutil.Discard)

		// fmt.Printf("logs called for ADR #%s -> %s\n", args[0], adrFile)
		// logger.Fatalln("Command <logs>: Not implemented yet!")
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
