package models

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Request struct {
	Id          int                 `json:"id"`
	Method      string              `json:"method"`
	Scheme      string              `json:"scheme"`
	Host        string              `json:"host"`
	Path        string              `json:"path"`
	QueryParams map[string][]string `json:"query_string,omitempty"`
	PostParams  map[string][]string `json:"post_params,omitempty"`
	Cookies     []string            `json:"cookies,omitempty"`
	Headers     map[string][]string `json:"header,omitempty"`
	Body        string              `json:"body,omitempty"`
}

type Response struct {
	Id        int                 `json:"id"`
	RequestId int                 `json:"request_id"`
	Code      int                 `json:"code"`
	Cookies   []string            `json:"cookies,omitempty"`
	Headers   map[string][]string `json:"header,omitempty"`
	Body      string              `json:"body,omitempty"`
}

type IPgxPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}
