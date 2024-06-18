package panes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/tui/messages"
	"github.com/gabrielfu/agora/tui/states"
	"github.com/gabrielfu/agora/tui/styles"
	"github.com/gabrielfu/agora/tui/views"
)

type responsePaneTab int

const (
	responseHeadersTab responsePaneTab = iota
	responseBodyTab
)

type ResponsePaneModel struct {
	width       int
	height      int
	borderColor string

	rctx        *states.RequestContext
	ready       bool
	fingerprint string
	viewport    viewport.Model

	tab  responsePaneTab
	list list.Model
}

func NewResponsePaneModel(rctx *states.RequestContext) ResponsePaneModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowFilter(false)
	return ResponsePaneModel{
		rctx: rctx,
		tab:  responseBodyTab,
		list: l,
	}
}

func (m *ResponsePaneModel) SetWidth(width int) {
	m.width = width
}

func (m *ResponsePaneModel) SetHeight(height int) {
	m.height = height
}

func (m *ResponsePaneModel) SetBorderColor(color string) {
	m.borderColor = color
}

func (m ResponsePaneModel) renderTabBar() string {
	tabs := []string{"Headers", "Body"}
	tabs[m.tab] = focusedStyle.Render(tabs[m.tab])
	separator := strings.Repeat(lipgloss.RoundedBorder().Bottom, m.width)
	separator = lipgloss.NewStyle().Foreground(lipgloss.Color(m.borderColor)).Render(separator)
	return strings.Join(tabs, " - ") + "\n" + separator + "\n"
}

func (m *ResponsePaneModel) switchTab(direction int) {
	m.tab = responsePaneTab((int(m.tab) + direction + 2) % 2)
}

func (m ResponsePaneModel) generateStyle() lipgloss.Style {
	var footer []string
	if m.tab == responseBodyTab && m.ready && m.viewport.TotalLineCount() > 0 {
		footer = append(footer, fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	}
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: []string{"[4]", "Response"}, Footer: footer},
		m.width,
	)
	// make the corner for the tab bar
	border.Left = border.Left + border.MiddleLeft + strings.Repeat(border.Left, m.height)
	border.Right = border.Right + border.MiddleRight + strings.Repeat(border.Right, m.height)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
}

func (m *ResponsePaneModel) Refresh() {
	items := make([]list.Item, 0)
	if m.rctx.Empty() {
		m.fingerprint = ""
		m.viewport.SetContent("")
		m.list.SetItems(items)
		return
	}

	if m.fingerprint != m.rctx.Fingerprint() {
		m.fingerprint = m.rctx.Fingerprint()

		var text string
		if err := m.rctx.Error(); err != nil {
			text = err.Error()
		} else if m.rctx.Response() != nil {
			text = m.rctx.Response().String()
		}
		text = styles.ColorizeJsonIfValid(text)
		text = lipgloss.NewStyle().Width(m.width - 2).Render(text)
		m.viewport.SetContent(text)

		if m.rctx.Response() != nil {
			for _, kv := range m.rctx.Response().Headers {
				items = append(items, kvItem{key: kv.Key, value: kv.Value})
			}
		}
		m.list.SetItems(items)
	}
}

func (m ResponsePaneModel) Update(msg tea.Msg) (ResponsePaneModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, messages.SetFocusCmd(views.CollectionPaneView)
		case "[", "shift+tab":
			m.switchTab(-1)
		case "]", "tab":
			m.switchTab(1)
		}
	case tea.WindowSizeMsg:
		verticalMarginHeight := 8
		if !m.ready {
			m.viewport = viewport.New(msg.Width-2, msg.Height-verticalMarginHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 2
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
		m.list.SetSize(m.width-2, m.height-3)
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ResponsePaneModel) View() string {
	var text string
	text = m.renderTabBar()
	switch m.tab {
	case responseHeadersTab:
		text += m.list.View()
	case responseBodyTab:
		if !m.ready {
			text += "Initializing..."
		} else {
			text += m.viewport.View()
		}
	}
	return m.generateStyle().Render(text)
}
