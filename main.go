package main

import "strconv"

type helloWorldHandler struct {
	greeting string
}

func NewHelloWorldHandler() *helloWorldHandler {
	return &helloWorldHandler{
		greeting: "hello there",
	}
}

func (h *helloWorldHandler) ServeHTTP(w ResponseWriter, req *Request) {
	if req.Method == "GET" {
		w.Header().Set("Content-Length", strconv.Itoa(len(h.greeting)))
		w.Write([]byte(h.greeting))
	} else {
		w.WriteHeader(StatusMethodNotAllowed)
	}
}

func main() {
	server := NewServer(":42069")
	helloHandler := NewHelloWorldHandler()
	server.Handle("/", helloHandler)
	server.ListenAndServe()

}
