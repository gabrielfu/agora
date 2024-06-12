// Manages the current selected request
package states

import "github.com/gabrielfu/tipi/internal"

type RequestContext struct {
	req *internal.Request
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
