package main

import (
	"log/slog"
	"net"
	"time"
)

const (
	defaultReadTimeout = 500 * time.Millisecond
)

// Handler provides the interface to process a request via `ServeHTTP`.
// This definition is inspired from the net/http package.
type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}

type Server struct {
	Addr   string
	Router map[string]Handler
}

func NewServer(address string) *Server {
	return &Server{
		Addr: address,
	}
}

func (s *Server) ListenAndServe() error {
	var addr = s.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("connection refused", slog.Any("err", err))
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	// TODO: Use configs to set HTTP timeouts
	conn.SetReadDeadline(time.Now().Add(defaultReadTimeout))
	req, err := Parse(conn)
	if err != nil {
		slog.Error("[handleConn] ParseRequest failed", slog.Any("err", err))
		return
	}
	slog.Info("Received Request", slog.Any("request", req))

}
