package httpapi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
)

type requestIDKey struct{}

var requestSequence atomic.Uint64

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

func (a *API) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Учебный номер уникален только в пределах этого процесса.
		id := fmt.Sprintf("req-%d", requestSequence.Add(1))
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), requestIDKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusWriter{ResponseWriter: w}
		completed := false
		defer func() {
			status := rec.status
			if status == 0 && completed {
				status = http.StatusOK
			}
			level := slog.LevelInfo
			if status >= 500 || !completed {
				level = slog.LevelError
			}
			a.logger.Log(r.Context(), level, "request",
				"request_id", RequestID(r.Context()), "method", r.Method,
				"path", r.URL.Path, "status", status,
				"duration_ms", float64(time.Since(start).Microseconds())/1000,
				"aborted", !completed)
		}()
		next.ServeHTTP(rec, r)
		completed = true
	})
}

func (a *API) recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if value := recover(); value != nil {
				if value == http.ErrAbortHandler {
					panic(value)
				}
				a.logger.ErrorContext(r.Context(), "panic recovered", "panic", value, "request_id", RequestID(r.Context()))
				// Уже отправленный ответ нельзя заменить на 500.
				if rec, ok := w.(*statusWriter); ok && rec.status != 0 {
					panic(http.ErrAbortHandler)
				}
				a.writeError(w, r, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Учебная обёртка для обычных ответов. Streaming и WebSocket требуют
// поддержки дополнительных интерфейсов ResponseWriter.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	if code >= 100 && code < 200 {
		w.ResponseWriter.WriteHeader(code)
		return
	}
	if w.status != 0 {
		return
	}
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

func (w *statusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
