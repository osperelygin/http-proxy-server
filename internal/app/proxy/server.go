package proxy

import (
	"http-proxy-server/internal/pkg/saver"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ProxyServer struct {
	cfg    Config
	saver  saver.Saver
	logger *logrus.Logger
}

func NewProxyServer(cfg Config, logger *logrus.Logger, saver saver.Saver) *ProxyServer {
	return &ProxyServer{
		cfg:    cfg,
		saver:  saver,
		logger: logger,
	}
}

func (ps ProxyServer) ListenAndServe() error {
	server := http.Server{
		Addr:    ":" + ps.cfg.Port,
		Handler: ps.handler(),
	}

	ps.logger.Infof("start listening at :%s", ps.cfg.Port)
	return server.ListenAndServe()
}
