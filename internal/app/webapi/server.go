package webapi

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg     Config
	logger  *logrus.Logger
	handler http.Handler
}

func NewServer(cfg Config, logger *logrus.Logger, handler http.Handler) Server {
	return Server{
		cfg:     cfg,
		logger:  logger,
		handler: handler,
	}
}

func (s Server) ListenAndServe() error {
	server := http.Server{
		Addr:    ":" + s.cfg.Port,
		Handler: s.handler,
	}

	s.logger.Infof("start listening at :%s", s.cfg.Port)
	return server.ListenAndServe()
}
