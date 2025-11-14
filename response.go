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

const (
	StatusOK         = 200
	StatusBadRequest = 400
	StatusNotFound   = 404
)

var codeToStatusMap = map[int]string{
	StatusOK:         "200 OK",
	StatusBadRequest: "400 Bad Request",
	StatusNotFound:   "404 Not Found",
}
