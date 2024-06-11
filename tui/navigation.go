package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NagivationModel struct {
	content string
}

func (m *NagivationModel) SetContent(content string) {
	m.content = content
}

func (m NagivationModel) Init() tea.Cmd {
	return nil
}

func (m NagivationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			m.content = "Up"
		case "j":
			m.content = "Down"
		default:
			m.content = "Unknown"
		}
	}
	return m, nil
}

func (m NagivationModel) View() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#61AFEF")).
		Render(m.content)
}
