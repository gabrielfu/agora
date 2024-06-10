package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	id      string
	Method  string
	URL     string
	Body    any // only supports json body for now
	Params  map[string]string
	Headers map[string]string
	Auth    string
}

func (r *Request) ID() string {
	return r.id
}

func makeJsonBodyReader(body any) (io.Reader, error) {
	marshalled, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("body is not a valid json: %w", err)
	}
	return bytes.NewReader(marshalled), nil
}

// Exec sends the request and returns the response.
// If error is not nil, the request is considered failed.
// Non-2xx status codes does not cause an error.
func (r *Request) Exec() (*http.Response, error) {
	body, err := makeJsonBodyReader(r.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range r.Params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	return client.Do(req)
}
