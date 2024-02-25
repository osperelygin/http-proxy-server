package webapi_usecase

import (
	"context"
	"crypto/tls"
	"http-proxy-server/internal/pkg/converter"
	"http-proxy-server/internal/pkg/models"
	"io"
	"net/http"
	"net/url"
	"strings"

	domain "http-proxy-server/internal/pkg/webapi"

	"github.com/sirupsen/logrus"
)

var scanCommands = []string{
	";cat /etc/passwd;",
	"|cat /etc/passwd|",
	"`cat /etc/passwd`",
}

type Usecase struct {
	repo   domain.Repo
	client http.Client
	logger *logrus.Logger
}

func New(repo domain.Repo, proxyURL string, logger *logrus.Logger) (*Usecase, error) {
	proxyUrl, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return &Usecase{
		client: client,
		repo:   repo,
		logger: logger,
	}, nil
}

func (uc *Usecase) convertRequest(reqModel models.Request) (*http.Request, error) {
	method, reqUrl, body := converter.ModelToRequest(reqModel)
	req, err := http.NewRequest(method, reqUrl, body)
	if err != nil {
		uc.logger.WithFields(logrus.Fields{
			"method": method,
			"url":    reqUrl,
		}).Errorln("NewRequest failed:", err.Error())

		return nil, err
	}

	req.Header = reqModel.Headers
	if reqModel.Cookies != nil {
		req.Header.Add("Cookie", strings.Join(reqModel.Cookies, "\n"))
	}

	return req, nil
}

func (uc *Usecase) isCommnadInjected(reader io.ReadCloser) (bool, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(body), "root:"), nil
}

func (uc *Usecase) scan(ctx context.Context, req models.Request) (bool, error) {
	r, err := uc.convertRequest(req)
	if err != nil {
		uc.logger.Errorln("convert request failed:", err.Error())
		return false, err
	}

	resp, err := uc.DoRequest(ctx, r)
	if err != nil {
		uc.logger.Errorln("do request failed:", err.Error())
		return false, err
	}

	ok, err := uc.isCommnadInjected(resp.Body)
	if err != nil {
		uc.logger.Errorln("is command inject failed:", err.Error())
		return false, err
	}

	return ok, nil
}

func (uc *Usecase) GetHTTPRequest(ctx context.Context, id int) (*http.Request, error) {
	reqModel, err := uc.repo.GetRequestByID(ctx, id)
	if err != nil {
		uc.logger.Errorln("GetRequestByID failed:", err.Error())
		return nil, err
	}

	req, err := uc.convertRequest(reqModel)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (uc *Usecase) Scan(ctx context.Context, req models.Request) (bool, error) {
	for key, values := range req.QueryParams {
		for idx := range values {
			for _, cmd := range scanCommands {
				req.QueryParams[key][idx] += cmd
				ok, err := uc.scan(ctx, req)
				req.QueryParams[key][idx] = strings.TrimSuffix(req.QueryParams[key][idx], cmd)
				if err != nil {
					return false, err
				}

				if ok {
					return ok, nil
				}
			}
		}
	}

	for key, values := range req.PostParams {
		for idx := range values {
			for _, cmd := range scanCommands {
				req.PostParams[key][idx] += cmd
				ok, err := uc.scan(ctx, req)
				req.PostParams[key][idx] = strings.TrimSuffix(req.PostParams[key][idx], cmd)
				if err != nil {
					return false, err
				}

				if ok {
					return ok, nil
				}
			}
		}
	}

	for key, values := range req.Headers {
		for idx := range values {
			for _, cmd := range scanCommands {
				req.Headers[key][idx] += cmd
				ok, err := uc.scan(ctx, req)
				req.Headers[key][idx] = strings.TrimSuffix(req.Headers[key][idx], cmd)
				if err != nil {
					return false, err
				}

				if ok {
					return ok, nil
				}
			}
		}
	}

	for idx := range req.Cookies {
		for _, cmd := range scanCommands {
			req.Cookies[idx] += cmd
			ok, err := uc.scan(ctx, req)
			req.Cookies[idx] = strings.TrimSuffix(req.Cookies[idx], cmd)
			if err != nil {
				return false, err
			}

			if ok {
				return ok, nil
			}
		}
	}

	return false, nil
}

func (uc *Usecase) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := uc.client.Do(req)
	if err != nil {
		uc.logger.WithFields(logrus.Fields{
			"cookies": req.Cookies(),
			"headers": req.Header,
		}).Errorln("Do failed:", err.Error())
		return nil, err
	}

	return resp, nil
}

func (uc *Usecase) GetRequestByID(ctx context.Context, id int) (models.Request, error) {
	return uc.repo.GetRequestByID(ctx, id)
}

func (uc *Usecase) GetAllRequest(ctx context.Context) ([]models.Request, error) {
	return uc.repo.GetAllRequest(ctx)
}

func (uc *Usecase) GetResponseByRequestID(ctx context.Context, reqID int) (models.Response, error) {
	return uc.repo.GetResponseByRequestID(ctx, reqID)
}
