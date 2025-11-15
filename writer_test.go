package main

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestHTTP1ResponseWriter_HelloWorld(t *testing.T) {
	var buf bytes.Buffer
	w := http1ResponseWriter{
		writer: &buf,
		header: make(Headers),
	}

	err := testHelloWorld(&w)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got := buf.String()

	if !strings.HasPrefix(got, "HTTP/1.1 200 OK\r\n") {
		t.Fatalf("missing or invalid status line, got: %q", got)
	}

	if !strings.Contains(strings.ToLower(got), "content-length: 11") {
		t.Fatalf("missing Content-Length header, got: %q", got)
	}

	if !strings.HasSuffix(got, "hello world") {
		t.Fatalf("missing or invalid body, got: %q", got)
	}
}

func TestHTTP1ResponseWriter_MethodNotAllowed(t *testing.T) {
	var buf bytes.Buffer
	w := http1ResponseWriter{
		writer: &buf,
		header: make(Headers),
	}

	err := testMethodNotAllowed(&w)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got := buf.String()

	if !strings.HasPrefix(got, "HTTP/1.1 405 Method Not Allowed\r\n") {
		t.Fatalf("missing or invalid status line, got: %q", got)
	}
}

func TestHTTP1ResponseWriter_Overflow(t *testing.T) {
	var buf bytes.Buffer
	w := http1ResponseWriter{
		writer: &buf,
		header: make(Headers),
	}

	err := testOverflow(&w)
	if !errors.Is(err, ErrWriteOverflow) {
		t.Fatalf("expected ErrWriteOverflow, got: %v", err)
	}
}

func testHelloWorld(w ResponseWriter) error {
	body := []byte("hello world")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	_, err := w.Write(body)
	return err
}

func testMethodNotAllowed(w ResponseWriter) error {
	w.WriteHeader(StatusMethodNotAllowed)
	return nil
}

func testOverflow(w ResponseWriter) error {
	body := []byte("overflow")
	w.Header().Set("Content-Length", "1")
	_, err := w.Write(body)
	return err
}
