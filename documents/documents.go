package documents

import (
	_ "embed"
)

//go:embed CHANGELOG.md
var Changelog string

//go:embed LICENSE
var License string

type DocumentFile struct {
	Filename string
	Content  string
}

var Docs = []DocumentFile{
	{Filename: "CHANGELOG.md", Content: Changelog},
	{Filename: "LICENSE", Content: License},
}
