package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type GenerateBorderOption struct {
	Title  []string
	Footer []string
}

func GenerateBorder(border lipgloss.Border, opt GenerateBorderOption, width int) lipgloss.Border {
	if len(opt.Title) > 0 {
		title := strings.Join(opt.Title, border.Top)
		border.Top = border.Top + title + strings.Repeat(border.Top, width)
	}
	if len(opt.Footer) > 0 {
		footer := strings.Join(opt.Footer, border.Bottom)
		repeatCount := width - len(footer) - 1
		if repeatCount < 0 {
			repeatCount = 0
		}
		border.Bottom = strings.Repeat(border.Bottom, repeatCount) + footer + border.Bottom
	}
	return border
}
