package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		l.logger.Debug("request received",
			"method", r.Method,
			"path", r.URL.String(),
			"origin", r.RemoteAddr,
		)
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		l.logger.Debug("request completed",
			"method", r.Method,
			"path", r.URL.String(),
			"origin", r.RemoteAddr,
			"latency", time.Since(t1).String(),
			"status", lrw.statusCode,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
