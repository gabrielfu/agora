package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui/dialogs"
	"github.com/gabrielfu/agora/tui/messages"
	"github.com/gabrielfu/agora/tui/states"
	"github.com/gabrielfu/agora/tui/styles"
	"github.com/gabrielfu/agora/tui/views"
)

func updateUrlCmd(url string) tea.Cmd {
	return messages.UpdateRequestCmd(func(r *internal.Request) {
		r.URL = url
	})
}

func updateNameCmd(name string) tea.Cmd {
	return messages.UpdateRequestCmd(func(r *internal.Request) {
		r.Name = name
	})
}

type UrlPaneModel struct {
	width       int
	height      int
	borderColor string

	rctx               *states.RequestContext
	dctx               *states.DialogContext
	selectMethodDialog dialogs.SelectMethodDialog
	editUrlDialog      dialogs.TextInputDialog
	editNameDialog     dialogs.TextInputDialog
}

func NewUrlPaneModel(rctx *states.RequestContext, dctx *states.DialogContext) UrlPaneModel {
	return UrlPaneModel{
		rctx:               rctx,
		dctx:               dctx,
		selectMethodDialog: dialogs.NewSelectMethodDialog(),
		editUrlDialog: dialogs.NewTextInputDialog(
			64,
			[]string{"URL"},
			nil,
			updateUrlCmd,
			views.UrlPaneView,
		),
		editNameDialog: dialogs.NewTextInputDialog(
			64,
			[]string{"Name"},
			nil,
			updateNameCmd,
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
	title := []string{"[2]", "URL"}
	if !m.rctx.Empty() && m.rctx.Request().Name != "" {
		title = append(title, "", "", "("+m.rctx.Request().Name+")")
	}
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: title},
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
		case "x":
			if !m.rctx.Empty() {
				return m, messages.ExecuteRequestCmd
			}
		case "m":
			if !m.rctx.Empty() {
				m.dctx.SetDialog(&m.selectMethodDialog)
			}
		case "enter":
			if !m.rctx.Empty() {
				m.editUrlDialog.SetValue(m.rctx.Request().URL)
				m.editUrlDialog.Focus()
				m.dctx.SetDialog(&m.editUrlDialog)
			}
		case "r":
			if !m.rctx.Empty() {
				m.editNameDialog.SetValue(m.rctx.Request().Name)
				m.editNameDialog.Focus()
				m.dctx.SetDialog(&m.editNameDialog)
			}
		case "esc":
			return m, messages.SetFocusCmd(views.CollectionPaneView)
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
