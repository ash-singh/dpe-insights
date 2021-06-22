package logger

import (
	"io"
	"log"
	"os"
)

// Configure logger.
func Configure() {
	mw := io.Writer(os.Stdout)
	log.SetOutput(mw)
}
