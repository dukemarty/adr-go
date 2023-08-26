/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
package logic

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/dukemarty/adr-go/data"
)

func GetStatusInteractively(pretext string) string {
	newStatus := ""
	prompt := &survey.Select{
		Message: pretext + " new status:",
		Options: data.SupportedStatus,
	}
	survey.AskOne(prompt, &newStatus)

	return newStatus
}
