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
	// Example usage with hello world handler
	helloHandler := NewHelloWorldHandler()
	// Example usage with shorthand handler function
	server.HandleFunc("/hi", func(w ResponseWriter, r *Request) {
		response := "hi!"
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		w.Write([]byte(response))
	})
	server.Handle("/", helloHandler)
	server.ListenAndServe()

}
