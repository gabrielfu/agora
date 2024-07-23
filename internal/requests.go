package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/gabrielfu/agora/tui/styles"
)

type KVPair struct {
	Key   string `json:"k" yaml:"key"`
	Value string `json:"v" yaml:"value"`
}

type KVPairs []KVPair

func (kvs KVPairs) Sort() {
	sort.Slice(kvs, func(i, j int) bool {
		if kvs[i].Key != kvs[j].Key {
			return kvs[i].Key < kvs[j].Key
		}
		return kvs[i].Value < kvs[j].Value
	})
}

func (kvs KVPairs) Add(key, value string) KVPairs {
	return append(kvs, KVPair{Key: key, Value: value})
}

// Remove removes the first occurrence of the key-value pair from the list.
func (kvs KVPairs) Remove(key, value string) KVPairs {
	var newKvs KVPairs
	removed := false
	for _, kv := range kvs {
		if !removed && (kv.Key == key && kv.Value == value) {
			removed = true
			continue
		}
		newKvs = append(newKvs, kv)
	}
	return newKvs
}

type Request struct {
	id     string `yaml:"id"`
	Name   string `yaml:"name"`
	Method string `yaml:"method"`
	URL    string `yaml:"url"`
	Body   []byte `yaml:"body"` // only supports json body for now
	// we use array of kv pair to preserve order
	Params  KVPairs `yaml:"params"`
	Headers KVPairs `yaml:"headers"`
	Auth    string  `yaml:"auth"`
}

// NewRequest creates a new request with a random id.
func NewRequest(method, url string) *Request {
	return &Request{
		id:     RandomID(),
		Method: method,
		URL:    url,
	}
}

func (r Request) Copy() Request {
	return Request{
		id:      r.id,
		Name:    r.Name,
		Method:  r.Method,
		URL:     r.URL,
		Body:    r.Body,
		Params:  r.Params,
		Headers: r.Headers,
		Auth:    r.Auth,
	}
}

func (r *Request) WithName(name string) *Request {
	r.Name = name
	return r
}

func (r *Request) WithBody(body []byte) *Request {
	body = styles.MinifyJsonBytes(body)
	r.Body = body
	return r
}

func (r *Request) WithParam(key, value string) *Request {
	r.Params = append(r.Params, KVPair{Key: key, Value: value})
	return r
}

func (r *Request) WithParams(params KVPairs) *Request {
	r.Params = params
	return r
}

func (r *Request) WithHeader(key, value string) *Request {
	r.Headers = append(r.Headers, KVPair{Key: key, Value: value})
	return r
}

func (r *Request) WithHeaders(headers KVPairs) *Request {
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

func (r *Request) RemoveParam(key, value string) {
	r.Params = r.Params.Remove(key, value)
}

func (r *Request) RemoveHeader(key, value string) {
	r.Headers = r.Headers.Remove(key, value)
}

func makeJsonBodyReader(body any) (io.Reader, error) {
	// todo: support other body dtype and content type
	switch body.(type) {
	case string, []byte:
		body = styles.MinifyJson(body.(string))
	}
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
	for _, kv := range r.Params {
		q.Add(kv.Key, kv.Value)
	}
	req.URL.RawQuery = q.Encode()

	for _, kv := range r.Headers {
		req.Header.Add(kv.Key, kv.Value)
	}
	client := &http.Client{}
	return client.Do(req)
}
