// The data package	contains data structures and functions for accessing and handling
// data objects, in particular the handled files: Config file for a project, and single
// ADR files.
package data

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dukemarty/adr-go/utils"
)

type AdrVars struct {
	NUMBER string
	TITLE  string
	DATE   string
}

type AdrInfo struct {
	RelativePath string
	Index        int
	Title        string
}

func LoadAdrInfo(logger *log.Logger, basepath string, adrFile string) (AdrInfo, error) {
	var res AdrInfo
	res.RelativePath = filepath.Join(basepath, adrFile)

	index, title, err := extractAdrBaseInfoFromFile(logger, res.RelativePath)
	if err != nil {
		return res, err
	}
	res.Index = index
	res.Title = title

	return res, nil
}

func extractAdrBaseInfoFromFile(logger *log.Logger, adrFile string) (int, string, error) {
	doc, err := utils.OpenMarkdownFile(adrFile)

	node := doc.Doc.FirstChild()
	for node != nil && !(node.Kind().String() == "Heading") {
		node = node.NextSibling()
	}
	logger.Printf("Extracted heading line: %s\n", string(node.Text(doc.Source)))

	// TODO: only very basic, error-prone parsing used here, improve!
	tokens := strings.Split(string(node.Text(doc.Source)), " ")

	indexPart := tokens[0][:len(tokens[0])-1]
	index, err := strconv.Atoi(indexPart)
	if err != nil {
		logger.Printf("Could not parse '%s' as index: %v\n", indexPart, err)
		return -1, "", err
	}
	title := strings.Join(tokens[1:], " ")

	return index, title, nil
}

type StatusChange struct {
	Date   string
	Status string
}

// ReadStatusEntries can be used to single out and read the status section of an ADR.
func ReadStatusEntries(logger *log.Logger, adrFile string) ([]StatusChange, error) {
	doc, err := utils.OpenMarkdownFile(adrFile)
	if err != nil {
		logger.Fatalf("Could not read data from ADR '%s': %v", adrFile, err)
	}

	section := doc.FindMarkdownSection("Status")

	res := make([]StatusChange, 0)
	nextChild := section.NextSibling()
	for nextChild != nil && !(nextChild.Kind().String() == "Heading") {
		line := string(nextChild.Text(doc.Source))
		tokens := strings.Split(line, " ")
		res = append(res, StatusChange{Date: tokens[0], Status: tokens[1]})
		nextChild = nextChild.NextSibling()
	}

	return res, nil
}

func AddStatusEntry(logger *log.Logger, adrFile string, newStatus string) {
	doc, err := utils.OpenMarkdownFile(adrFile)
	if err != nil {
		logger.Fatalf("Could not read data from ADR '%s': %v", adrFile, err)
	}

	prePart, postPart := doc.FindInsertAtEndOfSection("Status")

	err = os.Rename(adrFile, path.Join(adrFile+".bak"))
	if err != nil {
		logger.Printf("Could not rename ADR file '%s': %v\n", adrFile, err)
	}

	err = os.WriteFile(adrFile, []byte(fmt.Sprintf("%s\n%v %s\n%s", prePart, time.Now().Format("2006-01-02"), newStatus, postPart)), 0644)
	if err != nil {
		logger.Printf("Could not write changed ADR file '%s': %v\n", adrFile, err)
	}
}
