package main

import (
	"errors"
)

var ErrUnsupportedStatusCode = errors.New("unsupported status code")

type Response struct {
	Status     string
	StatusCode int
	Headers    Headers
	Proto      string
	ProtoMajor int
	ProtoMinor int

	Body          []byte
	ContentLength int
}

var codeToStatusMap = map[int]string{
	200: "200 OK",
	400: "400 Bad Request",
	404: "404 Not Found",
}
