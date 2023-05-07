/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/dukemarty/adr-go/logic"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <adr title>",
	Short: "Create new ADR",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		template, _ := cmd.Flags().GetString("template")

		fmt.Printf("new called with title '%s', explicit template?=%v ('%s')\n", args[0], len(template) > 0, template)

		am, err := logic.OpenAdrManager()
		if err != nil {
			log.Fatalf("Error opening ADR management: %v\n", err)
		}

		var adrFile string
		if len(template) > 0 {
			adrFile, err = am.AddAdrFromTemplate(args[0], template)
		} else {
			adrFile, err = am.AddAdr(args[0])
		}
		if err != nil {
			log.Fatalf("Error when creating new ADR: %v\n", err)
		}

		log.Printf("Created new ADR as %s\n", adrFile)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.
	newCmd.Flags().StringP("template", "t", "", "template to use for the new ADR")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
