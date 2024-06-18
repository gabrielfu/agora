package internal

type Response struct {
	id         string
	StatusCode int
	Content    []byte
	Headers    KVPairs
}

func NewResponse(statusCode int, content []byte, headers KVPairs) *Response {
	return &Response{
		id:         RandomID(),
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
	for _, kv := range r.Headers {
		if kv.Key == "Content-Type" {
			return kv.Value
		}
	}
	return ""
}
