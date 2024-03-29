/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
package cmd

import (
	"log"

	"github.com/dukemarty/adr-go/data"
	"github.com/dukemarty/adr-go/logic"

	"github.com/spf13/cobra"
)

var firstAdr = `# {{.NUMBER}}. Record architecture decisions

Date: {{.DATE}}

## Status

{{.DATE}} Accepted

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
	Args: cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		initCommon(cmd)

		path, _ := cmd.Flags().GetString("path")
		lang, _ := cmd.Flags().GetString("lang")
		prefix, _ := cmd.Flags().GetString("prefix")
		digits, _ := cmd.Flags().GetInt("digits")
		template, _ := cmd.Flags().GetString("template")
		newConfig := data.NewConfiguration(lang, path, prefix, digits, template)

		logger.Printf("Command 'init' called with: %+v.\n", *newConfig)

		// 1) Create config file and adr directory with standard templates
		am := logic.NewAdrManager(*newConfig)
		err := am.Init(logger)

		if err != nil {
			logger.Fatalf("Could not initialize ADRs: %v", err)
		}
		logger.Println("ADRs initialized.")

		// 2) Create initial ADR
		addFirst, _ := cmd.Flags().GetBool("addfirst")
		if addFirst {
			am.AddAdrWithContent("Record architecture decisions", firstAdr, logger)
			log.Println("Initial ADR created.")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().IntP("digits", "d", 4, "Number of digits for ADR numbering")
	initCmd.Flags().StringP("path", "p", "docs/adr/", "Path to directory where ADRs are stored")
	initCmd.Flags().StringP("prefix", "x", "", "Prefix for ADR numbers")
	initCmd.Flags().BoolP("addfirst", "a", true, "add initial adr about using adr's")
	initCmd.Flags().StringP("lang", "l", "en", "Language used, stored in config file")
	initCmd.Flags().StringP("template", "t", "template-short.md", "template to use for new ADRs")
}
