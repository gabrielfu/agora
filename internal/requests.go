package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type Request struct {
	id      string
	Name    string
	Method  string
	URL     string
	Body    any // only supports json body for now
	Params  map[string]string
	Headers map[string]string
	Auth    string
}

// NewRequest creates a new request with a random id.
func NewRequest(method, url string) *Request {
	return &Request{
		id:      randomID(),
		Method:  method,
		URL:     url,
		Params:  make(map[string]string),
		Headers: make(map[string]string),
	}
}

func randomID() string {
	return uuid.New().String()
}

func (r *Request) WithName(name string) *Request {
	r.Name = name
	return r
}

func (r *Request) WithBody(body any) *Request {
	r.Body = body
	return r
}

func (r *Request) WithParam(key, value string) *Request {
	r.Params[key] = value
	return r
}

func (r *Request) WithParams(params map[string]string) *Request {
	r.Params = params
	return r
}

func (r *Request) WithHeader(key, value string) *Request {
	r.Headers[key] = value
	return r
}

func (r *Request) WithHeaders(headers map[string]string) *Request {
	r.Headers = headers
	return r
}

func (r *Request) WithAuth(auth string) *Request {
	r.Auth = auth
	return r
}

func (r *Request) ID() string {
	return r.id
}

func (r Request) String() string {
	return fmt.Sprintf(
		"Request(ID=%s, Name=%s, Method=%s, URL=%s, Body=%v, Params=%v, Headers=%v, Auth=%s}",
		r.id, r.Name, r.Method, r.URL, r.Body, r.Params, r.Headers, r.Auth,
	)
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
