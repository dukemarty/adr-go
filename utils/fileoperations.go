/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
package utils

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		// filename exists
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		// filename does *not* exist
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false
	}
}

func EditFile(filename string, editor string, defaultEditor string, logger *log.Logger) {
	if len(editor) == 0 {
		if len(defaultEditor) == 0 {
			editor = os.Getenv("EDITOR")
		} else {
			editor = defaultEditor
		}
	}
	if len(editor) > 0 {
		logger.Printf("Opening '%s' with '%s'\n", filename, editor)
		cmd := exec.Command(editor, filename)
		err := cmd.Start()
		if err == nil {
			cmd.Process.Release()
		} else {
			logger.Printf("Error when trying to start editor: %v\n", err)
		}
	} else {
		logger.Println("EDITOR environment variable not set, therefor ADR can not be opened.")
	}

}
