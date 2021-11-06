package stub

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/ch629/mockservice/pkg/domain"
)

func RequestFromHTTP(r http.Request) domain.Request {
	bs, _ := io.ReadAll(r.Body)
	r.Body.Close()
	return domain.Request{
		QueryParameters: r.URL.Query(),
		Path:            r.URL.Path,
		Headers:         r.Header,
		Method:          r.Method,
		Body:            bs,
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
