package panes

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
)

type CollectionPaneModel struct {
	width       int
	height      int
	borderColor string

	requests []internal.Request
}

func (m *CollectionPaneModel) SetWidth(width int) {
	m.width = width
}

func (m *CollectionPaneModel) SetHeight(height int) {
	m.height = height
}

func (m *CollectionPaneModel) SetBorderColor(color string) {
	m.borderColor = color
}

func (m *CollectionPaneModel) SetRequests(requests []internal.Request) {
	m.requests = requests
}

func (m CollectionPaneModel) Update(msg tea.Msg) (CollectionPaneModel, tea.Cmd) {
	return m, nil
}

func (m CollectionPaneModel) generateStyle() lipgloss.Style {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: []string{"[1]", "Collection"}},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
}

func (m CollectionPaneModel) View() string {
	var text string
	for _, request := range m.requests {
		method := RenderMethod(request.Method)
		text += fmt.Sprintf("%s %s\n", method, request.URL)
	}
	return m.generateStyle().Render(text)
}
