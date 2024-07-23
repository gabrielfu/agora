// Manages the current selected request
package states

import (
	"io"
	"strings"
	"time"

	"github.com/gabrielfu/agora/internal"
)

type RequestContext struct {
	req         *internal.Request
	resp        *internal.Response
	err         error
	fingerprint string // not a real fingerprint, just a string to identify the state
	duration    time.Duration
}

func NewRequestContext() *RequestContext {
	return &RequestContext{}
}

func (c *RequestContext) Fingerprint() string {
	return c.fingerprint
}

func (c *RequestContext) newFingerprint() {
	c.fingerprint = internal.RandomID()
}

func (c *RequestContext) Empty() bool {
	return c.req == nil
}

func (c *RequestContext) Request() *internal.Request {
	return c.req
}

func (c *RequestContext) SetRequest(req *internal.Request) {
	c.req = req
	c.newFingerprint()
}

func (c *RequestContext) Response() *internal.Response {
	return c.resp
}

func (c *RequestContext) SetResponse(resp *internal.Response) {
	c.resp = resp
	c.newFingerprint()
}

func (c *RequestContext) Error() error {
	return c.err
}

func (c *RequestContext) SetError(err error) {
	c.err = err
	c.newFingerprint()
}

func (c *RequestContext) Duration() time.Duration {
	return c.duration
}

func (c *RequestContext) SetDuration(d time.Duration) {
	c.duration = d
	c.newFingerprint()
}

func (c *RequestContext) Clear() {
	c.req = nil
	c.resp = nil
	c.err = nil
	c.fingerprint = ""
	c.duration = 0
}

func (c *RequestContext) Exec() {
	start := time.Now()
	response, err := c.req.Exec()
	duration := time.Since(start)
	c.duration = duration
	defer c.newFingerprint()

	c.err = err
	if err != nil {
		return
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		c.err = err
		return
	}

	var headers internal.KVPairs = make([]internal.KVPair, 0)
	for k, v := range response.Header {
		headers = append(headers, internal.KVPair{
			Key:   k,
			Value: strings.Join(v, ", "),
		})
	}

	c.resp = internal.NewResponse(
		response.StatusCode,
		content,
		headers,
	)
}
