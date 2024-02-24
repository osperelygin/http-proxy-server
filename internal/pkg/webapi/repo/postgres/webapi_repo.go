package webapi_repo

import (
	"context"
	"http-proxy-server/internal/pkg/converter"
	"http-proxy-server/internal/pkg/models"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type WebapiRepo struct {
	pool   models.IPgxPool
	logger *logrus.Logger
}

func New(pool models.IPgxPool, logger *logrus.Logger) *WebapiRepo {
	return &WebapiRepo{
		pool:   pool,
		logger: logger,
	}
}

const (
	queryRequests = "select id, method, scheme, host, path, query_string, post_params, cookies, headers, encode(body, 'escape') from request"
)

func (db *WebapiRepo) scanRequest(row pgx.Row) (models.Request, error) {
	var req models.Request
	var headers, cookies, queryString, postParams string

	err := row.Scan(&req.Id, &req.Method, &req.Scheme, &req.Host, &req.Path, &queryString, &postParams, &cookies, &headers, &req.Body)
	if err != nil {
		db.logger.Errorln("Scan failed:", err.Error())
		return req, err
	}

	req.Headers = converter.StrToMap(headers)
	req.QueryParams = converter.StrToMap(queryString)
	req.Cookies = converter.StrToCookie(cookies)
	req.PostParams = converter.StrToMap(postParams)

	return req, nil
}

func (db *WebapiRepo) scanResponse(row pgx.Row) (models.Response, error) {
	var resp models.Response
	var headers, cookies string

	err := row.Scan(&resp.Id, &resp.RequestId, &resp.Code, &cookies, &headers, &resp.Body)
	if err != nil {
		db.logger.Errorln("Scan failed:", err.Error())
		return resp, err
	}

	resp.Headers = converter.StrToMap(headers)
	resp.Cookies = converter.StrToCookie(cookies)

	return resp, nil
}

func (db *WebapiRepo) GetAllRequest(ctx context.Context) ([]models.Request, error) {
	db.logger.Infoln("entered in GetAllRequests")

	rows, err := db.pool.Query(ctx, queryRequests)
	if err != nil {
		db.logger.Errorln("Query failed:", err.Error())
		return nil, err
	}

	defer rows.Close()

	requests := make([]models.Request, 0)
	for rows.Next() {
		req, err := db.scanRequest(rows)
		if err != nil {
			return nil, err
		}

		requests = append(requests, req)
	}

	db.logger.Infoln("successful executed GetAllRequest and exited")

	return requests, nil
}

func (db *WebapiRepo) GetRequestByID(ctx context.Context, id int) (models.Request, error) {
	db.logger.Infoln("entered in GetRequestByID")

	query := queryRequests + " where id = $1;"
	req, err := db.scanRequest(db.pool.QueryRow(ctx, query, id))
	if err != nil {
		return models.Request{}, err
	}

	db.logger.Infoln("successful executed GetRequestByID and exited")

	return req, nil
}

func (db *WebapiRepo) GetResponseByRequestID(ctx context.Context, reqID int) (models.Response, error) {
	db.logger.Infoln("entered in GetRequestByID")

	query := "select id, request_id, code, cookies, headers, encode(body, 'escape') from response where id = $1;"
	resp, err := db.scanResponse(db.pool.QueryRow(ctx, query, reqID))
	if err != nil {
		return models.Response{}, err
	}

	db.logger.Infoln("successful executed GetRequestByID and exited")

	return resp, nil
}
