/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
package logic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/dukemarty/adr-go/data"
	"github.com/dukemarty/adr-go/templates"
	"github.com/dukemarty/adr-go/utils"
)

var configFileName = ".adr.json"

var defaultTemplate = `# {{.NUMBER}}. {{.TITLE}}

Date: {{.DATE}}

## Status

{{.DATE}} proposed

## Context

Context here...

## Decision

Decision here...

## Consequences

Consequences here...`

type AdrManager struct {
	Config data.Configuration
}

// Constructor for a new AdrManager object with a given configuration config.
// The newly constructred AdrManager is returned.
func NewAdrManager(config data.Configuration) *AdrManager {
	am := AdrManager{
		Config: config,
	}

	return &am
}

// Constructor for an AdrManager based on a stored (initialized) ADR setup
// in the current directory.
// Return either the constructed AdrManager, or an error if it could not be opened.
func OpenAdrManager(logger *log.Logger) (*AdrManager, error) {
	config, err := data.LoadConfiguration()
	if err != nil {
		logger.Printf("Could not load ADR configuration, maybe project is not initialized: %v\n", err)
		return nil, errors.New("Could not load ADR configuration.")
	}

	am := AdrManager{
		Config: config,
	}

	return &am, nil
}

// Initialize ADR management in the current directory. Logging is performed
// via the logger provided as parameter, allowing to better control how much
// logging is done.
// Returns an error if for any reason initialization could not be performed,
// in particular if it was already initialized this counts as an error.
func (am AdrManager) Init(logger *log.Logger) error {

	if utils.FileExists(configFileName) {
		return errors.New("ADRs seem to be initialized already, config file '.adr.json' exists!")
	}

	// 1) Create config file
	am.Config.Store(configFileName)

	// 2) Create directory for ADRs
	if _, err := os.Stat(am.Config.Path); os.IsNotExist(err) {
		if err := os.MkdirAll(am.Config.Path, os.ModePerm); err != nil {
			return errors.New(fmt.Sprintf("Error when trying to create directory for adr's: %v", err))
		}
	}

	// 3) Create standard templates
	val, ok := templates.TemplatesLibrary[am.Config.Language]
	if !ok {
		val, _ = templates.TemplatesLibrary["en"]
	}
	pathShortTemplate := filepath.Join(am.Config.Path, "template-short.md")
	pathlongTemplate := filepath.Join(am.Config.Path, "template-long.md")
	errShort := ioutil.WriteFile(pathShortTemplate, []byte(val.Short), 0644)
	if errShort != nil {
		logger.Printf("Error when writing short ADR template: %v\n", errShort)
	}
	errLong := ioutil.WriteFile(pathlongTemplate, []byte(val.Long), 0644)
	if errLong != nil {
		logger.Printf("Error when writing long ADR template: %v\n", errLong)
	}
	if errShort == nil && errLong == nil {
		logger.Println("Successfully wrote short and long ADR templates.")
	}

	return nil
}

func (am AdrManager) AddAdr(title string, logger *log.Logger) (string, error) {
	templateContent := am.loadTemplateOrDefault(filepath.Join(am.Config.Path, am.Config.TemplateName), logger)
	return am.AddAdrWithContent(title, templateContent, logger)
}

func (am AdrManager) AddAdrFromTemplate(title string, templateFile string, logger *log.Logger) (string, error) {
	templContent := am.loadTemplateOrDefault(templateFile, logger)

	return am.AddAdrWithContent(title, templContent, logger)
}

// Add new ADR with the provided title and the also provided content.
// The content is processed as a template, replacing variables. It
// also generates/updates the TOC file.
//
// The replaced variables are '{.NUMBER}', '{.TITLE}' and '{.DATE}'.
func (am AdrManager) AddAdrWithContent(title string, content string, logger *log.Logger) (string, error) {
	newDate := createDateString()
	index := am.getNewIndexString(logger)
	fileName := index + "-" + generateBaseFileName(title) + ".md"

	// 	let newIndex = Utils.getNewIndexString()
	// 	let fileData = raw.replace(/{NUMBER}/g, Utils.getLatestIndex() + 1)
	// 	  .replace(/{TITLE}/g, name)
	// 	  .replace(/{DATE}/g, newDate)
	vars := data.AdrVars{
		NUMBER: index,
		TITLE:  title,
		DATE:   newDate,
	}
	logger.Printf("Identified template variables: %v\n", vars)

	tmpl, err := template.New("adr").Parse(content)
	if err != nil {
		panic(err)
	}

	am.createAdrFile(am.Config.Path, fileName, tmpl, vars)

	toc := am.GenerateToc(logger)
	os.WriteFile(filepath.Join(am.Config.Path, "README.md"), []byte(toc), 0644)

	return fileName, nil
}

func (am AdrManager) EditAdr(adrIndex int) {

}

// Get an ADR's filename for a given index number adrIndex.
// Returns either the found filename, or an error object if it could not find
// the respective ADR.
func (am AdrManager) GetAdrFilenameByIndex(adrIndex int, logger *log.Logger) (string, error) {
	allAdrFiles, err := am.getAdrFiles(logger)
	if err != nil {
		logger.Printf("Could not read any ADRs, in particular not found index %d: %v\n", adrIndex, err)
		return "", err
	}

	for _, filename := range allAdrFiles {
		index, err := am.ExtractAdrIndexFromFile(filename)
		if err != nil {
			continue
		}
		if index == adrIndex {
			return filename, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Could not find ADR with index %d", adrIndex))
}

func (am AdrManager) loadTemplateOrDefault(templateFile string, logger *log.Logger) string {
	rawTemplate, err := os.ReadFile(filepath.Join(am.Config.Path, templateFile))
	var templContent string
	if err != nil {
		logger.Printf("Could not read requested template file %s: %v\n", templateFile, err)
		logger.Println("Use standard template instead!")
		val, _ := templates.TemplatesLibrary["en"]
		templContent = val.Short
	} else {
		templContent = string(rawTemplate)
	}

	return templContent
}

// Generate table of content of all found ADRs and return
// it as a string.
func (am AdrManager) GenerateToc(logger *log.Logger) string {
	var sb strings.Builder

	// header
	sb.WriteString("# Architecture Decision Records\n\n")

	// body
	adrs, err := am.getAdrFiles(logger)
	if err != nil {

	}
	sort.Strings(adrs)
	for _, fn := range adrs {
		adrInfos, err := data.LoadAdrInfo(logger, am.Config.Path, fn)
		if err == nil {
			entry := "\n* [" + strconv.Itoa(adrInfos.Index) + ". " + adrInfos.Title + "](" + adrInfos.RelativePath + ")"
			sb.WriteString(entry)
		}
	}

	// footer
	sb.WriteString("\n")

	return sb.String()
}

type AdrStatus struct {
	Index        int
	Title        string
	LastModified string
	LastStatus   string
}

func (am AdrManager) GetListOfAllAdrsStatus(logger *log.Logger) ([]AdrStatus, error) {
	allAdrFiles, err := am.getAdrFiles(logger)
	if err != nil {
		logger.Printf("Could not load ADR files: %v\n", err)
		return nil, err
	}

	res := make([]AdrStatus, 0)
	for _, filename := range allAdrFiles {
		adrInfos, err := data.LoadAdrInfo(logger, am.Config.Path, filename)
		if err != nil {
			logger.Printf("Error loading basic info for %s\n", filename)
			continue
		}
		status, err := data.ReadStatusEntries(logger, adrInfos.RelativePath)
		if err != nil {
			logger.Printf("Error loading status for %s\n", filename)
			continue
		}
		res = append(res, AdrStatus{Index: adrInfos.Index, Title: adrInfos.Title, LastModified: status[len(status)-1].Date, LastStatus: status[len(status)-1].Status})
	}

	return res, nil
}

func (am AdrManager) extractAdrStatusFromFile(filename string) (AdrStatus, error) {
	// load ADR

	// if okay:

	testStatus := AdrStatus{Index: 1, Title: "Blabla Dings", LastModified: "1979-10-13", LastStatus: "1979-10-13 Birth"}

	return testStatus, nil
}

func (am AdrManager) ExtractAdrIndexFromFile(filename string) (int, error) {
	indexPart := filename[:am.Config.Digits]
	index, err := strconv.Atoi(indexPart)
	if err != nil {
		log.Printf("Could not parse '%s' as index: %v\n", indexPart, err)
		return -1, err
	}

	return index, nil
}

func (am AdrManager) createAdrFile(adrDirectory string, filename string, content *template.Template, data data.AdrVars) {

	f, err := os.Create(filepath.Join(adrDirectory, filename))

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// _, err2 := f.WriteString(content)
	err2 := content.Execute(f, data)

	if err2 != nil {
		log.Fatal(err2)
	}

	// 	let toc = generate('toc', { output: false })
	// 	fs.writeFileSync(savePath + 'README.md', toc + '\n')
}

func createDateString() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

func (am AdrManager) getNewIndexString(logger *log.Logger) string {
	lastIndex, err := am.getLatestIndex(logger)
	if err != nil {
		return am.createIndexByNumber(1, logger)
	}
	lastIndex = lastIndex + 1
	return am.createIndexByNumber(lastIndex, logger)
}

func (am AdrManager) getLatestIndex(logger *log.Logger) (int, error) {
	files, err := am.getAdrFiles(logger)

	if err != nil {
		logger.Printf("Error when trying to load existing ADR files: %v\n", err)
		return 0, err
	}

	if len(files) == 0 {
		logger.Println("Found no ADR files.")
		return 0, errors.New("No ADR files found.")
	}

	return am.getMaxIndex(files, logger), nil
}

func (am AdrManager) getMaxIndex(filenames []string, logger *log.Logger) int {
	maxNumber := 0

	for _, file := range filenames {
		logger.Printf("Trying to extract index from file of name '%s'\n", file)
		indexPart := file[:am.Config.Digits]
		logger.Printf("Found index parts: %s\n", indexPart)
		index, err := strconv.Atoi(indexPart)

		if err == nil && index > maxNumber {
			maxNumber = index
		}
	}

	return maxNumber
}

func (am AdrManager) getAdrFiles(logger *log.Logger) ([]string, error) {
	files, err := ioutil.ReadDir(am.Config.Path)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for _, file := range files {
		logger.Printf("Analyzing file %v, Name='%s', IsDir=%v with file extension='%s'\n", file, file.Name(), file.IsDir(), filepath.Ext(file.Name()))
		if !file.IsDir() && file.Name() != "README.md" && !strings.HasPrefix(file.Name(), "template-") && file.Name() != am.Config.TemplateName && filepath.Ext(file.Name()) == ".md" {
			res = append(res, file.Name())
		}
	}

	return res, nil
}

func (am AdrManager) createIndexByNumber(number int, logger *log.Logger) string {
	s := fmt.Sprintf("%020d", number)
	logger.Printf("Trying to create index by number: %s", s)
	return am.Config.Prefix + s[len(s)-am.Config.Digits:]
}

func generateBaseFileName(title string) string {
	mToHyphen1 := regexp.MustCompile(`[\s_-]+`)
	mRemove := regexp.MustCompile(`[#,.]+`)
	mToHyphen2 := regexp.MustCompile(`[:?]+`)

	filename := strings.Trim(strings.ToLower(title), " -")
	filename = mToHyphen1.ReplaceAllString(filename, "-")
	filename = mRemove.ReplaceAllString(filename, "")
	filename = mToHyphen2.ReplaceAllString(filename, "-")

	return filename
}
