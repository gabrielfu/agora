package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/states"
)

type RequestPaneModel struct {
	width       int
	height      int
	borderColor string

	ctx *states.RequestContext
}

func NewRequestPaneModel(ctx *states.RequestContext) RequestPaneModel {
	return RequestPaneModel{ctx: ctx}
}

func (m *RequestPaneModel) SetWidth(width int) {
	m.width = width
}

func (m *RequestPaneModel) SetHeight(height int) {
	m.height = height
}

func (m *RequestPaneModel) SetBorderColor(color string) {
	m.borderColor = color
}

func (m RequestPaneModel) Update(msg tea.Msg) (RequestPaneModel, tea.Cmd) {
	return m, nil
}

func (m RequestPaneModel) generateStyle() lipgloss.Style {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: []string{"[3]", "Request"}},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
}

func (m RequestPaneModel) View() string {
	var text string
	if !m.ctx.Empty() {
		request := m.ctx.Request()
		if request != nil {
			text = request.String()
		}
	}
	return m.generateStyle().Render(text)
}
