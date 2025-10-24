package middleware

import (
	"context"
	"net/http"
	"time"
)

type TimeoutConfig struct {
	DefaultTimeout time.Duration
}

func NewTimeout(config TimeoutConfig) func(http.Handler) http.Handler {
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = 30 * time.Second
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), config.DefaultTimeout)
			defer cancel()

			rw := &responseWriter{ResponseWriter: w}

			done := make(chan struct{})

			go func() {
				next.ServeHTTP(rw, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				if !rw.written {
					http.Error(w, "Request Timeout", http.StatusRequestTimeout)
				}
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	written bool
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.written = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.written = true
	return rw.ResponseWriter.Write(b)
}
