package proxy

import (
	"io"
	"net/http"
	"strings"
)

func parseFormURLEncoding(r *http.Request) error {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/x-www-form-urlencoded" {
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	r.Body = io.NopCloser(strings.NewReader(r.PostForm.Encode()))

	return nil
}
