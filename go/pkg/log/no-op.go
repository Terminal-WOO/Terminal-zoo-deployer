package log

import (
	"context"
)

// NoOpLogger is a logger that does nothing (used when logging is disabled)
type NoOpLogger struct{}

func NewNoOpLogger() *NoOpLogger { return new(NoOpLogger) }

// Implement all log methods with empty bodies
func (n *NoOpLogger) Debug(msg string, args ...any)                         {}
func (n *NoOpLogger) DebugCtx(ctx context.Context, msg string, args ...any) {}
func (n *NoOpLogger) Info(msg string, args ...any)                          {}
func (n *NoOpLogger) InfoCtx(ctx context.Context, msg string, args ...any)  {}
func (n *NoOpLogger) Warn(msg string, args ...any)                          {}
func (n *NoOpLogger) WarnCtx(ctx context.Context, msg string, args ...any)  {}
func (n *NoOpLogger) Error(msg string, args ...any)                         {}
func (n *NoOpLogger) ErrorCtx(ctx context.Context, msg string, args ...any) {}
func (n *NoOpLogger) With(args ...any) Logger                               { return n }
