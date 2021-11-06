package stub

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// TODO: Cookies
// TODO: Form data
// TODO: Auth
type Request struct {
	Body            io.ReadCloser
	QueryParameters map[string][]string
	Headers         map[string][]string
	Path            string
	Method          string
}

func RequestFromHTTP(r http.Request) Request {
	return Request{
		QueryParameters: r.URL.Query(),
		Path:            r.URL.Path,
		Headers:         r.Header,
		Method:          r.Method,
		Body:            r.Body,
	}
}

type Response struct {
	Headers map[string]string
	Body    []byte
	Status  int
}

func (r Response) WriteTo(rw http.ResponseWriter) error {
	rw.WriteHeader(r.Status)
	for key, value := range r.Headers {
		rw.Header().Add(key, value)
	}
	if _, err := io.Copy(rw, bytes.NewReader(r.Body)); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	return nil
}
