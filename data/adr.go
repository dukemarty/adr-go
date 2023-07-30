// The data package	contains data structures and functions for accessing and handling
// data objects, in particular the handled files: Config file for a project, and single
// ADR files.
package data

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"

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

func LoadAdrInfo(basepath string, adrFile string) (AdrInfo, error) {
	var res AdrInfo
	res.RelativePath = filepath.Join(basepath, adrFile)

	index, title, err := extractAdrBaseInfoFromFile(res.RelativePath)
	if err != nil {
		return res, err
	}
	res.Index = index
	res.Title = title

	return res, nil
}

func extractAdrBaseInfoFromFile(adrFile string) (int, string, error) {
	doc, err := utils.OpenMarkdownFile(adrFile)

	node := doc.Doc.FirstChild()
	for node != nil && !(node.Kind().String() == "Heading") {
		// fmt.Printf("** Skip NODE %d:\n", i)
		node = node.NextSibling()
	}

	log.Printf("Extracted heading line: %s\n", string(node.Text(doc.Source)))
	// TODO: only very basic, error-prone parsing used here, improve!
	tokens := strings.Split(string(node.Text(doc.Source)), " ")

	indexPart := tokens[0][:len(tokens[0])-1]
	index, err := strconv.Atoi(indexPart)
	if err != nil {
		log.Printf("Could not parse '%s' as index: %v\n", indexPart, err)
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
func ReadStatusEntries(adrFile string) ([]StatusChange, error) {
	doc, err := utils.OpenMarkdownFile(adrFile)
	if err != nil {
		log.Fatalf("Could not read data from ADR '%s': %v", adrFile, err)
	}
	// md := goldmark.New(goldmark.WithParserOptions(parser.WithBlockParsers()))
	// doc := md.Parser().Parse(text.NewReader(data))

	section := utils.FindMarkdownSection(doc, "Status")
	// node := doc.FirstChild()
	// i := 0
	// for node != nil && !(node.Kind().String() == "Heading" && strings.ToUpper(string(node.Text(data))) == "STATUS") {
	// 	fmt.Printf("** Skip NODE %d:\n", i)
	// 	node = node.NextSibling()
	// 	i = i + 1
	// }
	// fmt.Println("Found Status block:")
	// node.Dump(data, 4)
	// // fmt.Printf("%v\n", doc.OwnerDocument().)

	res := make([]StatusChange, 0)
	// tbl := tablewriter.NewWriter(os.Stdout)
	// tbl.SetHeader([]string{"Date of Change", "Status"})
	nextChild := section.NextSibling()
	for nextChild != nil && !(nextChild.Kind().String() == "Heading") {
		// blockLines := nextChild.Lines()
		// log.Printf("Lines: %v\n", blockLines)
		// section := string(nextChild.Text(doc.Source))
		// log.Printf("Section to parse: <%s>\n", section)
		// log.Printf("Section to parse: <%v>\n", []byte(section))
		// log.Printf("Section to parse: <%v>\n", nextChild.Text(doc.Source))
		// lines := strings.Split(section, "\r")
		// for _, line := range lines {
		line := string(nextChild.Text(doc.Source))
		tokens := strings.Split(line, " ")
		// tbl.Append([]string{tokens[0], tokens[1]})
		res = append(res, StatusChange{Date: tokens[0], Status: tokens[1]})
		nextChild = nextChild.NextSibling()
		// }
	}

	// tbl.Render()

	if err != nil {
		log.Fatalf("Error converting data: %v", err)
		return nil, err
	}
	// md.Parser().Parse(data)

	// return nil, errors.New("Not implemented yet!")
	return res, nil
}
