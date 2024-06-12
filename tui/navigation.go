package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NagivationModel struct {
	content string
	focus   View
}

func (m *NagivationModel) SetContent(content string) {
	m.content = content
}

func (m *NagivationModel) SetFocus(focus View) {
	m.focus = focus
	m.updateNagivationContent()
}

func (m *NagivationModel) renderKeymap(keymap *Keymap) string {
	var strs []string
	for el := keymap.Front(); el != nil; el = el.Next() {
		strs = append(strs, el.Value+": "+el.Key)
	}
	return strings.Join(strs, " | ")
}

func (m *NagivationModel) updateNagivationContent() {
	var keymap *Keymap = EmptyKeymap
	switch m.focus {
	case CollectionPaneView:
		keymap = CollectionPaneKeymap
	case UrlPaneView:
	case RequestPaneView:
	case ResponsePaneView:
	}
	m.content = m.renderKeymap(keymap)
}

func (m NagivationModel) Update(msg tea.Msg) (NagivationModel, tea.Cmd) {
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
