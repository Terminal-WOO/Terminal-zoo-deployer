package main

import (
	"context"

	"github.com/ClappFormOrg/AI-CO/go/pkg/log"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Entry
	level  log.Level
}

// NewLogrusLogger creates a new Logrus-based Logger with a specified log level
func NewLogrusLogger(level log.Level) *LogrusLogger {
	logger := logrus.New()
	// Set log level based on the input Level enum
	switch level {
	case log.LevelDebug:
		logger.SetLevel(logrus.DebugLevel)
	case log.LevelInfo:
		logger.SetLevel(logrus.InfoLevel)
	case log.LevelWarn:
		logger.SetLevel(logrus.WarnLevel)
	case log.LevelError:
		logger.SetLevel(logrus.ErrorLevel)
	}
	return &LogrusLogger{
		logger: logrus.NewEntry(logger),
		level:  level,
	}
}

func (l *LogrusLogger) fieldsFromArgs(args ...any) logrus.Fields {
	// Preallocate map based on the length of args, assuming it's an even number of elements (key, value pairs)
	fields := make(logrus.Fields, len(args)/2)

	// Iterate through the args slice
	for i := 0; i < len(args); i += 2 {
		// Check if args[i] is a string (key) and args[i+1] is a valid value
		if key, ok := args[i].(string); ok && i+1 < len(args) {
			// Only add valid key-value pairs to the fields map
			fields[key] = args[i+1]
		}
	}

	return fields
}

// Debug logs a debug message
func (l *LogrusLogger) Debug(msg string, args ...any) {
	if l.level <= log.LevelDebug {
		l.logger.WithFields(l.fieldsFromArgs(args...)).Debug(msg)
	}
}

// DebugCtx logs a debug message with a context
func (l *LogrusLogger) DebugCtx(ctx context.Context, msg string, args ...any) {
	if l.level <= log.LevelDebug {
		// Here, you could include context-specific info (e.g., request IDs)
		// For example, attaching context values to the log can be done here if needed
		l.logger.WithContext(ctx).WithFields(l.fieldsFromArgs(args...)).Debug(msg)
	}
}

// Info logs an info message
func (l *LogrusLogger) Info(msg string, args ...any) {
	if l.level <= log.LevelInfo {
		l.logger.WithFields(l.fieldsFromArgs(args...)).Info(msg)
	}
}

// InfoCtx logs an info message with a context
func (l *LogrusLogger) InfoCtx(ctx context.Context, msg string, args ...any) {
	if l.level <= log.LevelInfo {
		l.logger.WithContext(ctx).WithFields(l.fieldsFromArgs(args...)).Info(msg)
	}
}

// Warn logs a warning message
func (l *LogrusLogger) Warn(msg string, args ...any) {
	if l.level <= log.LevelWarn {
		l.logger.WithFields(l.fieldsFromArgs(args...)).Warn(msg)
	}
}

// WarnCtx logs a warning message with a context
func (l *LogrusLogger) WarnCtx(ctx context.Context, msg string, args ...any) {
	if l.level <= log.LevelWarn {
		l.logger.WithContext(ctx).WithFields(l.fieldsFromArgs(args...)).Warn(msg)
	}
}

// Error logs an error message
func (l *LogrusLogger) Error(msg string, args ...any) {
	if l.level <= log.LevelError {
		l.logger.WithFields(l.fieldsFromArgs(args...)).Error(msg)
	}
}

// ErrorCtx logs an error message with a context
func (l *LogrusLogger) ErrorCtx(ctx context.Context, msg string, args ...any) {
	if l.level <= log.LevelError {
		l.logger.WithContext(ctx).WithFields(l.fieldsFromArgs(args...)).Error(msg)
	}
}

// With adds additional fields to the logger
func (l *LogrusLogger) With(args ...any) log.Logger {
	return &LogrusLogger{
		logger: l.logger.WithFields(l.fieldsFromArgs(args...)),
		level:  l.level,
	}
}
