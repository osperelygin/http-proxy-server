package proxy

import (
	"http-proxy-server/internal/pkg/mw"
	"net/http"
)

func (ps ProxyServer) handler() http.Handler {
	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			ps.proxyHTTPS(w, r)
			return
		}

		ps.proxyHTTP(w, r)
	})

	handler := mw.AccessLog(ps.logger, router)
	return mw.RequestID(handler)
}
