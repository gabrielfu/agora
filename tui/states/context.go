// Manages the current selected request
package states

import (
	"fmt"
	"io"
	"strings"

	"github.com/gabrielfu/tipi/internal"
)

type RequestContext struct {
	req  *internal.Request
	resp *internal.Response
	err  error
}

func NewRequestContext() *RequestContext {
	return &RequestContext{}
}

func (c *RequestContext) Empty() bool {
	return c.req == nil
}

func (c *RequestContext) Request() *internal.Request {
	return c.req
}

func (c *RequestContext) SetRequest(req *internal.Request) {
	c.req = req
}

func (c *RequestContext) Response() *internal.Response {
	return c.resp
}

func (c *RequestContext) SetResponse(resp *internal.Response) {
	c.resp = resp
}

func (c *RequestContext) Error() error {
	return c.err
}

func (c *RequestContext) SetError(err error) {
	c.err = err
}

func (c *RequestContext) Clear() {
	c.req = nil
	c.resp = nil
	c.err = nil
}

func (c *RequestContext) Exec() {
	response, err := c.req.Exec()
	c.err = err
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		c.err = err
		return
	}

	headers := make(map[string]string)
	for k, v := range response.Header {
		headers[k] = strings.Join(v, ", ")
	}

	c.resp = internal.NewResponse(
		response.StatusCode,
		content,
		headers,
	)
}
