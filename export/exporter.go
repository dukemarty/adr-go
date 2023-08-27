package adrexport

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dukemarty/adr-go/logic"
)

// Interface for exporters, transforming a list of ADRs (status infos) into a string.
type AdrListExporter interface {
	Export(*log.Logger, []logic.AdrStatus) string
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
		logger.Println("Exporter type 'html' not implemented yet!")
		resErr = errors.New("Exporter type 'html' not implemented yet!")
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

func (CsvExporter) Export(logger *log.Logger, entries []logic.AdrStatus) string {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	w.Write([]string{"Index", "Decision", "Last Modified Date", "Last Status"})
	if err := w.Error(); err != nil {
		logger.Fatal("Error writing csv: %v\n", err)
	}

	for _, e := range entries {
		w.Write([]string{fmt.Sprintf("%04d", e.Index), e.Title, e.LastModified, e.LastStatus})
		if err := w.Error(); err != nil {
			logger.Fatal("Error writing csv: %v\n", err)
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

func (JsonExporter) Export(logger *log.Logger, entries []logic.AdrStatus) string {
	data := make([]JsonAdrData, 0)
	for _, e := range entries {
		nextEntry := JsonAdrData{Index: e.Index, Decision: e.Title, LastModified: e.LastModified, LastStatus: e.LastStatus}
		data = append(data, nextEntry)
	}

	jsonData, _ := json.Marshal(data)

	return string(jsonData)
}
