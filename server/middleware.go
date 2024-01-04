package server

import (
	"github.com/umtdemr/go-todo/logger"
	"net/http"
	"time"
)

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := logger.Get()

		defer func() {
			l.
				Info().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Dur("elapsed_ms", time.Since(start)).
				Msg("request")
		}()

		next.ServeHTTP(w, r)
	})
}
