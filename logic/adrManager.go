package logic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/dukemarty/adr-go/data"
	"github.com/dukemarty/adr-go/templates"
	"github.com/dukemarty/adr-go/utils"
)

var configFilenName = ".adr.json"

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

func NewAdrManager(config data.Configuration) *AdrManager {
	am := AdrManager{
		Config: config,
	}

	return &am
}

func OpenAdrManager() (*AdrManager, error) {
	config, err := data.LoadConfiguration()
	if err != nil {
		log.Printf("Could not load ADR configuration, maybe project is not initialized: %v\n", err)
		return nil, errors.New("Could not load ADR configuration.")
	}

	am := AdrManager{
		Config: config,
	}

	return &am, nil
}

func (am AdrManager) Init() error {

	if utils.FileExists(configFilenName) {
		return errors.New("ADRs seem to be initialized already, config file '.adr.json' exists!")
	}

	// 1) Create config file
	am.Config.Store(".adr.json")

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
		log.Printf("Error when writing short ADR template: %v\n", errShort)
	}
	errLong := ioutil.WriteFile(pathlongTemplate, []byte(val.Long), 0644)
	if errLong != nil {
		log.Printf("Error when writing long ADR template: %v\n", errLong)
	}

	return nil
}

func (am AdrManager) AddAdr(title string) (string, error) {
	templateContent := am.loadTemplateOrDefault(filepath.Join(am.Config.Path, am.Config.TemplateName))
	return am.AddAdrWithContent(title, templateContent)
}

func (am AdrManager) AddAdrFromTemplate(title string, templateFile string) (string, error) {
	templContent := am.loadTemplateOrDefault(templateFile)

	return am.AddAdrWithContent(title, templContent)
}

func (am AdrManager) loadTemplateOrDefault(templateFile string) string {
	rawTemplate, err := os.ReadFile(filepath.Join(am.Config.Path, templateFile))
	var templContent string
	if err != nil {
		log.Printf("Could not read requested template file %s: %v\n", templateFile, err)
		log.Println("Use standard template instead!")
		val, _ := templates.TemplatesLibrary["en"]
		templContent = val.Short
	} else {
		templContent = string(rawTemplate)
	}

	return templContent
}

func (am AdrManager) AddAdrWithContent(title string, content string) (string, error) {
	newDate := createDateString()
	index := am.getNewIndexString()
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
	log.Printf("Identified template variables: %v\n", vars)

	tmpl, err := template.New("adr").Parse(content)
	if err != nil {
		panic(err)
	}

	am.createAdrFile(am.Config.Path, fileName, tmpl, vars)
	return fileName, nil
	// fileData := template.New("test").Parse(content).Execute()

	// return "", nil
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

// function createDecisions (name: string, savePath: string | any | void): string {
// 	let language = Config.getLanguage()
// 	let raw = fs.readFileSync(getTemplatePath(language), 'utf8')

// 	let filePath = savePath + newIndex + '-' + fileName + '.md'
// 	fs.writeFileSync(filePath, fileData)

// 	return filePath
//   }

//   export function create (name: string) {
// 	let savePath = Config.getSavePath()
// 	let i18n = Utils.getI18n()
// 	console.log(i18n.logSavePath + savePath)
// 	mkdirp.sync(savePath)

// 	const filePath = createDecisions(name, savePath)
// 	Utils.openInEditor(path.join(process.cwd(), filePath))

//   }

func createDateString() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

func (am AdrManager) getNewIndexString() string {
	lastIndex, err := am.getLatestIndex()
	if err != nil {
		return am.createIndexByNumber(1)
	}
	lastIndex = lastIndex + 1
	return am.createIndexByNumber(lastIndex)
}

func (am AdrManager) getLatestIndex() (int, error) {
	files, err := am.getAdrFiles()

	if err != nil {
		log.Printf("Error when trying to load existing ADR files: %v\n", err)
		return 0, err
	}

	if len(files) == 0 {
		log.Println("Found no ADR files.")
		return 0, errors.New("No ADR files found.")
	}

	return am.getMaxIndex(files), nil
}

func (am AdrManager) getMaxIndex(filenames []string) int {
	maxNumber := 0

	for _, file := range filenames {
		log.Printf("Trying to extract index from file of name '%s'\n", file)
		indexPart := file[:am.Config.Digits]
		log.Printf("Found index parts: %s\n", indexPart)
		index, err := strconv.Atoi(indexPart)

		if err == nil && index > maxNumber {
			maxNumber = index
		}
	}

	return maxNumber
}

func (am AdrManager) getAdrFiles() ([]string, error) {
	files, err := ioutil.ReadDir(am.Config.Path)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for _, file := range files {
		fmt.Printf("Analyzing file %v, Name='%s', IsDir=%v with file extension='%s'\n", file, file.Name(), file.IsDir(), filepath.Ext(file.Name()))
		if !file.IsDir() && file.Name() != "README.md" && !strings.HasPrefix(file.Name(), "template-") && file.Name() != am.Config.TemplateName && filepath.Ext(file.Name()) == ".md" {
			res = append(res, file.Name())
		}
	}

	return res, nil
}

func (am AdrManager) createIndexByNumber(number int) string {
	s := fmt.Sprintf("%020d", number)
	log.Printf("Trying to create index by number: %s", s)
	return am.Config.Prefix + s[len(s)-am.Config.Digits:]
}

func generateBaseFileName(title string) string {

	return strings.Trim(strings.ToLower(title), " ")
}

// function generateFileName (originFileName) {
// 	return originFileName.toLowerCase().trim()
// 	  .replace(/[\s_-]+/g, '-') // swap any length of whitespace, underscore, hyphen characters with a single _
// 	  .replace(/^-+|-+$/g, '') // remove leading, trailing -
// 	  .replace(/，/g, '')
// 	  .replace(/。/g, '')
// 	  .replace(/ /g, '-')
// 	  .replace(/\?/g, '-')
// 	  .replace(/#/g, '')
// 	  .replace(/:/g, '')
// 	  .replace(/# /g, '')
//   }
