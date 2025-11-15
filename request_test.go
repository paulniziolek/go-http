package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		input       string
		ok          bool
		expectedErr error
	}{
		{
			input: "GET /resource HTTP/1.1\r\n\r\n",
			ok:    true,
		},
		{
			input: "POST /resource HTTP/1.1\r\nContent-Length: 5\r\n\r\nHello",
			ok:    true,
		},
		{
			input:       "NOMETHOD / HTTP/1.1\r\n\r\n",
			ok:          false,
			expectedErr: ErrInvalidMethod,
		},
		{
			input:       "GET /HTTP/1.1\r\n\r\n",
			ok:          false,
			expectedErr: ErrMalformedRequestLine,
		},
		{
			input: "GET / HTTP/1.1 extra\r\n\r\n",
			ok:    false,
			// expectedErr here should be ErrMalformedRequestLine
			expectedErr: ErrInvalidHTTPVersion,
		},
		{
			input:       "GET / HTTP/0.1\r\n\r\n",
			ok:          false,
			expectedErr: ErrInvalidHTTPVersion,
		},
		{
			input:       "POST /resource HTTP/1.1\r\nContent-Length: notANumber\r\n\r\nHello",
			ok:          false,
			expectedErr: ErrInvalidContentLength,
		},
		{
			input:       "POST /resource HTTP/1.1\r\nContent-Length 5\r\n\r\nHello",
			ok:          false,
			expectedErr: ErrMalformedHeaderLine,
		},
		{
			input:       "POST /resource HTTP/1.1\r\n:5\r\n\r\n",
			ok:          false,
			expectedErr: ErrMalformedHeaderLine,
		},
		{
			input: "POST /resource HTTP/1.1\r\nContent-Length: 9000\r\n\r\nHello",
			ok:    false,
		},
		{
			input:       "GET / HTTP/1.0\r\nTransfer-Encoding: chunked\r\n\r\n",
			ok:          false,
			expectedErr: ErrUnsupportedEncoding,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			req := NewRequest()
			_, done, err := req.ParseRequest([]byte(tt.input))
			if tt.expectedErr == nil && err != nil {
				t.Fatalf("expected no error but received error: %v", err)
			}
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error: %v, but got: %v", tt.expectedErr, err)
			}
			if tt.ok != done {
				t.Fatalf("expected ok: %t but got: %t", tt.ok, done)
			}
		})
	}
}

func FuzzParse(f *testing.F) {
	seeds := [][]byte{
		[]byte("GET /resource HTTP/1.1\r\n\r\n"),
		[]byte("POST /resource HTTP/1.1\r\nContent-Length: 5\r\n\r\nHello"),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		r := bytes.NewReader(data)
		_, _ = Parse(r)
	})
}
