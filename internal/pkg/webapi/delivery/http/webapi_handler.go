package webapi_handler

import (
	"encoding/json"
	domain "http-proxy-server/internal/pkg/webapi"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	usecase domain.Usecase
	logger  *logrus.Logger
}

func New(usecase domain.Usecase, logger *logrus.Logger) Handler {
	return Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h Handler) encodeResponse(body any, w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(body); err != nil {
		h.logger.Errorln("Encode failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infoln("successful encode")
}

func (h Handler) GetAllRequest(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("entered in GetAllRequest")

	requests, err := h.usecase.GetAllRequest(r.Context())
	if err != nil {
		h.logger.Errorln("GetAllRequest failed", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.encodeResponse(requests, w)
}

func (h Handler) GetRequestByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("entered in GetRequestByID")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Errorln("request id parsing from path failed:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request, err := h.usecase.GetRequestByID(r.Context(), id)
	if err != nil {
		h.logger.Errorln("GetRequestByID failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.encodeResponse(request, w)
}

func (h Handler) GetRespByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("entered in GetResponseByID")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Errorln("request id parsing from path failed:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.usecase.GetResponseByRequestID(r.Context(), id)
	if err != nil {
		h.logger.Errorln("GetResponseByRequestID failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.encodeResponse(response, w)
}

func (h Handler) RepeatRequest(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("entered in RepeatRequest")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Errorln("request id parsing from path failed:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := h.usecase.GetHTTPRequest(r.Context(), id)
	if err != nil {
		h.logger.Errorln("get request by id failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := h.usecase.DoRequest(r.Context(), req)
	if err != nil {
		h.logger.Errorln("DoRequest failed:", err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infoln("successful executed RepeatRequest")
}

func (h Handler) ScanRequest(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("entered in RepeatRequest")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Errorln("request id parsing from path failed:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rctx := r.Context()
	req, err := h.usecase.GetRequestByID(rctx, id)
	if err != nil {
		h.logger.Errorln("get http request failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, err := h.usecase.Scan(rctx, req)
	if err != nil {
		h.logger.Errorln("scan failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok {
		w.Write([]byte("!request is vulnerable!\n"))
		return
	}

	w.Write([]byte("vulnerabilies was not found\n"))
}
