package proxy

import (
	"io"
	"net/http"
	"strings"
)

func parseFormURLEncoding(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	r.Body = io.NopCloser(strings.NewReader(r.PostForm.Encode()))

	return nil
}
