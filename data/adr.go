package data

import (
	"log"
	"os"
	"strings"

	"github.com/dukemarty/adr-go/utils"
	"github.com/olekukonko/tablewriter"
)

type AdrVars struct {
	NUMBER string
	TITLE  string
	DATE   string
}

func ReadStatusSection(adrFile string) {
	doc, err := utils.OpenMarkdownFile(adrFile)
	// data, err := os.ReadFile(adrFile)
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

	tbl := tablewriter.NewWriter(os.Stdout)
	tbl.SetHeader([]string{"Date of Change", "Status"})
	nextChild := section.NextSibling()
	for nextChild != nil && !(nextChild.Kind().String() == "Heading") {
		line := string(nextChild.Text(doc.Source))
		tokens := strings.Split(line, " ")
		tbl.Append([]string{tokens[0], tokens[1]})
		nextChild = nextChild.NextSibling()
	}

	tbl.Render()

	if err != nil {
		log.Fatalf("Error converting data: %v", err)
	}
	// md.Parser().Parse(data)

}
