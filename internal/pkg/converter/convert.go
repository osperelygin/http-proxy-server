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
	cookieDelimiter = "\r\r"
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
		str += val.String() + cookieDelimiter
	}

	return str
}

func StrToCookie(str string) []string {
	cookies := strings.Split(str, cookieDelimiter)
	if len(cookies) == 1 && cookies[0] == str {
		return nil
	}

	return cookies
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
