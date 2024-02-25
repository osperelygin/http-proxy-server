package webapi

import (
	"context"
	"http-proxy-server/internal/pkg/models"
	"net/http"
)

type Repo interface {
	GetAllRequest(ctx context.Context) ([]models.Request, error)
	GetRequestByID(ctx context.Context, id int) (models.Request, error)
	GetResponseByRequestID(ctx context.Context, reqID int) (models.Response, error)
}

type Usecase interface {
	GetHTTPRequest(ctx context.Context, id int) (*http.Request, error)
	Scan(ctx context.Context, req models.Request) (bool, error)
	DoRequest(ctx context.Context, req *http.Request) (*http.Response, error)
	GetRequestByID(ctx context.Context, id int) (models.Request, error)
	GetAllRequest(ctx context.Context) ([]models.Request, error)
	GetResponseByRequestID(ctx context.Context, reqID int) (models.Response, error)
}
