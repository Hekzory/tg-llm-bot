package logging

import (
	"errors"
	"log/slog"
	"os"
)

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	// DEBUG is for detailed debugging information.
	DEBUG LogLevel = iota
	// INFO is for general informational messages.
	INFO
	// WARNING is for potential issues that don't necessarily indicate an error.
	WARNING
	// ERROR is for errors that should be investigated.
	ERROR
	// FATAL is for critical errors that cause the application to terminate.
	FATAL
)

// Logger provides logging functionality with configurable log levels.
type Logger struct {
	logger *slog.Logger
}

// NewLogger creates a new Logger with the specified log level.
func NewLogger(logLevel string) (*Logger, error) {
	var handler slog.Handler
	switch logLevel {
	case "DEBUG":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "INFO":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	case "WARNING":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})
	case "ERROR":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})
	case "FATAL":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}) // FATAL is treated as ERROR in slog
	default:
		return nil, errors.New("Could not create logger, log level " + logLevel + " is not supported")
	}

	logger := slog.New(handler)
	return &Logger{logger: logger}, nil
}

// Debug logs a message at the DEBUG level.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.logger.Debug(format, args...)
}

// Info logs a message at the INFO level.
func (l *Logger) Info(format string, args ...interface{}) {
	l.logger.Info(format, args...)
}

// Warning logs a message at the WARNING level.
func (l *Logger) Warning(format string, args ...interface{}) {
	l.logger.Warn(format, args...)
}

// Error logs a message at the ERROR level.
func (l *Logger) Error(format string, args ...interface{}) {
	l.logger.Error(format, args...)
}

// Fatal logs a message at the FATAL level and exits the application.
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logger.Error(format, args...) // FATAL is treated as ERROR in slog
	os.Exit(1)
}
