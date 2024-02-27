package postgres_saver

import (
	"bytes"
	"context"
	"http-proxy-server/internal/pkg/converter"
	"http-proxy-server/internal/pkg/models"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type PostgresSaver struct {
	pool   models.IPgxPool
	logger *logrus.Logger
}

func New(pool models.IPgxPool, logger *logrus.Logger) *PostgresSaver {
	return &PostgresSaver{
		pool:   pool,
		logger: logger,
	}
}

func copyHeaders(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		out[key] = values
	}

	return out
}

func (db *PostgresSaver) SaveRequest(ctx context.Context, r *http.Request) (int, error) {
	db.logger.Infoln("entered in SaveRequest")

	var id int
	body, err := io.ReadAll(r.Body)
	if err != nil {
		db.logger.Errorln("ReadAll failed:", err.Error())
		return id, err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	cookies := converter.CookieToStr(r.Cookies())
	queryStrig := converter.MapToStr(r.URL.Query())
	postParams := converter.MapToStr(r.PostForm)

	copyHeaders := copyHeaders(r.Header)
	delete(copyHeaders, "Cookie")
	headers := converter.MapToStr(copyHeaders)

	query := "insert into request (method, scheme, host, path, query_string, post_params, cookies, headers, body) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id"
	row := db.pool.QueryRow(ctx, query, r.Method, r.URL.Scheme, r.URL.Host, r.URL.Path, queryStrig, postParams, cookies, headers, body)
	if err := row.Scan(&id); err != nil {
		db.logger.WithFields(logrus.Fields{
			"Method":  r.Method,
			"Scheme":  r.URL.Scheme,
			"Host":    r.URL.Host,
			"Path":    r.URL.Path,
			"Cookies": cookies,
			"Headers": headers,
		}).Errorln("Scan failed:", err.Error())

		return id, err
	}

	db.logger.Infoln("successful executed SaveReqeust and exited")

	return id, nil
}

func (db *PostgresSaver) SaveResponse(ctx context.Context, reqID int, resp *http.Response) error {
	db.logger.Infoln("entered in SaveResponse")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		db.logger.Errorln("ReadAll failed:", err.Error())
		return err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	cookies := converter.CookieToStr(resp.Cookies())

	copyHeaders := copyHeaders(resp.Header)
	delete(copyHeaders, "Set-Cookie")
	headers := converter.MapToStr(copyHeaders)

	query := "insert into response (request_id, code, cookies, headers, body) values ($1, $2, $3, $4, $5)"
	_, err = db.pool.Exec(ctx, query, reqID, resp.StatusCode, cookies, headers, body)
	if err != nil {
		db.logger.WithFields(logrus.Fields{
			"reqID":   reqID,
			"Code":    resp.StatusCode,
			"Cookies": cookies,
			"headers": headers,
		}).Errorln("Exec failed:", err.Error())

		return err
	}

	db.logger.Infoln("successful executed SaveReqeust and exited")

	return nil
}
