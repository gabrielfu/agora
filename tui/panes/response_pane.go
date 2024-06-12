package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/states"
)

type ResponsePaneModel struct {
	width       int
	height      int
	borderColor string

	ctx *states.RequestContext
}

func NewResponsePaneModel(ctx *states.RequestContext) ResponsePaneModel {
	return ResponsePaneModel{ctx: ctx}
}

func (m *ResponsePaneModel) SetWidth(width int) {
	m.width = width
}

func (m *ResponsePaneModel) SetHeight(height int) {
	m.height = height
}

func (m *ResponsePaneModel) SetBorderColor(color string) {
	m.borderColor = color
}
func (m ResponsePaneModel) generateStyle() lipgloss.Style {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: []string{"[4]", "Response"}},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
}

func (m ResponsePaneModel) Update(msg tea.Msg) (ResponsePaneModel, tea.Cmd) {
	return m, nil
}

func (m ResponsePaneModel) View() string {
	var text string
	if !m.ctx.Empty() {
		err := m.ctx.Error()
		if err != nil {
			text = err.Error()
		} else {
			response := m.ctx.Response()
			if response != nil {
				text = response.String()
			}
		}
	}
	return m.generateStyle().Render(text)
}
