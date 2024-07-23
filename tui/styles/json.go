package styles

import (
	"encoding/json"

	"github.com/tidwall/pretty"
)

func IsValidJson(s string) bool {
	var m any
	return json.Unmarshal([]byte(s), &m) == nil
}

func PrettifyJson(s string) string {
	s = string(pretty.Pretty([]byte(s)))
	return s
}

func PrettifyJsonIfValid(s string) string {
	if IsValidJson(s) {
		return PrettifyJson(s)
	}
	return s
}

func ColorizeJson(s string) string {
	return string(pretty.Color([]byte(s), nil))
}

func ColorizeJsonIfValid(s string) string {
	if IsValidJson(s) {
		return ColorizeJson(s)
	}
	return s
}

func MinifyJson(s string) string {
	return string(pretty.Ugly([]byte(s)))
}

func MinifyJsonBytes(b []byte) []byte {
  return pretty.Ugly(b)
}
