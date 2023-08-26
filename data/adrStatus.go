package data

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var SupportedStatus = []string{"Proposed", "Accepted", "Done", "Deprecated", "Superseded"}
var statusForComparisons = []string{"PROPOSED", "ACCEPTED", "DONE", "DEPRECATED", "SUPERSEDED"}

type AdrStatus string

// String is used both by fmt.Print and by Cobra in help text
func (e *AdrStatus) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *AdrStatus) Set(v string) error {
	if slices.Contains(statusForComparisons, strings.ToUpper(v)) {
		caser := cases.Title(language.English)
		*e = AdrStatus(caser.String(v))
		return nil
	} else {
		return errors.New(fmt.Sprintf(`must be one of: %v`, SupportedStatus))
	}
}

// Type is only used in help text
func (e *AdrStatus) Type() string {
	return "AdrStatus"
}
