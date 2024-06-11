package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
)

type View uint

const (
	// Pane views
	SidebarView View = iota
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

	sidebar      *SidebarModel
	urlPane      *UrlPaneModel
	requestPane  *RequestPaneModel
	responsePane *ResponsePaneModel
	navigation   *NagivationModel

	focus View

	width        int
	height       int
	sidebarWidth float32
}

type Options func(*RootModel)

func WithSidebarWidth(width float32) Options {
	return func(m *RootModel) {
		m.sidebarWidth = width
	}
}

func NewRootModel(db *internal.RequestDatabase, opts ...Options) *RootModel {
	m := &RootModel{
		db:           db,
		sidebar:      &SidebarModel{},
		urlPane:      &UrlPaneModel{},
		requestPane:  &RequestPaneModel{},
		responsePane: &ResponsePaneModel{},
		navigation:   &NagivationModel{},
		focus:        SidebarView,
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
	m.sidebar.SetBorderColor("#ffffff")
	m.urlPane.SetBorderColor("#ffffff")
	m.requestPane.SetBorderColor("#ffffff")
	m.responsePane.SetBorderColor("#ffffff")

	switch v {
	case SidebarView:
		m.sidebar.SetBorderColor("#98C379")
	case UrlPaneView:
		m.urlPane.SetBorderColor("#98C379")
	case RequestPaneView:
		m.requestPane.SetBorderColor("#98C379")
	case ResponsePaneView:
		m.responsePane.SetBorderColor("#98C379")
	}
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		if isPaneView(m.focus) {
			switch msg.String() {
			case "1":
				m.setFocus(SidebarView)
			case "2":
				m.setFocus(UrlPaneView)
			case "3":
				m.setFocus(RequestPaneView)
			case "4":
				m.setFocus(ResponsePaneView)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height - 3

		sidebarWidth := int(float32(m.width) * m.sidebarWidth)
		m.sidebar.SetWidth(sidebarWidth)
		m.sidebar.SetHeight(m.height)

		urlPaneWidth := m.width - sidebarWidth - 2
		m.urlPane.SetWidth(urlPaneWidth)

		requestPaneWidth := urlPaneWidth / 2
		m.requestPane.SetWidth(requestPaneWidth)
		m.requestPane.SetHeight(m.height - 3)

		responsePaneWidth := urlPaneWidth - requestPaneWidth - 2
		m.responsePane.SetWidth(responsePaneWidth)
		m.responsePane.SetHeight(m.height - 3)
	}
	return m, nil
}

func (m RootModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.sidebar.View(),
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
