package panes

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var methodColors = map[string]string{
	"GET":     "#68D696",
	"POST":    "#EED577",
	"PUT":     "#74AEF6",
	"PATCH":   "#C0A8E1",
	"DELETE":  "#EF968A",
	"HEAD":    "#68D696",
	"OPTIONS": "#E55AA8",
}

var methodShort = map[string]string{
	"GET":     "GET  ",
	"POST":    "POST ",
	"PUT":     "PUT  ",
	"PATCH":   "PATCH",
	"DELETE":  "DEL  ",
	"HEAD":    "HEAD ",
	"OPTIONS": "OPT  ",
}

func getMethodColor(method string) string {
	if color, ok := methodColors[method]; ok {
		return color
	}
	return methodColors["GET"]
}

func getMethodShort(method string) string {
	if short, ok := methodShort[method]; ok {
		return short
	}
	short := method[:min(3, len(method))]
	return short + strings.Repeat(" ", 5-len(short))
}

func RenderMethod(method string) string {
	color := getMethodColor(method)
	short := getMethodShort(method)
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(short)
}
