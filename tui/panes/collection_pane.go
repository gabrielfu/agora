package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CollectionPaneModel struct {
	width       int
	height      int
	borderColor string
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

func (m CollectionPaneModel) Update(msg tea.Msg) (CollectionPaneModel, tea.Cmd) {
	return m, nil
}

func (m CollectionPaneModel) View() string {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: []string{"[1]", "Collection"}},
		m.width,
	)
	style := lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
	return style.Render()
}
