package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RequestPaneModel struct {
	width  int
	height int
}

func (m *RequestPaneModel) SetWidth(width int) {
	m.width = width
}

func (m *RequestPaneModel) SetHeight(height int) {
	m.height = height
}

func (m RequestPaneModel) Init() tea.Cmd {
	return nil
}

func (m RequestPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m RequestPaneModel) View() string {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: []string{"[3]", "Request"}},
		m.width,
	)
	style := lipgloss.NewStyle().
		BorderStyle(border).
		Width(m.width).
		Height(m.height)
	return style.Render()
}
