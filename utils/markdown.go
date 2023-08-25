/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
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

func (doc *MarkdownDoc) FindMarkdownSection(sectionName string) mdast.Node {
	node := doc.Doc.FirstChild()
	for node != nil && !(node.Kind().String() == "Heading" && strings.ToUpper(string(node.Text(doc.Source))) == strings.ToUpper(sectionName)) {
		node = node.NextSibling()
	}

	return node
}

func (doc *MarkdownDoc) FindInsertAtEndOfSection(sectionName string) (preText string, postText string) {

	var cutPoint int
	node := doc.Doc.FirstChild()
	for node != nil && !(node.Kind().String() == "Heading" && strings.ToUpper(string(node.Text(doc.Source))) == strings.ToUpper(sectionName)) {
		node = node.NextSibling()
	}
	for node != nil && node.Kind().String() == "Heading" {
		node = node.NextSibling()
	}
	for node != nil && node.Kind().String() != "Heading" {
		node = node.NextSibling()
	}

	node = node.PreviousSibling()
	cutPoint = node.Lines().At(node.Lines().Len()-1).Stop + 1

	return string(doc.Source[:cutPoint]), string(doc.Source[cutPoint:])
}
