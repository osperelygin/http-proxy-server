package webapi

import (
	"context"
	"http-proxy-server/internal/pkg/models"
)

type Repo interface {
	GetAllRequest(ctx context.Context) ([]models.Request, error)
	GetRequestByID(ctx context.Context, id int) (models.Request, error)
	GetResponseByRequestID(ctx context.Context, reqID int) (models.Response, error)
}
