package logger

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func New(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("middleware logger enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := log.With(
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
			)

			ctx := r.Context()
			ctx = context.WithValue(ctx, "log", log)
			r = r.WithContext(ctx)

			t1 := time.Now()

			defer func() {
				log.Info("request completed", zap.Duration("duration", time.Since(t1)))
			}()

			next.ServeHTTP(w, r)

		}
		return http.HandlerFunc(fn)
	}
}
