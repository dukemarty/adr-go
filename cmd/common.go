package cmd

import (
	"log"

	"github.com/dukemarty/adr-go/logic"
	"github.com/dukemarty/adr-go/utils"
	"github.com/spf13/cobra"
)

var logger *log.Logger

func initCommon(cmd *cobra.Command) {
	verbose, _ := cmd.Flags().GetBool("verbose")

	logger = utils.SetupLogger(verbose)
}

func loadAdrData(logger *log.Logger) (string, []logic.AdrStatus) {
	am, err := logic.OpenAdrManager(logger)
	if err != nil {
		logger.Printf("Error opening ADR management: %v\n", err)
		return "", []logic.AdrStatus{}
	}

	allAdrs, err := am.GetListOfAllAdrsStatus(logger)
	if err != nil {
		logger.Printf("Error while loading ADR status': %v\n", err)
		return am.Config.Path, []logic.AdrStatus{}
	}
	logger.Printf("Number of parsed and loaded ADRs: %d\n", len(allAdrs))

	return am.Config.Path, allAdrs
}
