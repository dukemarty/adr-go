package logic

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
)

func GetAdrFilePathByIndexString(adrIndexStr string, logger *log.Logger) (string, error) {

	adrIdx, err := strconv.Atoi(adrIndexStr)
	if err != nil {
		logger.Printf("Provided ADR index '%s' must be number, could not be parsed: %v", adrIndexStr, err)
		return "", errors.New(fmt.Sprintf("Provided ADR index '%s' must be number, could not be parsed: %v", adrIndexStr, err))
	}

	am, err := OpenAdrManager(logger)
	if err != nil {
		logger.Printf("Error opening ADR management: %v", err)
		return "", errors.New(fmt.Sprintf("Error opening ADR management: %v", err))
	}

	adrFile, err := am.GetAdrFilenameByIndex(adrIdx, logger)
	if err != nil {
		logger.Printf("Could not find ADR for index %d: %v", adrIdx, err)
		return "", errors.New(fmt.Sprintf("Could not find ADR for index %d: %v", adrIdx, err))
	}

	return filepath.Join(am.Config.Path, adrFile), nil
}
