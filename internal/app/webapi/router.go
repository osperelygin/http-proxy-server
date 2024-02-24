package webapi

import (
	"http-proxy-server/internal/pkg/mw"
	webapi_handler "http-proxy-server/internal/pkg/webapi/delivery/http"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func getRouter(handler webapi_handler.Handler, logger *logrus.Logger) http.Handler {
	mux := mux.NewRouter()
	webapi := mux.PathPrefix("/api/v1").Subrouter()
	{
		webapi.HandleFunc("/requests", handler.GetAllRequest).Methods(http.MethodGet)
		webapi.HandleFunc("/request/{id}", handler.GetRequestByID).Methods(http.MethodGet)
		webapi.HandleFunc("/response/{id}", handler.GetRespByID).Methods(http.MethodGet)
		webapi.HandleFunc("/repeat/{id}", handler.RepeatRequest).Methods(http.MethodGet)
	}
	router := mw.AccessLog(logger, webapi)
	return mw.RequestID(router)
}
