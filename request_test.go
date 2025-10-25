package main

import (
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
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			req := NewRequest()
			_, done, err := req.ParseRequest([]byte(tt.input))
			if tt.ok != done {
				t.Fatalf("expected ok: %t but got: %t", tt.ok, done)
			}
			if !errors.Is(tt.expectedErr, err) {
				t.Fatalf("expected error: %v, but got: %v", tt.expectedErr, err)
			}
		})
	}
}
