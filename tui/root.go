package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
	"github.com/gabrielfu/tipi/tui/panes"
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

func isModalView(v View) bool {
	return !isPaneView(v)
}

// RootModel implements tea.RootModel interface
type RootModel struct {
	db *internal.RequestDatabase

	collectionPane panes.CollectionPaneModel
	urlPane        panes.UrlPaneModel
	requestPane    panes.RequestPaneModel
	responsePane   panes.ResponsePaneModel
	navigation     NagivationModel

	focus View

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

func NewRootModel(db *internal.RequestDatabase, opts ...Options) *RootModel {
	m := &RootModel{
		db:             db,
		collectionPane: panes.CollectionPaneModel{},
		urlPane:        panes.UrlPaneModel{},
		requestPane:    panes.RequestPaneModel{},
		responsePane:   panes.ResponsePaneModel{},
		navigation:     NagivationModel{},
		focus:          CollectionPaneView,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m RootModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RootModel) setFocus(v View) {
	m.focus = v
	m.collectionPane.SetBorderColor("#ffffff")
	m.urlPane.SetBorderColor("#ffffff")
	m.requestPane.SetBorderColor("#ffffff")
	m.responsePane.SetBorderColor("#ffffff")

	switch v {
	case CollectionPaneView:
		m.collectionPane.SetBorderColor("#98C379")
	case UrlPaneView:
		m.urlPane.SetBorderColor("#98C379")
	case RequestPaneView:
		m.requestPane.SetBorderColor("#98C379")
	case ResponsePaneView:
		m.responsePane.SetBorderColor("#98C379")
	}
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
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
