package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/vishenosik/web/api"
	attrs "github.com/vishenosik/web/log"
)

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

func RequestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			timeStart := time.Now()
			lrw := newLoggingResponseWriter(w)

			log := logger.With(
				slog.String("method", fmt.Sprintf("%s %s", r.Method, r.URL.Path)),
			)

			defer func() {
				if api.IsClientError(lrw.statusCode) || api.IsServerError(lrw.statusCode) {
					log.Error("request failed with error",
						slog.Int("code", lrw.statusCode),
						//TODO slog.String("RequestID", attrs.RequestID(r)),
						attrs.Took(timeStart),
					)
				} else if api.IsRedirect(lrw.statusCode) {
					log.Warn("request redirected",
						slog.Int("code", lrw.statusCode),
						//TODO slog.String("RequestID", attrs.RequestID(r)),
						attrs.Took(timeStart),
					)
				} else {
					logger.Info("request accepted",
						slog.Int("code", lrw.statusCode),
						//TODO slog.String("RequestID", attrs.RequestID(r)),
						attrs.Took(timeStart),
					)
				}
			}()
			next.ServeHTTP(lrw, r)
		}
		return http.HandlerFunc(fn)
	}
}
