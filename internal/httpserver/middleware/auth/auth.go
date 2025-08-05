package auth

import (
	"encoding/base64"
	"github.com/qrv1t9/webdav/internal/webdav/users"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func New(username string, password string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := r.Context().Value("log").(*zap.Logger)

			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Basic ") {
				w.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				log.Info("request without authorization")
				return
			}

			payload, err := base64.StdEncoding.DecodeString(auth[len("Basic "):])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				log.Info("failed to decode payload")
				return
			}

			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				log.Info("decoded payload has invalid format")
				return
			}

			if username != pair[0] {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				log.Info("invalid username")
				return
			}

			p, err := users.CheckPassword(password, pair[1])
			if err != nil {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				log.Info("failed to check password: " + err.Error())
				return
			}
			if p != true {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				log.Info("invalid password")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
