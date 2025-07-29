package logging

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger creates and configures a new zerolog.Logger.
func NewLogger() zerolog.Logger {
	// Configure lumberjack for log rotation.
	logRotator := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // compress old log files
	}

	// Create a multi-writer to write to both stdout and the log file.
	multi := io.MultiWriter(os.Stdout, logRotator)

	// Create a ConsoleWriter for human-readable output.
	consoleWriter := zerolog.ConsoleWriter{
		Out:        multi,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    true, // No colors for cleaner file logs
	}

	// Create the logger.
	return zerolog.New(consoleWriter).With().Timestamp().Logger()
}
