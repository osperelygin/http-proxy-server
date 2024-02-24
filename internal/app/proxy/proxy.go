package proxy

import (
	"bufio"
	"crypto/tls"
	"http-proxy-server/internal/pkg/mw"
	"io"
	"net/http"
	"net/http/httputil"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func (ps ProxyServer) hostTLSConfig(host string) (*tls.Config, error) {
	if err := exec.Command(ps.cfg.Script, host).Run(); err != nil {
		ps.logger.WithFields(logrus.Fields{
			"script": ps.cfg.Script,
			"host":   host,
		}).Errorln("exec command failed:", err.Error())

		return nil, err
	}

	tlsCert, err := tls.LoadX509KeyPair(ps.cfg.CertFile, ps.cfg.KeyFile)
	if err != nil {
		ps.logger.WithFields(logrus.Fields{
			"cert file": ps.cfg.CertFile,
			"key file":  ps.cfg.KeyFile,
		}).Errorln("LoadX509KeyPair failed:", err.Error())

		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}, nil
}

func (ps ProxyServer) proxyHTTP(w http.ResponseWriter, r *http.Request) {
	reqID := mw.GetRequestID(r.Context())

	ps.logger.WithField("reqID", reqID).Infoln("entered in proxyHTTP")

	r.Header.Del("Proxy-Connection")

	if err := parseFormURLEncoding(r); err != nil {
		ps.logger.Errorln("parseFormURLEncoding failed:", err.Error())
	}

	id := ps.saveRequest(r)

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	resp.Cookies()
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	ps.saveResponse(resp, id)

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		ps.logger.WithField("reqID", reqID).Errorln("io.Copy failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ps.logger.WithField("reqID", reqID).Infoln("exited from proxyHTTP")
}

func (ps ProxyServer) proxyHTTPS(w http.ResponseWriter, r *http.Request) {
	reqID := mw.GetRequestID(r.Context())
	ps.logger.WithField("reqID", reqID).Infoln("entered in proxyHTTPS")

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		ps.logger.WithField("reqID", reqID).Errorln("hijacking not supported")
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	localConn, _, err := hijacker.Hijack()
	if err != nil {
		ps.logger.WithField("reqID", reqID).Errorln("hijack failed:", err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	if _, err := localConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n")); err != nil {
		ps.logger.WithField("reqID", reqID).Errorln("write to local connection failed:", err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		localConn.Close()
		return
	}

	defer localConn.Close()

	tlsConfig, err := ps.hostTLSConfig(strings.Split(r.Host, ":")[0])
	if err != nil {
		ps.logger.WithField("reqID", reqID).Errorln("hostTLSConfig failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tlsLocalConn := tls.Server(localConn, tlsConfig)
	if err := tlsLocalConn.Handshake(); err != nil {
		ps.logger.WithField("reqID", reqID).Errorln("tls handshake failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer tlsLocalConn.Close()

	reader := bufio.NewReader(tlsLocalConn)
	request, err := http.ReadRequest(reader)
	if err != nil {
		ps.logger.WithField("reqID", reqID).Errorln("read request failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	remoteConn, err := tls.Dial("tcp", r.Host, tlsConfig)
	if err != nil {
		ps.logger.WithField("reqID", reqID).Errorln("tls dial failed:", err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer remoteConn.Close()

	request.URL.Scheme = "https"
	request.URL.Host = r.URL.Host
	id := ps.saveRequest(request)

	request.Header.Set("Accept-Encoding", "identity;q=0")

	requestByte, err := httputil.DumpRequest(request, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = remoteConn.Write(requestByte)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serverReader := bufio.NewReader(remoteConn)
	resp, err := http.ReadResponse(serverReader, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawResponse, err := httputil.DumpResponse(resp, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tlsLocalConn.Write(rawResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ps.saveResponse(resp, id)

	ps.logger.WithField("reqID", reqID).Infoln("exited from proxyHTTPS")
}
