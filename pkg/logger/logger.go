package logger

import (
	"context"
	"os"

	"github.com/sagikazarmark/slog-shim"
	"go.opentelemetry.io/otel/trace"
)

type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type Logger struct {
	logger *slog.Logger
	ctx    *context.Context
}

// returns a new logger.
func NewLogger() *Logger {
	return &Logger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		ctx:    nil,
	}
}

// returns a new logger with the given context.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		logger: l.logger,
		ctx:    &ctx,
	}
}

// returns a new logger with the given attributes.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
		ctx:    l.ctx,
	}
}

func (l *Logger) Debug(msg string) {
	l.Log(LevelDebug, msg)
}

func (l *Logger) Info(msg string) {
	l.Log(LevelInfo, msg)
}

func (l *Logger) Warn(msg string) {
	l.Log(LevelWarn, msg)
}

func (l *Logger) Error(msg string) {
	l.Log(LevelError, msg)
}

func (l *Logger) Log(level Level, msg string) {
	if l.ctx == nil {
		l.logger.Log(context.Background(), level, msg)
		return
	}
	span := trace.SpanFromContext(*l.ctx)
	spanContext := span.SpanContext()
	traceId := spanContext.TraceID()
	spanId := spanContext.SpanID()
	l.logger.With(
		"trace_id", traceId,
		"span_id", spanId,
	).Log(*l.ctx, level, msg)
}
