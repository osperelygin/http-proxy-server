package converter

import (
	"bufio"
	"bytes"
	"io"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

func ToUtf8Encoding(in io.Reader) (io.Reader, error) {
	body, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	peek, err := bufio.NewReader(bytes.NewReader(body)).Peek(256)
	if err != nil {
		return nil, err
	}

	encoder, _, _ := charset.DetermineEncoding(peek, "")

	return transform.NewReader(bytes.NewReader(body), encoder.NewDecoder()), nil
}
