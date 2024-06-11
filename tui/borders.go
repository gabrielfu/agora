package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type GenerateBorderOption struct {
	Title  string
	Footer string
}

func generateBorder(border lipgloss.Border, opt GenerateBorderOption, width int) lipgloss.Border {
	if opt.Title != "" {
		border.Top = border.Top + opt.Title + strings.Repeat(border.Top, width)
	}
	if opt.Footer != "" {
		repeatCount := width - len(opt.Footer) - 1
		if repeatCount < 0 {
			repeatCount = 0
		}
		border.Bottom = strings.Repeat(border.Bottom, repeatCount) + opt.Footer + border.Bottom
	}
	return border
}
