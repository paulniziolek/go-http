package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrWriteOverflow = errors.New("attempted to write more than content-length specified")
)

// The ResponseWriter is inspired from the net/http package.
type ResponseWriter interface {
	// Header() returns the Headers map of the Response. Headers are written
	// after a call to Write or WriteHeader. Modifying the headers after
	// will be ignored. See Headers.
	Header() Headers
	// WriteHeader() will write the headers of the response to the client along
	// with a given status code.
	WriteHeader(code int)
	// Write() will write the response body to the client. If Write is called
	// before WriteHeader, then Write will also write the current Header to the
	// client along with the assumption of a "200 OK" response.
	Write([]byte) (int, error)
}

// http1ResponseWriter is the default ResponseWriter for HTTP/1.1.
type http1ResponseWriter struct {
	req    *Request
	header Headers
	writer io.Writer

	contentLength    int
	hasContentLength bool
	status           string
	wroteHeader      bool
	wrote            int
	chunked          bool
}

func (w *http1ResponseWriter) Header() Headers {
	return w.header
}

func (w *http1ResponseWriter) WriteHeader(code int) {
	status, ok := codeToStatusMap[code]
	if !ok {
		slog.Error("[WriteHeader] Invalid or unsupported HTTP code", slog.Int("code", code))
		return
	}
	if w.wroteHeader {
		slog.Error("[WriteHeader] Header has already been written")
		return
	}
	if code >= 100 && code < 200 || code == 204 || code == 304 {
		// By HTTP Semantics, we must refuse a body in these responses.
		w.Header().Set("Content-Length", "0")
	}

	if cl, ok := w.header.Get("Content-Length"); ok {
		// TODO: Enforce all values to be the same.
		parsedCL, err := strconv.Atoi(cl)
		if err != nil {
			slog.Error("[WriteHeader] Invalid Content-Length set", slog.String("content-length", cl))
			return
		}
		w.contentLength = parsedCL
		w.hasContentLength = true
	}
	proto := "HTTP/1.1"
	if w.req != nil && w.req.Proto != "" {
		proto = w.req.Proto
	}

	if proto == "HTTP/1.0" {
		if w.req != nil && w.req.Headers.ContainsValue("Connection", "keep-alive") {
			w.Header().Set("Connection", "keep-alive")
		} else {
			w.Header().Set("Connection", "close")
		}
	}

	if !w.header.ContainsKey("Date") {
		w.header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	}

	// For HTTP/1.1 requests without any set Content-Length, we must default
	// response bodies to chunked encoding.
	if proto == "HTTP/1.1" && !w.hasContentLength && !w.Header().ContainsValue("Transfer-Encoding", "chunked") {
		w.Header().Add("Transfer-Encoding", "chunked")
	}

	w.chunked = w.Header().ContainsValue("Transfer-Encoding", "chunked")

	w.status = status
	fmt.Fprintf(w.writer, "%s %s", proto, w.status)
	w.writer.Write(CRLF)

	w.header.ForEach(func(key string, value string) {
		w.writer.Write([]byte(fmt.Sprintf("%s: %s", key, value)))
		w.writer.Write(CRLF)
	})
	w.writer.Write(CRLF)
	w.wroteHeader = true
}

func (w *http1ResponseWriter) Write(b []byte) (int, error) {
	// TODO: Support writing in chunked mode.
	if !w.wroteHeader {
		w.WriteHeader(StatusOK)
	}
	remaining := w.contentLength - w.wrote
	if len(b) > remaining {
		slog.Error("[Write] Attempted to write more than specified Content-Length")
		return 0, ErrWriteOverflow
	}
	n, err := w.writer.Write(b)
	if err != nil {
		return n, err
	}
	w.wrote += n
	return n, nil
}
