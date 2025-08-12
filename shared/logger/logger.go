package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var globalLogger = zerolog.New(os.Stdout).
	With().
	Timestamp().
	Logger()

// Get returns the global logger instance.
func Get() *zerolog.Logger {
	return &globalLogger
}
