package utils

import (
	"github.com/AlecAivazis/survey/v2"
)

func GetStatusInteractively(pretext string) string {
	newStatus := ""
	prompt := &survey.Select{
		Message: pretext + " new status:",
		Options: []string{"Proposed", "Accepted", "Done", "Deprecated", "Superseded"},
	}
	survey.AskOne(prompt, &newStatus)

	return newStatus
}
