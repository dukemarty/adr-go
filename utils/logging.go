/*
Copyright © 2023 Martin Loesch <development@martinloesch.net>
*/
package utils

import (
	"io"
	"log"
)

// Create and returns a new logger object, depending on flag
// verbose either with normal output or discarding all messages.
func SetupLogger(verbose bool) *log.Logger {
	logger := log.Default()
	if !verbose {
		logger.SetFlags(0)
		logger.SetOutput(io.Discard)
	}

	return logger
}
