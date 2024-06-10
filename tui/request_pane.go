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
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Width(m.width).
		Height(m.height)
	return style.Render("Request pane")
}
