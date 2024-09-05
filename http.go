package shuai_utils

import (
	"compress/gzip"
	"io"
	"net/http"
)

func HttpRequestGetBody(req *http.Request) ([]byte, error) {
	var (
		reader  io.ReadCloser
		err     error
		reqBody []byte
	)
	switch req.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(req.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = req.Body
	}
	defer reader.Close()
	reqBody, err = io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return reqBody, nil
}
