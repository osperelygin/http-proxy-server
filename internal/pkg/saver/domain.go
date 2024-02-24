package saver

import (
	"context"
	"net/http"
)

type Saver interface {
	SaveRequest(ctx context.Context, r *http.Request) (int, error)
	SaveResponse(ctx context.Context, requestID int, resp *http.Response) error
}
