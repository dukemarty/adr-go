/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
package utils

import (
	"errors"
	"os"
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
