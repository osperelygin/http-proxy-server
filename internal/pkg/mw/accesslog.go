package mw

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func AccessLog(logger *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := GetRequestID(r.Context())
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
			"header": r.Header,
			"reqID":  reqID,
		}).Infoln("start request processing")

		next.ServeHTTP(w, r)

		logger.WithFields(logrus.Fields{
			"latency": time.Since(start),
			"reqID":   reqID,
		}).Infoln("request processed")
	})
}
