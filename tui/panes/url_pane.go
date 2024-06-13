package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
	"github.com/gabrielfu/tipi/tui/dialogs"
	"github.com/gabrielfu/tipi/tui/messages"
	"github.com/gabrielfu/tipi/tui/states"
	"github.com/gabrielfu/tipi/tui/styles"
	"github.com/gabrielfu/tipi/tui/views"
)

func updateUrlCmd(url string) tea.Cmd {
	return func() tea.Msg {
		return messages.UpdateRequestMsg{
			Func: func(r *internal.Request) {
				r.URL = url
			},
		}
	}
}

type UrlPaneModel struct {
	width       int
	height      int
	borderColor string

	rctx               *states.RequestContext
	dctx               *states.DialogContext
	selectMethodDialog dialogs.SelectMethodDialog
	textInputDialog    dialogs.TextInputDialog
}

func NewUrlPaneModel(rctx *states.RequestContext, dctx *states.DialogContext) UrlPaneModel {
	return UrlPaneModel{
		rctx:               rctx,
		dctx:               dctx,
		selectMethodDialog: dialogs.NewSelectMethodDialog(),
		textInputDialog: dialogs.NewTextInputDialog(
			64,
			[]string{"URL"},
			nil,
			updateUrlCmd,
			views.UrlPaneView,
		),
	}
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

func (m UrlPaneModel) generateStyle() lipgloss.Style {
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: []string{"[2]", "URL"}},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
}

func (m UrlPaneModel) Update(msg tea.Msg) (UrlPaneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "m":
			if !m.rctx.Empty() {
				m.dctx.SetDialog(&m.selectMethodDialog)
			}
		case "u":
			if !m.rctx.Empty() {
				m.textInputDialog.SetValue(m.rctx.Request().URL)
				m.textInputDialog.Focus()
				m.dctx.SetDialog(&m.textInputDialog)
			}
		case "esc":
			return m, func() tea.Msg { return messages.SetFocusMsg{View: views.CollectionPaneView} }
		}
	}
	return m, nil
}

func (m UrlPaneModel) View() string {
	var text string
	if !m.rctx.Empty() {
		request := m.rctx.Request()
		method := styles.RenderMethodWithColor(request.Method)
		u := request.URL
		text = method + " " + u
	}
	return m.generateStyle().Render(text)
}
