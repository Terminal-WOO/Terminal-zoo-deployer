package log

import (
	"context"
)

type ComponentLogger struct {
	logger Logger
}

func NewComponentLogger(logger Logger, component string) *ComponentLogger {
	return &ComponentLogger{
		logger: logger.With(ComponentKey, component),
	}
}

func (l *ComponentLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *ComponentLogger) DebugCtx(ctx context.Context, msg string, args ...any) {
	l.logger.DebugCtx(ctx, msg, args...)
}

func (l *ComponentLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *ComponentLogger) InfoCtx(ctx context.Context, msg string, args ...any) {
	l.logger.InfoCtx(ctx, msg, args...)
}

func (l *ComponentLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *ComponentLogger) WarnCtx(ctx context.Context, msg string, args ...any) {
	l.logger.WarnCtx(ctx, msg, args...)
}

func (l *ComponentLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *ComponentLogger) ErrorCtx(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorCtx(ctx, msg, args...)
}

func (l *ComponentLogger) With(args ...any) Logger {
	return l.logger.With(args...)
}
