package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UrlPaneModel struct {
	width       int
	height      int
	borderColor string
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

func (m UrlPaneModel) Init() tea.Cmd {
	return nil
}

func (m UrlPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return style.Render()
}
