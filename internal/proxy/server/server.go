package server

import (
	"http-proxy-server/internal/pkg/config"
	"http-proxy-server/internal/pkg/mw"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ProxyServer struct {
	port   string
	host   string
	logger *logrus.Logger
}

func New(cfg config.SrvConfig, logger *logrus.Logger) *ProxyServer {
	return &ProxyServer{
		port:   cfg.Port,
		host:   cfg.Host,
		logger: logger,
	}
}

func (ps ProxyServer) ListenAndServe() error {
	h := mw.AccessLog(ps.logger, http.HandlerFunc(ps.proxyHTTP))
	h = mw.RequestID(h)
	server := http.Server{
		Addr:    ps.host + ":" + ps.port,
		Handler: h,
	}

	ps.logger.Infof("start listening at %s:%s", ps.host, ps.port)

	return server.ListenAndServe()
}

func (ps ProxyServer) proxyHTTP(w http.ResponseWriter, r *http.Request) {
	reqID := mw.GetRequestID(r.Context())
	ps.logger.WithField("reqID", reqID).Infoln("entered in proxyHTTP")

	r.Header.Del("Proxy-Connection")

	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		ps.logger.WithField("reqID", reqID).Errorln("round trip failed:", err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer res.Body.Close()

	res.Cookies()
	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(res.StatusCode)
	if _, err := io.Copy(w, res.Body); err != nil {
		ps.logger.WithField("reqID", reqID).Errorln("io copy failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ps.logger.WithField("reqID", reqID).Infoln("exited from proxyHTTP")
}
