package main

import (
	"flag"
	"fmt"
	"github.com/qrv1t9/webdav/internal/config"
	"github.com/qrv1t9/webdav/internal/httpserver"
	"github.com/qrv1t9/webdav/internal/webdav"
	"go.uber.org/zap"
	"os"
)

func main() {
	cfg := config.MustLoad(fetchConfigPath())

	log, err := initLogger(cfg.AppEnv)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	defer log.Sync()

	log.Info("webdav server is starting", zap.Int("port", cfg.Port), zap.String("path", cfg.Path), zap.Bool("TLS", cfg.TLS))

	dav := webdav.New(cfg.Prefix, cfg.Path)

	server := httpserver.New(log, dav, cfg)
	server.Run()
}

func initLogger(env string) (*zap.Logger, error) {
	switch env {
	case "development":
		return zap.NewDevelopment()
	case "production":
		return zap.NewProduction()
	}
	return nil, fmt.Errorf("unknown environment %q", env)
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
