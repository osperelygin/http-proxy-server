package converter

import (
	"bytes"
	"http-proxy-server/internal/pkg/models"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	kvDelimiter     = ":"
	valueDelimiter  = ","
	headerDelimiter = "\n"
)

func MapToStr(m map[string][]string) string {
	var str string
	for key, values := range m {
		str += key + kvDelimiter
		for _, value := range values {
			str += value + valueDelimiter
		}

		str = strings.TrimRight(str, valueDelimiter) + headerDelimiter
	}

	return str
}

func StrToMap(str string) map[string][]string {
	res := make(map[string][]string)

	for _, header := range strings.Split(str, headerDelimiter) {
		idx := strings.Index(header, kvDelimiter)
		if idx == -1 {
			continue
		}

		key := header[:idx]
		res[key] = append(res[key], strings.Split(header[idx+1:], valueDelimiter)...)
	}

	return res
}

func CookieToStr(cookie []*http.Cookie) string {
	var str string
	for _, val := range cookie {
		str += val.String() + headerDelimiter
	}

	return str
}

func StrToCookie(str string) []string {
	return strings.Split(str, headerDelimiter)
}

func ModelToRequest(req models.Request) (string, string, io.Reader) {
	reqURL := url.URL{
		Scheme:   req.Scheme,
		Host:     req.Host,
		Path:     req.Path,
		RawQuery: url.Values(req.QueryParams).Encode(),
	}

	return req.Method, reqURL.String(), bytes.NewBufferString(req.Body)
}
