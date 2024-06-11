package panes

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func RenderURL(url string, width int) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	return ansi.Truncate(url, width, "...")
}
