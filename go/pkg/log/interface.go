package log

import "context"

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger interface with all log levels
type Logger interface {
	Debug(msg string, args ...any)
	DebugCtx(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoCtx(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnCtx(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorCtx(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}
