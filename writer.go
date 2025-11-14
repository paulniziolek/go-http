package main

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

