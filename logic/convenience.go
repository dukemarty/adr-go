/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
package logic

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func GetAdrFilenamesFiltered(keywords []string, caseSensitive bool, logger *log.Logger) ([]string, error) {
	am, err := OpenAdrManager(logger)
	if err != nil {
		logger.Printf("Error opening ADR management: %v", err)
		return nil, errors.New(fmt.Sprintf("Error opening ADR management: %v", err))
	}

	regexes, err := compileKeywordsToRegexes(keywords, caseSensitive)
	if err != nil {
		logger.Printf("Error compiling regexes for keyword search: %v\n", err)
	}
	allAdrFiles, err := am.GetAllAdrFileNames(logger)
	if err != nil {
		logger.Printf("Error loading ADR file names: %v\n", err)
	}

	res := make([]string, 0)
FILELOOP:
	for _, adrFile := range allAdrFiles {
		rawContent, err := os.ReadFile(filepath.Join(am.Config.Path, adrFile))
		if err != nil {
			logger.Printf("Error reading ADR from '%s': %v\n", adrFile, err)
			return nil, err
		}
		content := string(rawContent)
		for _, r := range regexes {
			if !r.MatchString(content) {
				continue FILELOOP
			}
		}

		res = append(res, adrFile)
	}

	return res, nil
}

func GetStatusFromListOfAdrFiles(files []string, logger *log.Logger) ([]AdrStatus, error) {
	am, err := OpenAdrManager(logger)
	if err != nil {
		logger.Printf("Error opening ADR management: %v", err)
		return nil, errors.New(fmt.Sprintf("Error opening ADR management: %v", err))
	}

	res, err := am.GetStatusFromListOfAdrFiles(files, logger)

	return res, err
}

func compileKeywordsToRegexes(keywords []string, caseSensitive bool) ([]*regexp.Regexp, error) {
	prefix := ""
	if !caseSensitive {
		prefix = "(?i)"
	}

	res := make([]*regexp.Regexp, 0)
	for _, kw := range keywords {
		r, err := regexp.Compile(prefix + kw)
		res = append(res, r)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
