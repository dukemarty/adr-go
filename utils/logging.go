package utils

import (
	"io"
	"log"
)

func SetupLogger(verbose bool) *log.Logger {
	logger := log.Default()
	if !verbose {
		logger.SetFlags(0)
		logger.SetOutput(io.Discard)
	}

	return logger
}
