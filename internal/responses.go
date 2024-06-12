package internal

type Response struct {
	StatusCode int
	Content    []byte
	Headers    map[string]string
}

func (r Response) ContentType() string {
	v, ok := r.Headers["Content-Type"]
	if !ok {
		return ""
	}
	return v
}
