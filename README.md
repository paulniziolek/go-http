# go-http

A small, personal HTTP/1.1 server implementation in Go. This repository implements a minimal HTTP parser, a tiny router and server loop, a simple ResponseWriter inspired by net/http, and basic header utilities. It is intended as an educational and experimental implementation rather than a production-ready server.

## Features

- Minimal HTTP/1.0 and HTTP/1.1 request parsing
    - Request line parsing, header parsing, Content-Length handling
    - Supports GET and POST (other methods return either unsupported or invalid)
- Simple server & routing
    - Server type with Register/Handle/HandleFunc style routing
- Lightweight ResponseWriter for HTTP/1.1
    - Header management and content-length based writes

## Quick usage (programmatic)

This shows how to create a server and register handlers. This is a minimal sketch — see code for full behavior and details.

```go
srv := NewServer(":8080")

srv.HandleFunc("/hello", func(w ResponseWriter, r *Request) {
    body := "hello"
    w.Header().Set("Content-Length", strconv.Itoa(len(body)))
    w.Write([]byte(body))
})

if err := srv.ListenAndServe(); err != nil {
    panic(err)
}
```

Notes:
- Handlers receive a ResponseWriter and the parsed *Request.
- ResponseWriter requires you to set a Content-Length header if you want to avoid write errors (chunked encoding is not implemented).
- For unsupported/unknown routes use WriteHeader with an appropriate status constant.

## API highlights

- type Server
    - NewServer(address string) *Server
    - Handle(pattern string, handler Handler)
    - HandleFunc(pattern string, func(ResponseWriter, *Request))
    - ListenAndServe() error — accepts and handles connections

- type ResponseWriter
    - Header() Headers
    - WriteHeader(code int)
    - Write([]byte) (int, error)
    - Available status codes: StatusOK, StatusBadRequest, StatusNotFound, StatusMethodNotAllowed; unsupported codes are ignored by the writer.

- type Handler
    - ServeHTTP(w ResponseWriter, r *Request)

- type HandlerFunc

- type Headers
    - Add, Set, Get, GetAll, ForEach; header keys are canonicalized (case-insensitive).

## Limitations / Known gaps

- No TLS support.
- No persistent connection lifecycle; server reads a single request per connection in the current loop.
- Transfer-Encoding: chunked is not implemented.
- Some HTTP methods are intentionally not supported and will return ErrUnsupportedMethod.
- Some HTTP status codes are still missing. 
- No advanced header validation or full RFC compliance.
- Writer requires Content-Length for non-chunked writes; write overflow returns ErrWriteOverflow.
- No corresponding built-in client. 

## Running tests

Run unit tests with:

```
go test ./...
```

## Building

Build as a command (example):

```
go build ./...
```

or build a binary that uses this package (see example in main.go).

## License

This project is MIT licensed — see LICENSE.

Contributions, issues and suggestions are welcome.