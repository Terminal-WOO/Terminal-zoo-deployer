package log

import (
	"context"
	"log/slog"
	"os"
)

const DefaultPackageKey string = "component"

var ComponentKey string = DefaultPackageKey

// DefaultLogger wraps slog.Logger to implement our Logger interface
type DefaultLogger struct {
	logger *slog.Logger
}

// NewDefaultLogger creates a JSON-based slog logger
func NewDefaultLogger() *DefaultLogger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return &DefaultLogger{
		logger: slog.New(handler),
	}
}

// Implement all logging methods
func (l *DefaultLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *DefaultLogger) DebugCtx(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *DefaultLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *DefaultLogger) InfoCtx(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *DefaultLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *DefaultLogger) WarnCtx(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *DefaultLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *DefaultLogger) ErrorCtx(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *DefaultLogger) With(args ...any) Logger {
	return &DefaultLogger{
		logger: l.logger.With(args...),
	}
}
