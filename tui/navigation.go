package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/tui/views"
)

type NagivationModel struct {
	content string
	focus   views.View
}

func (m *NagivationModel) SetContent(content string) {
	m.content = content
}

func (m *NagivationModel) SetFocus(focus views.View) {
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
	case views.CollectionPaneView:
		keymap = CollectionPaneKeymap
	case views.CollectionListPaneView:
		keymap = CollectionListPaneKeymap
	case views.UrlPaneView:
		keymap = UrlPaneKeymap
	case views.RequestPaneView:
		keymap = RequestPaneKeymap
	case views.ResponsePaneView:
		keymap = ResponsePaneKeymap
	case views.SelectMethodDialogView:
		keymap = SelectMethodDialogKeymap
	case views.TextInputDialogView:
		keymap = TextInputDialogKeymap
	case views.TextAreaDialogView:
		keymap = TextAreaDialogKeymap
	}
	m.content = m.renderKeymap(keymap)
}

func (m NagivationModel) Update(msg tea.Msg) (NagivationModel, tea.Cmd) {
	return m, nil
}

func (m NagivationModel) View() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#61AFEF")).
		Render(m.content)
}
