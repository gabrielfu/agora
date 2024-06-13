package dialogs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/messages"
	"github.com/gabrielfu/tipi/tui/styles"
	"github.com/gabrielfu/tipi/tui/views"
)

type SelectMethodDialog struct {
	width  int
	height int
}

func NewSelectMethodDialog() SelectMethodDialog {
	return SelectMethodDialog{}
}

func (m *SelectMethodDialog) SetWidth(width int) {
	m.width = width
}

func (m *SelectMethodDialog) SetHeight(height int) {
	m.height = height
}

func (m SelectMethodDialog) generateStyle() lipgloss.Style {
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: []string{"Method"}},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(styles.FocusBorderColor)).
		Width(m.width).
		Height(m.height).
		Padding(0, 1)
}

func (m SelectMethodDialog) Prev() views.View {
	return views.UrlPaneView
}

func (m SelectMethodDialog) Update(msg tea.Msg) (any, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd := func() tea.Msg {
				return messages.ExitDialogMsg{Dest: views.UrlPaneView}
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m SelectMethodDialog) View() string {
	text := "GET\nPOST\nPUT\nPATCH\nDELETE\nHEAD\nOPTIONS"
	return m.generateStyle().Render(text)
}
