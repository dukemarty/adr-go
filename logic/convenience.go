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

// Get path to an ADR file bye its index.
//
// Takes the index as parseable string and a logger object.
// Returns either the path if a fitting ADR was found, or
// an error object.
func GetAdrFilePathByIndexString(adrIndexStr string, logger *log.Logger) (string, error) {

	adrIdx, err := strconv.Atoi(adrIndexStr)
	if err != nil {
		logger.Printf("Provided ADR index '%s' must be number, could not be parsed: %v", adrIndexStr, err)
		return "", errors.New(fmt.Sprintf("Provided ADR index '%s' must be number, could not be parsed: %v", adrIndexStr, err))
	}

	return GetAdrFilePathByIndex(adrIdx, logger)
}

func GetAdrFilePathByIndex(adrIndex int, logger *log.Logger) (string, error) {
	am, err := OpenAdrManager(logger)
	if err != nil {
		logger.Printf("Error opening ADR management: %v", err)
		return "", errors.New(fmt.Sprintf("Error opening ADR management: %v", err))
	}

	adrFile, err := am.GetAdrFilenameByIndex(adrIndex, logger)
	if err != nil {
		logger.Printf("Could not find ADR for index %d: %v", adrIndex, err)
		return "", errors.New(fmt.Sprintf("Could not find ADR for index %d: %v", adrIndex, err))
	}

	return filepath.Join(am.Config.Path, adrFile), nil
}

// Get (relative) paths for all ADRs in the repository.
//
// Takes a logger as parameter, returns either a list of string
// (all ADR paths) or an error.
func GetAllAdrFilePaths(logger *log.Logger) ([]string, error) {
	am, err := OpenAdrManager(logger)
	if err != nil {
		logger.Printf("Error opening ADR management: %v", err)
		return nil, errors.New(fmt.Sprintf("Error opening ADR management: %v", err))
	}

	filenames, err := am.GetAllAdrFileNames(logger)
	if err != nil {
		logger.Printf("Error reading all ADR filenames: %v", err)
		return nil, errors.New(fmt.Sprintf("Error reading all ADR filenames: %v", err))
	}

	res := make([]string, 0)
	for _, f := range filenames {
		res = append(res, filepath.Join(am.Config.Path, f))
	}

	return res, nil
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

// Update all ADRs in a repository. "Update" here means to compare
// the filename with the configured format and the actual name of
// the ADR extracted from the file content. After doing this for
// all ADRs and renaming files if necessary, the README is updated.
//
// The function takes a logger as parameter, and returns an error
// if something went wrong.
func UpdateAdrRepository(logger *log.Logger) error {
	am, err := OpenAdrManager(logger)
	if err != nil {
		logger.Printf("Error opening ADR management: %v", err)
		return errors.New(fmt.Sprintf("Error opening ADR management: %v", err))
	}

	filenames, err := am.GetAllAdrFileNames(logger)
	if err != nil {
		logger.Printf("Error reading all ADR filenames: %v", err)
		return errors.New(fmt.Sprintf("Error reading all ADR filenames: %v", err))
	}

	for _, f := range filenames {
		err := am.UpdateFilenameByTitle(f, logger)
		if err != nil {
			logger.Printf("Could not update ADR '%s': %v\n", f, err)
		}
	}

	// update toc
	toc := am.GenerateToc(logger)
	os.WriteFile(filepath.Join(am.Config.Path, "README.md"), []byte(toc), 0644)

	return nil
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
