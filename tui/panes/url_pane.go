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

func (m UrlPaneModel) generateStyle() lipgloss.Style {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: []string{"[2]", "URL"}},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
}

func (m UrlPaneModel) View() string {
	var text string
	if !m.ctx.Empty() {
		request := m.ctx.Request()
		method := RenderMethodWithColor(request.Method)
		u := RenderURL(request.URL)
		text = method + " " + u
	}
	return m.generateStyle().Render(text)
}
