package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui/dialogs"
	"github.com/gabrielfu/agora/tui/messages"
	"github.com/gabrielfu/agora/tui/panes"
	"github.com/gabrielfu/agora/tui/states"
	"github.com/gabrielfu/agora/tui/styles"
	"github.com/gabrielfu/agora/tui/views"
)

// RootModel implements tea.RootModel interface
type RootModel struct {
	store *internal.RequestFileStore

	collectionPane panes.CollectionPaneModel
	urlPane        panes.UrlPaneModel
	requestPane    panes.RequestPaneModel
	responsePane   panes.ResponsePaneModel
	navigation     NagivationModel

	focus views.View
	rctx  *states.RequestContext
	dctx  *states.DialogContext

	width               int
	height              int
	collectionPaneWidth float32
	enoughSpace         bool
}

type Options func(*RootModel)

func WithCollectionPaneWidth(width float32) Options {
	return func(m *RootModel) {
		m.collectionPaneWidth = width
	}
}

func NewRootModel(store *internal.RequestFileStore, opts ...Options) *RootModel {
	rctx := states.NewRequestContext()
	dctx := states.NewDialogContext()
	m := &RootModel{
		store:          store,
		collectionPane: panes.NewCollectionPaneModel(rctx, dctx),
		urlPane:        panes.NewUrlPaneModel(rctx, dctx),
		requestPane:    panes.NewRequestPaneModel(rctx, dctx),
		responsePane:   panes.NewResponsePaneModel(rctx),
		navigation:     NagivationModel{},
		focus:          views.CollectionPaneView,
		rctx:           rctx,
		dctx:           dctx,
		enoughSpace:    true,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		messages.SetFocusCmd(views.CollectionPaneView),
	)
}

func (m *RootModel) setFocus(v views.View) {
	m.focus = v
	m.collectionPane.SetBorderColor(styles.DefaultBorderColor)
	m.urlPane.SetBorderColor(styles.DefaultBorderColor)
	m.requestPane.SetBorderColor(styles.DefaultBorderColor)
	m.responsePane.SetBorderColor(styles.DefaultBorderColor)
	m.collectionPane.Blur()

	switch v {
	case views.CollectionPaneView:
		m.collectionPane.SetBorderColor(styles.FocusBorderColor)
		m.collectionPane.Focus()
	case views.UrlPaneView:
		m.urlPane.SetBorderColor(styles.FocusBorderColor)
	case views.RequestPaneView:
		m.requestPane.SetBorderColor(styles.FocusBorderColor)
	case views.ResponsePaneView:
		m.responsePane.SetBorderColor(styles.FocusBorderColor)
	}

	m.navigation.SetFocus(v)
}

func (m *RootModel) updatePanes(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.collectionPane, cmd = m.collectionPane.Update(msg)
	cmds = append(cmds, cmd)
	m.urlPane, cmd = m.urlPane.Update(msg)
	cmds = append(cmds, cmd)
	m.requestPane, cmd = m.requestPane.Update(msg)
	cmds = append(cmds, cmd)
	m.responsePane, cmd = m.responsePane.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *RootModel) updateDialogFocus() {
	if !m.dctx.Empty() {
		switch m.dctx.Dialog().(type) {
		case *dialogs.SelectMethodDialog:
			m.setFocus(views.SelectMethodDialogView)
		case *dialogs.TextInputDialog, *dialogs.DoubleTextInputDialog:
			m.setFocus(views.TextInputDialogView)
		case *dialogs.TextAreaDialog:
			m.setFocus(views.TextAreaDialogView)
		}
	}
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case messages.SetFocusMsg:
		m.setFocus(msg.View)
	case messages.ExitDialogMsg:
		m.dctx.Clear()
		m.setFocus(msg.Dest)
	case messages.ExecuteRequestMsg:
		m.rctx.Exec()
	case messages.UpdateRequestMsg:
		req := m.rctx.Request().Copy()
		msg.Func(&req)
		m.store.UpdateRequest(req)
	case messages.CreateRequestMsg:
		m.store.CreateRequest(msg.Req)
	case messages.DeleteRequestMsg:
		m.store.DeleteRequest(msg.ID)
		m.rctx.Clear()
	case messages.CopyRequestMsg:
		newReq := msg.Req.CopyWithNewID()
		m.store.CreateRequest(newReq)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.dctx.Empty() {
				return m, tea.Quit
			}
		}
		if !m.enoughSpace {
			break
		}
		// switch focus
		if views.IsPaneView(m.focus) {
			switch msg.String() {
			case "1":
				m.setFocus(views.CollectionPaneView)
			case "2":
				m.setFocus(views.UrlPaneView)
			case "3":
				m.setFocus(views.RequestPaneView)
			case "4":
				m.setFocus(views.ResponsePaneView)
			}
		}
		if !m.dctx.Empty() {
			dialog, cmd := m.dctx.Dialog().Update(msg)
			m.dctx.SetDialog(dialog.(dialogs.Dialog))
			cmds = append(cmds, cmd)
		}
		// update focused pane
		switch m.focus {
		case views.CollectionPaneView:
			m.collectionPane, cmd = m.collectionPane.Update(msg)
		case views.UrlPaneView:
			m.urlPane, cmd = m.urlPane.Update(msg)
		case views.RequestPaneView:
			m.requestPane, cmd = m.requestPane.Update(msg)
		case views.ResponsePaneView:
			m.responsePane, cmd = m.responsePane.Update(msg)
		}
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height - 3

		m.enoughSpace = false
		if m.width >= 78 && m.height >= 8 {
			m.enoughSpace = true
		}

		collectionPaneWidth := int(float32(m.width) * m.collectionPaneWidth)
		m.collectionPane.SetWidth(collectionPaneWidth)
		m.collectionPane.SetHeight(m.height)

		urlPaneWidth := m.width - collectionPaneWidth - 2
		m.urlPane.SetWidth(urlPaneWidth)

		requestPaneWidth := urlPaneWidth / 2
		m.requestPane.SetWidth(requestPaneWidth)
		m.requestPane.SetHeight(m.height - 3)

		responsePaneWidth := urlPaneWidth - requestPaneWidth - 2
		m.responsePane.SetWidth(responsePaneWidth)
		m.responsePane.SetHeight(m.height - 3)

		m.dctx.SetDialogWidth(m.width)
		m.dctx.SetDialogHeight(m.height)

		cmd = m.updatePanes(msg)
		cmds = append(cmds, cmd)
	}

	// Set requests for collection pane
	reqs, err := m.store.ListRequests()
	if err != nil {
		return m, tea.Quit
	}
	m.collectionPane.SetRequests(reqs)
	m.requestPane.Refresh()
	m.responsePane.Refresh()
	m.updateDialogFocus()

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	if !m.enoughSpace {
		return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(
			lipgloss.Place(
				m.width, m.height,
				lipgloss.Center, lipgloss.Center,
				"Not enough space to render panels",
			),
		)
	}
	var content string
	if !m.dctx.Empty() {
		content = lipgloss.Place(
			m.width+2, m.height+3,
			lipgloss.Center, lipgloss.Center,
			m.dctx.View(),
		)
	} else {
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.collectionPane.View(),
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.urlPane.View(),
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					m.requestPane.View(),
					m.responsePane.View(),
				),
			),
		)
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		m.navigation.View(),
	)
}
