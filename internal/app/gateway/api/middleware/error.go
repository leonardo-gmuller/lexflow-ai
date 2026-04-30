package middleware

import (
	"net/http"
	"time"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/library/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func HTTPErrorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		if rec.status >= 400 {
			logger.Error.Printf(
				"[HTTP_ERROR] method=%s path=%s status=%d ip=%s duration=%s user_agent=%s error=%s\n",
				r.Method,
				r.URL.Path,
				rec.status,
				r.RemoteAddr,
				time.Since(start),
				r.UserAgent(),
				http.StatusText(rec.status),
			)
		}
	})
}
