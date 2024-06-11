package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResponsePaneModel struct {
	width  int
	height int
}

func (m *ResponsePaneModel) SetWidth(width int) {
	m.width = width
}

func (m *ResponsePaneModel) SetHeight(height int) {
	m.height = height
}

func (m ResponsePaneModel) Init() tea.Cmd {
	return nil
}

func (m ResponsePaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ResponsePaneModel) View() string {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: "Response"},
		m.width,
	)
	style := lipgloss.NewStyle().
		BorderStyle(border).
		Width(m.width).
		Height(m.height)
	return style.Render()
}
