package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UrlPaneModel struct {
	width  int
	height int
}

func (m *UrlPaneModel) SetWidth(width int) {
	m.width = width
}

func (m *UrlPaneModel) SetHeight(height int) {
	m.height = height
}

func (m UrlPaneModel) Init() tea.Cmd {
	return nil
}

func (m UrlPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m UrlPaneModel) View() string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Width(m.width).
		Height(m.height)
	return style.Render("URL pane")
}
