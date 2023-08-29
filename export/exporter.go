package adrexport

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	_ "embed"

	"github.com/dukemarty/adr-go/logic"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/toc"
)

// Interface for exporters, transforming a list of ADRs (status infos) into a string.
type AdrListExporter interface {
	Export(logger *log.Logger, data []logic.AdrStatus, dataPath string) string
}

// List of supported exporter types.
var SupportedExporters = []string{"csv", "json", "markdown", "md", "html"}

func CreateExporter(logger *log.Logger, expType string) (AdrListExporter, error) {
	var res AdrListExporter = nil
	var resErr error = nil
	switch strings.ToLower(expType) {
	case "csv":
		res = CsvExporter{}
	case "json":
		res = JsonExporter{}
	case "md":
	case "markdown":
		logger.Println("Exporter type 'markdown' not implemented yet!")
		resErr = errors.New("Exporter type 'markdown' not implemented yet!")
	case "html":
		res = HtmlExporter{}
	default:
		logger.Printf("Exporter type '%s' not supported!\n", expType)
		resErr = errors.New(fmt.Sprintf("Exporter type '%s' not supported!", expType))
	}

	return res, resErr
}

// ----------------------------------------------------------------------------
// Implementation of an AdrListExporter for CSV data

// Empty struct to represent an exporter of csv data.
type CsvExporter struct{}

func (CsvExporter) Export(logger *log.Logger, entries []logic.AdrStatus, _ string) string {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	w.Write([]string{"Index", "Decision", "Last Modified Date", "Last Status"})
	if err := w.Error(); err != nil {
		logger.Printf("Error writing csv: %v\n", err)
		return ""
	}

	for _, e := range entries {
		w.Write([]string{fmt.Sprintf("%04d", e.Index), e.Title, e.LastModified, e.LastStatus})
		if err := w.Error(); err != nil {
			logger.Printf("Error writing csv: %v\n", err)
			return ""
		}
	}

	// have to Flush() explicitly, because only .Write() is used...
	w.Flush()

	return buf.String()
}

// ----------------------------------------------------------------------------
// Implementation of an AdrListExporter for JSON data

// Empty struct to represent an exporter of json data.
type JsonAdrData struct {
	Index        int    `json:"index"`
	Decision     string `json:"decision"`
	LastModified string `json:"modifiedDate"`
	LastStatus   string `json:"lastStatus"`
}

type JsonExporter struct{}

func (JsonExporter) Export(logger *log.Logger, entries []logic.AdrStatus, _ string) string {
	data := make([]JsonAdrData, 0)
	for _, e := range entries {
		nextEntry := JsonAdrData{Index: e.Index, Decision: e.Title, LastModified: e.LastModified, LastStatus: e.LastStatus}
		data = append(data, nextEntry)
	}

	jsonData, _ := json.Marshal(data)

	return string(jsonData)
}

// ----------------------------------------------------------------------------
// Implementation of an AdrListExporter for JSON data

//go:embed export_template.html
var HtmlTemplate string

type HtmlExportVars struct {
	HTMLTOC     template.HTML
	HTMLCONTENT template.HTML
}

type HtmlExporter struct{}

// ByIndex implements sort.Interface based on the Index field for AdrStatus slices.
type ByIndex []logic.AdrStatus

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (HtmlExporter) Export(logger *log.Logger, entries []logic.AdrStatus, dataPath string) string {

	// resort entries based on their index
	sort.Sort(ByIndex(entries))

	// assemble all ADRs into single in-memory document
	var sb strings.Builder

	for _, e := range entries {
		buf, err := os.ReadFile(path.Join(dataPath, e.Filename))
		if err != nil {
			logger.Printf("Could not read ADR from file '%s': %v\n", e.Filename, err)
		}
		sb.Write(buf)
		sb.WriteString("\n\n")
	}
	source := []byte(sb.String())

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	// create content html
	var contentBuf bytes.Buffer
	if err := md.Convert(source, &contentBuf); err != nil {
		logger.Printf("Could not render ADRs content as html: %v\n", err)
		return ""
	}

	// create toc html
	doc := md.Parser().Parse(text.NewReader(source))
	tree, err := toc.Inspect(doc, source)
	if err != nil {
		logger.Printf("Could not render ADRs toc as html: %v\n", err)
		return ""
	}
	list := toc.RenderList(tree)

	var tocBuf bytes.Buffer
	if list != nil {
		// list will be nil if the table of contents is empty
		// because there were no headings in the document.
		md.Renderer().Render(&tocBuf, source, list)
	}

	// assemble content and toc html parts in export html template
	vars := HtmlExportVars{
		HTMLCONTENT: template.HTML(contentBuf.String()),
		HTMLTOC:     template.HTML(tocBuf.String()),
	}

	tmpl, err := template.New("htmlexport").Parse(HtmlTemplate)
	if err != nil {
		panic(err)
	}

	var completeExportBuf bytes.Buffer
	err = tmpl.Execute(&completeExportBuf, vars)
	if err != nil {
		log.Fatal(err)
	}

	return completeExportBuf.String()
}
