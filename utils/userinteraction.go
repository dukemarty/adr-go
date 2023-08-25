/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
package utils

import (
	"github.com/AlecAivazis/survey/v2"
)

var SupportedStatus = []string{"Proposed", "Accepted", "Done", "Deprecated", "Superseded"}

func GetStatusInteractively(pretext string) string {
	newStatus := ""
	prompt := &survey.Select{
		Message: pretext + " new status:",
		Options: SupportedStatus,
	}
	survey.AskOne(prompt, &newStatus)

	return newStatus
}
