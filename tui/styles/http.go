package styles

func StatusCodeColor(code int) string {
	switch {
	case code >= 100 && code < 200:
		return StatusCode100Color
	case code >= 200 && code < 300:
		return StatusCode200Color
	case code >= 300 && code < 400:
		return StatusCode300Color
	case code >= 400 && code < 500:
		return StatusCode400Color
	case code >= 500 && code < 600:
		return StatusCode500Color
	default:
		return StatusCodeUnknownColor
	}
}
