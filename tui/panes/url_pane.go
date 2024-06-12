package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/states"
)

type UrlPaneModel struct {
	width       int
	height      int
	borderColor string

	ctx *states.RequestContext
}

func NewUrlPaneModel(ctx *states.RequestContext) UrlPaneModel {
	return UrlPaneModel{ctx: ctx}
}

func (m *UrlPaneModel) SetWidth(width int) {
	m.width = width
}

func (m *UrlPaneModel) SetHeight(height int) {
	m.height = height
}

func (m *UrlPaneModel) SetBorderColor(color string) {
	m.borderColor = color
}

func (m UrlPaneModel) Update(msg tea.Msg) (UrlPaneModel, tea.Cmd) {
	return m, nil
}

func (m UrlPaneModel) View() string {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: []string{"[2]", "URL"}},
		m.width,
	)
	style := lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
	var text string
	if !m.ctx.Empty() {
		text = m.ctx.Request().URL
	}
	return style.Render(text)
}
