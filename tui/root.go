package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
	"github.com/gabrielfu/tipi/tui/panes"
	"github.com/gabrielfu/tipi/tui/states"
)

type View uint

const (
	// Pane views
	CollectionPaneView View = iota
	UrlPaneView
	RequestPaneView
	ResponsePaneView
	// Modal views
	// ...
)

func isPaneView(v View) bool {
	return v <= ResponsePaneView
}

const (
	DefaultColor     = "#DCDFE4"
	FocusBorderColor = "#98C379"
)

// RootModel implements tea.RootModel interface
type RootModel struct {
	db *internal.RequestDatabase

	collectionPane panes.CollectionPaneModel
	urlPane        panes.UrlPaneModel
	requestPane    panes.RequestPaneModel
	responsePane   panes.ResponsePaneModel
	navigation     NagivationModel

	focus View
	ctx   *states.RequestContext

	width               int
	height              int
	collectionPaneWidth float32
}

type Options func(*RootModel)

func WithCollectionPaneWidth(width float32) Options {
	return func(m *RootModel) {
		m.collectionPaneWidth = width
	}
}

type SetFocusMsg struct {
	View View
}

func NewRootModel(db *internal.RequestDatabase, opts ...Options) *RootModel {
	ctx := states.NewRequestContext()
	m := &RootModel{
		db:             db,
		collectionPane: panes.NewCollectionPaneModel(ctx),
		urlPane:        panes.NewUrlPaneModel(ctx),
		requestPane:    panes.NewRequestPaneModel(ctx),
		responsePane:   panes.NewResponsePaneModel(ctx),
		navigation:     NagivationModel{},
		focus:          CollectionPaneView,
		ctx:            ctx,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		func() tea.Msg {
			return SetFocusMsg{View: CollectionPaneView}
		},
	)
}

func (m *RootModel) setFocus(v View) {
	m.focus = v
	m.collectionPane.SetBorderColor(DefaultColor)
	m.urlPane.SetBorderColor(DefaultColor)
	m.requestPane.SetBorderColor(DefaultColor)
	m.responsePane.SetBorderColor(DefaultColor)

	switch v {
	case CollectionPaneView:
		m.collectionPane.SetBorderColor(FocusBorderColor)
	case UrlPaneView:
		m.urlPane.SetBorderColor(FocusBorderColor)
	case RequestPaneView:
		m.requestPane.SetBorderColor(FocusBorderColor)
	case ResponsePaneView:
		m.responsePane.SetBorderColor(FocusBorderColor)
	}

	m.navigation.SetFocus(v)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case SetFocusMsg:
		m.setFocus(msg.View)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		if isPaneView(m.focus) {
			switch msg.String() {
			case "1":
				m.setFocus(CollectionPaneView)
			case "2":
				m.setFocus(UrlPaneView)
			case "3":
				m.setFocus(RequestPaneView)
			case "4":
				m.setFocus(ResponsePaneView)
			}
		}
		switch m.focus {
		case CollectionPaneView:
			m.collectionPane, cmd = m.collectionPane.Update(msg)
		case UrlPaneView:
			m.urlPane, cmd = m.urlPane.Update(msg)
		case RequestPaneView:
			m.requestPane, cmd = m.requestPane.Update(msg)
		case ResponsePaneView:
			m.responsePane, cmd = m.responsePane.Update(msg)
		}
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height - 3

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
	}

	// Set requests for collection pane
	reqs, err := m.db.ListRequests()
	if err != nil {
		return m, tea.Quit
	}
	m.collectionPane.SetRequests(reqs)

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
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
		),
		m.navigation.View(),
	)
}
