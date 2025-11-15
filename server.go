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

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}

type Server struct {
	Addr   string
	Router map[string]Handler
}

func NewServer(address string) *Server {
	return &Server{
		Addr:   address,
		Router: make(map[string]Handler),
	}
}

func (s *Server) Handle(pattern string, handler Handler) {
	if _, ok := s.Router[pattern]; ok {
		slog.Error("[Handle] Handler for pattern already defined", slog.String("pattern", pattern))
		return
	}
	s.Router[pattern] = handler
}

func (s *Server) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	s.Handle(pattern, HandlerFunc(handler))
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
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	// TODO: Use configs to set HTTP timeouts
	// TODO: Process HTTP requests in a loop until reached Closed request/response.
	conn.SetReadDeadline(time.Now().Add(defaultReadTimeout))
	req, err := Parse(conn)
	if err != nil {
		slog.Error("[handleConn] ParseRequest failed", slog.Any("err", err))
		return
	}
	slog.Info("Received Request", slog.Any("request", req))
	// TODO: ServeHTTP should use default writers based on the HTTP protocol.
	// Currently, only HTTP/1.1 is supported so that is the defaulted protocol.
	w := &http1ResponseWriter{
		req:    req,
		writer: conn,
		header: make(map[string][]string),
	}

	handler, ok := s.Router[req.Target]
	if !ok {
		// The requested target resource doesn't exist.
		slog.Error("Requested target resource is not found", slog.String("target", req.Target))
		w.WriteHeader(StatusNotFound)
		return
	}

	handler.ServeHTTP(w, req)
}
