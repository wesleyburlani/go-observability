package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

func ParseLevel(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	default:
		return LevelInfo, fmt.Errorf("invalid log level: %s", level)
	}
}

type Options struct {
	Enabled bool
	Level   Level
}

type Logger struct {
	logger *slog.Logger
}

// returns a new logger.
func NewLogger(options Options) *Logger {
	w := io.Discard
	if options.Enabled {
		w = os.Stdout
	}

	return &Logger{
		logger: slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: options.Level})),
	}
}

// returns a new logger with the given attributes.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

func (l *Logger) Debug(ctx context.Context, msg string) {
	l.Log(ctx, LevelDebug, msg)
}

func (l *Logger) Info(ctx context.Context, msg string) {
	l.Log(ctx, LevelInfo, msg)
}

func (l *Logger) Warn(ctx context.Context, msg string) {
	l.Log(ctx, LevelWarn, msg)
}

func (l *Logger) Error(ctx context.Context, msg string) {
	l.Log(ctx, LevelError, msg)
}

func (l *Logger) Log(ctx context.Context, level Level, msg string) {
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()
	traceId := spanContext.TraceID()
	spanId := spanContext.SpanID()
	l.logger.With(
		"trace_id", traceId,
		"span_id", spanId,
	).Log(ctx, level, msg)
}
