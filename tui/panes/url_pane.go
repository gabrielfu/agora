package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/dialogs"
	"github.com/gabrielfu/tipi/tui/states"
	"github.com/gabrielfu/tipi/tui/styles"
)

type UrlPaneModel struct {
	width       int
	height      int
	borderColor string

	rctx               *states.RequestContext
	dctx               *states.DialogContext
	selectMethodDialog dialogs.SelectMethodDialog
}

func NewUrlPaneModel(rctx *states.RequestContext, dctx *states.DialogContext) UrlPaneModel {
	return UrlPaneModel{rctx: rctx, dctx: dctx, selectMethodDialog: dialogs.NewSelectMethodDialog()}
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
			m.dctx.SetDialog(m.selectMethodDialog)
		}
	}
	return m, nil
}

func (m UrlPaneModel) View() string {
	var text string
	if !m.rctx.Empty() {
		request := m.rctx.Request()
		method := RenderMethodWithColor(request.Method)
		u := RenderURL(request.URL)
		text = method + " " + u
	}
	return m.generateStyle().Render(text)
}
