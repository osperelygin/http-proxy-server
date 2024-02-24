package webapi

import (
	"context"
	"flag"
	init_postgres "http-proxy-server/internal/init/postgres"
	"http-proxy-server/internal/pkg/logger"
	webapi_handler "http-proxy-server/internal/pkg/webapi/delivery/http"
	webapi_repo "http-proxy-server/internal/pkg/webapi/repo/postgres"
)

var loggerSingleton logger.Singleton

func Start() error {
	var cfgPath string
	flag.StringVar(&cfgPath, "c", "configs/config.yaml", "path to config file")

	cfg, err := getConfig(cfgPath)
	if err != nil {
		return err
	}

	pool, err := init_postgres.Init(context.Background(), "DATABASE_URL")
	if err != nil {
		return err
	}

	logger := loggerSingleton.GetLogger()
	repo := webapi_repo.New(pool, logger)
	handler, err := webapi_handler.New(repo, cfg.ProxyURL, logger)
	if err != nil {
		return err
	}

	router := getRouter(handler, logger)

	srv := NewServer(cfg, logger, router)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
