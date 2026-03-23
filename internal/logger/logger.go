// Package logger provides functions to internal logger.
package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type logLevel int

const (
	// DEBUG level 1.
	DEBUG logLevel = iota + 1
	// INFO level 2.
	INFO
	// WARN level 3.
	WARN
	// ERROR level 4.
	ERROR
)

// Logger is main logger.
type Logger struct {
	Level  logLevel
	output io.Writer
}

// NewLogger returns new instance of the logger.
func NewLogger(level string) *Logger {
	switch level {
	case "DEBUG":
		return &Logger{Level: DEBUG, output: os.Stdout}
	case "INFO":
		return &Logger{Level: INFO, output: os.Stdout}
	case "WARN":
		return &Logger{Level: WARN, output: os.Stdout}
	case "ERROR":
		return &Logger{Level: ERROR, output: os.Stdout}
	default:
		return &Logger{Level: INFO, output: os.Stdout}
	}
}

func (l *Logger) log(level logLevel, levelName, msg string, args ...any) {
	if level < l.Level {
		return
	}
	timestamp := time.Now().UTC().Format(time.RFC3339)

	formattedMsg := fmt.Sprintf(msg, args...)

	fmt.Fprintf(l.output, "%s [%s]: %s\n", timestamp, levelName, formattedMsg)
}

// Debug message.
func (l *Logger) Debug(msg string, args ...any) {
	l.log(DEBUG, "DEBUG", msg, args...)
}

// Info message.
func (l *Logger) Info(msg string, args ...any) {
	l.log(INFO, "INFO", msg, args...)
}

// Warn message.
func (l *Logger) Warn(msg string, args ...any) {
	l.log(WARN, "WARN", msg, args...)
}

// Error message.
func (l *Logger) Error(msg string, args ...any) {
	l.log(ERROR, "ERROR", msg, args...)
}
