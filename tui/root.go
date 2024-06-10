package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
)

// RootModel implements tea.RootModel interface
type RootModel struct {
	db *internal.RequestDatabase

	sidebar      *SidebarModel
	urlPane      *UrlPaneModel
	requestPane  *RequestPaneModel
	responsePane *ResponsePaneModel

	width  int
	height int

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
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m RootModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height - 2

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
	return lipgloss.JoinHorizontal(
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
	)
}
