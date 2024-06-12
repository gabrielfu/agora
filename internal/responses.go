package internal

type Response struct {
	id         string
	StatusCode int
	Content    []byte
	Headers    map[string]string
}

func NewResponse(statusCode int, content []byte, headers map[string]string) *Response {
	return &Response{
		id:         randomID(),
		StatusCode: statusCode,
		Content:    content,
		Headers:    headers,
	}
}

func (r Response) String() string {
	return string(r.Content)
}

func (r Response) ID() string {
	return r.id
}

func (r Response) ContentType() string {
	v, ok := r.Headers["Content-Type"]
	if !ok {
		return ""
	}
	return v
}
