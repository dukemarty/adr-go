package templates

import (
	_ "embed"
)

type TemplatesSet struct {
	Short string
	Long  string
}

//go:embed en-template-short.md
var shortStandardTemplateEn string

//go:embed en-template-long.md
var longStandardTemplateEn string

var TemplatesLibrary = make(map[string]TemplatesSet)

func init() {
	TemplatesLibrary["en"] = TemplatesSet{Short: shortStandardTemplateEn, Long: longStandardTemplateEn}
}
