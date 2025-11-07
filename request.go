package main

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"
)

var (
	CRLF = []byte("\r\n")

	ErrMalformedRequestLine = errors.New("malformed request line")
	ErrInvalidHTTPVersion   = errors.New("invalid http version")
	ErrInvalidMethod        = errors.New("invalid method")
	ErrUnsupportedMethod    = errors.New("http method not yet supported")
	ErrMalformedHeaderLine  = errors.New("field/header line is malformed")
	ErrInvalidContentLength = errors.New("invalid content-length")
	ErrUnsupportedEncoding  = errors.New("unsupported transfer-encoding")
)

type Request struct {
	Method        string
	Target        string
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Headers       Headers
	ContentLength int
	// TODO: Convert Body to an `io.Reader`.
	Body string

	state parseState
}

type parseState string

const (
	maxBufferSize = 8 << 10

	stateRequestLine parseState = "request line"
	stateFieldLines  parseState = "field lines"
	stateBody        parseState = "body"
	stateDone        parseState = "done"
)

// Parse reads from a connection until a request has been successfully parsed
// and returns an error for malformed requests.
func Parse(conn io.Reader) (*Request, error) {
	req := NewRequest()
	buf := make([]byte, maxBufferSize)
	read := 0
	consumed := 0

	// TODO: Add buffer compaction and limit header sizing.
	for {
		n, err := conn.Read(buf[read:])
		if err != nil {
			slog.Error("[ParseRequest] read error", slog.Any("error", err))
			return nil, err
		}
		read += n

		c, done, err := req.ParseRequest(buf[consumed:read])
		if err != nil {
			return nil, err
		}
		consumed += c
		if done {
			return req, nil
		}
	}
}

// ParseRequest returns consumedBytes, isDone, err
func (req *Request) ParseRequest(data []byte) (int, bool, error) {
	consumed := 0

	for {
		switch req.state {
		case stateRequestLine:
			end := bytes.Index(data, CRLF)
			if end == -1 {
				return consumed, false, nil
			}
			requestLine := data[:end]
			consumed += end + len(CRLF)
			if err := parseRequestLine(req, requestLine); err != nil {
				return consumed, false, err
			}
			req.state = stateFieldLines
		case stateFieldLines:
			idx := bytes.Index(data[consumed:], CRLF)
			if idx == -1 {
				return consumed, false, nil
			}
			if idx == 0 {
				if err := req.applyHeaderSemantics(); err != nil {
					return consumed, false, err
				}
				req.state = stateBody
				consumed += len(CRLF)
				continue
			}
			fieldLine := data[consumed : consumed+idx]
			if err := parseFieldLine(req, fieldLine); err != nil {
				return consumed, false, err
			}
			consumed += idx + len(CRLF)
		case stateBody:
			contentLengthValue, ok := req.Headers.Get("content-length")
			if ok {
				contentLength, err := strconv.Atoi(contentLengthValue)
				if err != nil {
					return consumed, false, ErrInvalidContentLength
				}
				// TODO: Handle cases with bigger body than specified.
				if len(data)-consumed < contentLength {
					return consumed, false, nil
				}
				req.Body = string(data[consumed : consumed+contentLength])
				consumed += contentLength
			}

			_, ok = req.Headers.Get("transfer-encoding")
			if ok {
				// TODO: for "Transfer-Encoding: chunked", support reading a chunked body.
				return consumed, false, ErrUnsupportedEncoding
			}

			req.state = stateDone
		case stateDone:
			return consumed, true, nil
		}
	}
}

// applyHeaderSemantics processes all header field lines.
func (req *Request) applyHeaderSemantics() error {
	// TODO: Content-Length can appear multiple times with the same value,
	// so we need to enforce that it's the same for all appearances.
	if contentLength, ok := req.Headers.Get("content-length"); ok {
		n, err := strconv.Atoi(contentLength)
		if err != nil {
			return ErrInvalidContentLength
		}
		req.ContentLength = n
	}
	return nil
}

func parseRequestLine(req *Request, line []byte) error {
	method, rest, ok := bytes.Cut(line, []byte(" "))
	target, proto, ok2 := bytes.Cut(rest, []byte(" "))
	if !ok || !ok2 {
		return ErrMalformedRequestLine
	}

	req.Method = string(method)
	req.Target = string(target)
	req.Proto = string(proto)

	if err := validateMethod(string(method)); err != nil {
		return err
	}

	major, minor, valid := parseHTTPVersion(string(proto))
	req.ProtoMajor = major
	req.ProtoMinor = minor
	if !valid {
		return ErrInvalidHTTPVersion
	}
	return nil
}

func parseHTTPVersion(version string) (int, int, bool) {
	switch version {
	case "HTTP/1.1":
		return 1, 1, true
	case "HTTP/1.0":
		return 1, 0, true
	default:
		return 0, 0, false
	}
}

func validateMethod(method string) error {
	switch method {
	case "GET", "POST":
		return nil
	case "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE":
		return ErrUnsupportedMethod
	default:
		return ErrInvalidMethod
	}
}

func parseFieldLine(req *Request, line []byte) error {
	rawName, rawValue, ok := bytes.Cut(line, []byte(":"))
	if !ok {
		return ErrMalformedHeaderLine
	}
	// TODO: enforce RFC's token set on field name.
	if len(rawName) == 0 {
		return ErrMalformedHeaderLine
	}
	fieldName := string(rawName)
	// TODO: Check if we need to validate the value in any way.
	fieldValue := strings.TrimSpace(string(rawValue))
	req.Headers.Add(fieldName, fieldValue)

	return nil
}

func NewRequest() *Request {
	return &Request{
		Headers: make(map[string][]string),
		state:   stateRequestLine,
	}
}
