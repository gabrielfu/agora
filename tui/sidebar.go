package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SidebarModel struct {
	width  int
	height int
}

func (m *SidebarModel) SetWidth(width int) {
	m.width = width
}

func (m *SidebarModel) SetHeight(height int) {
	m.height = height
}

func (m SidebarModel) Init() tea.Cmd {
	return nil
}

func (m SidebarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m SidebarModel) View() string {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: "Sidebar"},
		m.width,
	)
	style := lipgloss.NewStyle().
		BorderStyle(border).
		Width(m.width).
		Height(m.height)
	return style.Render()
}
