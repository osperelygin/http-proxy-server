package proxy

import (
	"context"
	"flag"
	init_postgres "http-proxy-server/internal/init/postgres"
	"http-proxy-server/internal/pkg/logger"
	postgres_saver "http-proxy-server/internal/pkg/saver/repo/postgres"
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

	saver := postgres_saver.New(pool, logger)

	srv := NewProxyServer(cfg, logger, saver)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
