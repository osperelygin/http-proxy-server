package webapi_handler

import (
	"crypto/tls"
	"encoding/json"
	"http-proxy-server/internal/pkg/converter"
	domain "http-proxy-server/internal/pkg/webapi"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo   domain.Repo
	logger *logrus.Logger
	client http.Client
}

func New(repo domain.Repo, proxyURL string, logger *logrus.Logger) (Handler, error) {
	proxyUrl, err := url.Parse(proxyURL)
	if err != nil {
		return Handler{}, err
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return Handler{
		repo:   repo,
		client: client,
		logger: logger,
	}, nil
}

func (h Handler) GetAllRequest(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("entered in GetAllRequest")

	requests, err := h.repo.GetAllRequest(r.Context())
	if err != nil {
		h.logger.Errorln("GetAllRequest failed", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(requests); err != nil {
		h.logger.Errorln("Encode failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infoln("successful executed GetAllRequest and exited")
}

func (h Handler) GetRequestByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("entered in GetRequestByID")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Errorln("request id parsing from path failed:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request, err := h.repo.GetRequestByID(r.Context(), id)
	if err != nil {
		h.logger.Errorln("GetRequestByID failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(request); err != nil {
		h.logger.Errorln("Encode request failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infoln("successful executed GetRequestByID and exited")
}

func (h Handler) GetRespByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("entered in GetResponseByID")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Errorln("request id parsing from path failed:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.repo.GetResponseByRequestID(r.Context(), id)
	if err != nil {
		h.logger.Errorln("GetResponseByRequestID failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		h.logger.Errorln("Encode request failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infoln("successful executed GetRequestByID and exited")
}

func (h Handler) RepeatRequest(w http.ResponseWriter, r *http.Request) {
	h.logger.Infoln("entered in RepeatRequest")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Errorln("request id parsing from path failed:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reqModel, err := h.repo.GetRequestByID(r.Context(), id)
	if err != nil {
		h.logger.Errorln("GetRequestByID failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	method, reqUrl, body := converter.ModelToRequest(reqModel)
	req, err := http.NewRequest(method, reqUrl, body)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"method": method,
			"url":    reqUrl,
		}).Errorln("NewRequest failed:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header = reqModel.Headers
	req.Header.Add("Cookie", strings.Join(reqModel.Cookies, "\n"))

	resp, err := h.client.Do(req)
	if err != nil {
		h.logger.Errorln("Do failed:", err.Error())
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
