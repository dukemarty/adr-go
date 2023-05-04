/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/dukemarty/adr-go/data"
	"github.com/dukemarty/adr-go/logic"

	"github.com/spf13/cobra"
)

var firstAdr = `# {{NUMBER}}. Record architecture decisions

Date: {{DATE}}

## Status

{{DATE}} Accepted

## Context

We need to record the architectural decisions made on this project.

## Decision

We will use Architecture Decision Records, as described by Michael Nygard in this article: http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions

## Consequences

See Michael Nygard's article, linked above.`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init ADR repository",
	Long: `Initialize ADR repository.
	
	This involves setting up a folder for the ADRs, adding a configuration
	file for the ADR tool, and adding initial ADRs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called, todos:")

		// 1) Create config file with the provided flags
		path, _ := cmd.Flags().GetString("path")
		lang, _ := cmd.Flags().GetString("lang")
		prefix, _ := cmd.Flags().GetString("prefix")
		digits, _ := cmd.Flags().GetInt("digits")
		newConfig := data.NewConfiguration(lang, path, prefix, digits)
		newConfig.Store(".adr.json")

		// 2) Create adr directory
		// if err := os.MkdirAll(path, os.ModePerm); err != nil {
		// 	log.Fatalf("Error when trying to create directory for adr's: %v", err)
		// }
		am := logic.NewAdrManager(*newConfig)
		am.Init()

		addFirst, _ := cmd.Flags().GetBool("addfirst")
		if addFirst {
			am.AddAdrWithContent("Record architecture decisions", firstAdr)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.
	initCmd.Flags().IntP("digits", "d", 4, "Number of digits for ADR numbering")
	initCmd.Flags().StringP("path", "p", "docs/adr/", "Path to directory where ADRs are stored")
	initCmd.Flags().StringP("prefix", "x", "", "Prefix for ADR numbers")
	initCmd.Flags().BoolP("addfirst", "a", true, "add initial adr about using adr's")
	initCmd.Flags().StringP("lang", "l", "en", "Language used, stored in config file")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
