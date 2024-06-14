package panes

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/states"
	"github.com/gabrielfu/tipi/tui/styles"
)

type ResponsePaneModel struct {
	width       int
	height      int
	borderColor string

	rctx        *states.RequestContext
	ready       bool
	fingerprint string
	viewport    viewport.Model
}

func NewResponsePaneModel(rctx *states.RequestContext) ResponsePaneModel {
	return ResponsePaneModel{rctx: rctx}
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

func (m ResponsePaneModel) generateStyle() lipgloss.Style {
	var footer []string
	if m.ready && m.viewport.TotalLineCount() > 0 {
		footer = append(footer, fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	}
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: []string{"[4]", "Response"}, Footer: footer},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
}

func (m *ResponsePaneModel) Refresh() {
	if m.rctx.Empty() {
		m.fingerprint = ""
		m.viewport.SetContent("")
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
		m.viewport.SetContent(text)
	}
}

func (m ResponsePaneModel) Update(msg tea.Msg) (ResponsePaneModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		verticalMarginHeight := 8
		if !m.ready {
			m.viewport = viewport.New(msg.Width-2, msg.Height-verticalMarginHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 2
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ResponsePaneModel) View() string {
	var text string
	if !m.ready {
		text = "Initializing..."
	} else {
		text = m.viewport.View()
	}
	return m.generateStyle().Render(text)
}
