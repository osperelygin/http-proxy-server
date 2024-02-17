package main

import (
	"http-proxy-server/internal/pkg/app"
	"http-proxy-server/internal/pkg/logger"
	"http-proxy-server/internal/proxy/config"
	"http-proxy-server/internal/proxy/server"
)

var loggerSingleton logger.Singleton

func main() {
	app := app.Init()

	cfg := config.GetConfig(app.ConfigPath)
	logger := loggerSingleton.GetLogger()

	srv := server.New(cfg, logger)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalln(err.Error())
	}
}
