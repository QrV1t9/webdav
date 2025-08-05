package httpserver

import (
	"fmt"
	"github.com/qrv1t9/webdav/internal/config"
	"github.com/qrv1t9/webdav/internal/httpserver/middleware/auth"
	"github.com/qrv1t9/webdav/internal/httpserver/middleware/logger"
	"go.uber.org/zap"
	"golang.org/x/net/webdav"
	"net/http"
	"path"
)

type Server struct {
	log *zap.Logger
	cfg config.Config
	Dav *webdav.Handler
}

func New(log *zap.Logger, dav *webdav.Handler, cfg config.Config) *Server {
	return &Server{log: log, Dav: dav, cfg: cfg}
}

func (s *Server) Run() {
	handler := logger.New(s.log)(auth.New(s.cfg.User.Username, s.cfg.User.Password)(s.Dav))

	if s.cfg.TLS {
		cert := path.Join(s.cfg.Path, "cert.pem")
		key := path.Join(s.cfg.Path, "key.pem")
		err := http.ListenAndServeTLS(fmt.Sprintf(":%v", s.cfg.Port), cert, key, handler)
		if err != nil {
			panic("failed to start http server: " + err.Error())
		}
	} else {
		err := http.ListenAndServe(fmt.Sprintf(":%v", s.cfg.Port), handler)
		if err != nil {
			panic("failed to start http server: " + err.Error())
		}
	}
}
