package panes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui/messages"
	"github.com/gabrielfu/agora/tui/states"
	"github.com/gabrielfu/agora/tui/styles"
	"github.com/gabrielfu/agora/tui/views"
	"github.com/mattn/go-runewidth"
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
	fingerprint string

	tab      responsePaneTab
	viewport viewport.Model
	table    table.Model
}

func NewResponsePaneModel(rctx *states.RequestContext) ResponsePaneModel {
	t := table.New(
		table.WithColumns(makeKeyValueColumns(0)),
		table.WithRows(make([]table.Row, 0)),
		table.WithFocused(true),
		table.WithStyles(tableStyles()),
	)
	t.KeyMap.HalfPageUp.SetEnabled(false)
	t.KeyMap.HalfPageDown.SetEnabled(false)
	return ResponsePaneModel{
		rctx:     rctx,
		tab:      responseBodyTab,
		table:    t,
		viewport: viewport.New(0, 0),
	}
}

func (m *ResponsePaneModel) SetWidth(width int) {
	m.width = width
	m.table.SetWidth(width)
	m.table.SetColumns(makeKeyValueColumns(width))
	m.viewport.Width = width - 2
}

func (m *ResponsePaneModel) SetHeight(height int) {
	m.height = height
	m.table.SetHeight(height - 4)
	verticalMarginHeight := 10
	m.viewport.Height = height - verticalMarginHeight
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

func (m ResponsePaneModel) renderStatus() string {
	if m.rctx.Empty() {
		return ""
	}
	var text, color string
	if err := m.rctx.Error(); err != nil {
		text = "Error"
		color = styles.StatusErrorColor
	} else if m.rctx.Response() != nil {
		statusCode := m.rctx.Response().StatusCode
		text = fmt.Sprintf("%d %s", statusCode, internal.StatusText(statusCode))
		color = styles.StatusCodeColor(statusCode)
	} else {
		return ""
	}
	text = runewidth.Truncate(text, m.width-2, "…")
	return lipgloss.NewStyle().
		Width(m.width-2).
		MaxHeight(1).
		Foreground(lipgloss.Color(color)).
		Render(text) + "\n"
}

func (m ResponsePaneModel) renderDuration() string {
	if m.rctx.Response() == nil {
		return ""
	}
	text := m.rctx.Duration().String()
	return lipgloss.NewStyle().
		Width(m.width-2).
		MaxHeight(1).
		Foreground(lipgloss.Color(styles.StatusCode200Color)).
		Render(text) + "\n"
}

func (m ResponsePaneModel) generateStyle() lipgloss.Style {
	var footer []string
	if m.tab == responseBodyTab && m.viewport.TotalLineCount() > 0 {
		footer = append(footer, fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	} else if m.tab == responseHeadersTab {
		footer = append(footer, tableFooter(&m.table))
	}
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: []string{"[5]", "Response"}, Footer: footer},
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
	rows := make([]table.Row, 0)
	if m.rctx.Empty() {
		m.fingerprint = ""
		m.viewport.SetContent("")
		m.table.SetRows(rows)
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
				rows = append(rows, table.Row{kv.Key, kv.Value})
			}
		}
		m.table.SetRows(rows)
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
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ResponsePaneModel) View() string {
	var text string
	text = m.renderTabBar()
	text += m.renderStatus()
	text += m.renderDuration()
	switch m.tab {
	case responseHeadersTab:
		text += renderTableWithoutHeader(&m.table)
	case responseBodyTab:
		text += m.viewport.View()
	}
	return m.generateStyle().Render(text)
}
