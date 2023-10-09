package middlewares

import (
	"net/http"
	"time"

	"github.com/wesleyburlani/go-rest/pkg/logger"
)

type Logger struct {
	logger *logger.Logger
}

func NewLogger(logger *logger.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		logger := l.logger.With(
			"method", r.Method,
			"path", r.URL.String(),
			"origin", r.RemoteAddr,
		)
		logger.WithContext(r.Context()).Debug("request received")
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		logger.With(
			"latency", time.Since(t1).String(),
			"status", lrw.statusCode,
		).WithContext(r.Context()).Debug("request completed")
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
