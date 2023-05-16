package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	mdast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type MarkdownDoc struct {
	Source []byte
	Doc    mdast.Node
}

func OpenMarkdownFile(filename string) (*MarkdownDoc, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read data from Markdown '%s': %v", filename, err))
	}
	md := goldmark.New(goldmark.WithParserOptions(parser.WithBlockParsers()))
	doc := md.Parser().Parse(text.NewReader(data))
	res := MarkdownDoc{Source: data, Doc: doc}

	return &res, nil
}

func FindMarkdownSection(doc *MarkdownDoc, sectionName string) mdast.Node {
	node := doc.Doc.FirstChild()
	i := 0
	for node != nil && !(node.Kind().String() == "Heading" && strings.ToUpper(string(node.Text(doc.Source))) == strings.ToUpper(sectionName)) {
		// fmt.Printf("** Skip NODE %d:\n", i)
		node = node.NextSibling()
		i = i + 1
	}

	return node
}
