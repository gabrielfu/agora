package panes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui/dialogs"
	"github.com/gabrielfu/agora/tui/messages"
	"github.com/gabrielfu/agora/tui/states"
	"github.com/gabrielfu/agora/tui/styles"
	"github.com/gabrielfu/agora/tui/views"
)

var focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.FocusBorderColor))

type requestPaneTab int

const (
	requestParamsTab requestPaneTab = iota
	requestHeadersTab
	requestBodyTab
)

func updateParamCmdFunc(cursor int, key string) dialogs.TextInputCmdFunc {
	return func(value string) tea.Cmd {
		return messages.UpdateRequestCmd(func(r *internal.Request) {
			r.UpdateParam(cursor, key, value)
		})
	}
}

var newParamCmdFunc dialogs.DoubleTextInputCmdFunc = func(key, value string) tea.Cmd {
	return messages.UpdateRequestCmd(func(r *internal.Request) {
		r.WithParam(key, value)
	})
}

func updateHeaderCmdFunc(cursor int, key string) dialogs.TextInputCmdFunc {
	return func(value string) tea.Cmd {
		return messages.UpdateRequestCmd(func(r *internal.Request) {
			r.UpdateHeader(cursor, key, value)
		})
	}
}

var updateBodyCmdFunc dialogs.TextAreaCmdFunc = func(body string) tea.Cmd {
	return messages.UpdateRequestCmd(func(r *internal.Request) {
		r.Body = []byte(body)
	})
}

var newHeaderCmdFunc dialogs.DoubleTextInputCmdFunc = func(key, value string) tea.Cmd {
	return messages.UpdateRequestCmd(func(r *internal.Request) {
		r.WithHeader(key, value)
	})
}

type RequestPaneModel struct {
	width       int
	height      int
	borderColor string

	rctx *states.RequestContext
	dctx *states.DialogContext

	tab                   requestPaneTab
	textInputDialog       dialogs.TextInputDialog
	doubleTextInputDialog dialogs.DoubleTextInputDialog
	textAreaDialog        dialogs.TextAreaDialog
	viewport              viewport.Model
	table                 table.Model
}

func NewRequestPaneModel(rctx *states.RequestContext, dctx *states.DialogContext) RequestPaneModel {
	t := table.New(
		table.WithColumns(makeKeyValueColumns(0)),
		table.WithRows(make([]table.Row, 0)),
		table.WithFocused(true),
		table.WithStyles(tableStyles()),
	)
	t.KeyMap.HalfPageUp.SetEnabled(false)
	t.KeyMap.HalfPageDown.SetEnabled(false)
	return RequestPaneModel{
		rctx: rctx,
		dctx: dctx,
		tab:  requestParamsTab,
		textInputDialog: dialogs.NewTextInputDialog(
			64,
			[]string{"Param"},
			nil,
			nil,
			views.RequestPaneView,
		),
		doubleTextInputDialog: dialogs.NewDoubleTextInputDialog(
			64,
			[]string{"Key"},
			nil,
			[]string{"Value"},
			nil,
			nil,
			views.RequestPaneView,
		),
		textAreaDialog: dialogs.NewTextAreaDialog(
			64,
			7,
			[]string{"Body"},
			nil,
			nil,
			views.RequestPaneView,
		),
		table:    t,
		viewport: viewport.New(0, 0),
	}
}

func (m *RequestPaneModel) SetWidth(width int) {
	m.width = width
	m.table.SetWidth(width)
	m.table.SetColumns(makeKeyValueColumns(width))
}

func (m *RequestPaneModel) SetHeight(height int) {
	m.height = height
	m.table.SetHeight(height - 2)
}

func (m *RequestPaneModel) SetBorderColor(color string) {
	m.borderColor = color
}

func (m RequestPaneModel) renderTabBar() string {
	tabs := []string{"Params", "Headers", "Body"}
	tabs[m.tab] = focusedStyle.Render(tabs[m.tab])
	separator := strings.Repeat(lipgloss.RoundedBorder().Bottom, m.width)
	separator = lipgloss.NewStyle().Foreground(lipgloss.Color(m.borderColor)).Render(separator)
	return strings.Join(tabs, " - ") + "\n" + separator + "\n"
}

func (m *RequestPaneModel) switchTab(direction int) {
	m.tab = requestPaneTab((int(m.tab) + direction + 3) % 3)
}

func (m RequestPaneModel) generateStyle() lipgloss.Style {
	var footer []string
	if m.tab == requestBodyTab && m.viewport.TotalLineCount() > 0 {
		footer = append(footer, fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	} else if m.tab == requestParamsTab || m.tab == requestHeadersTab {
		footer = append(footer, tableFooter(&m.table))
	}
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{
			Title:  []string{"[4]", "Request"},
			Footer: footer,
		},
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

func (m *RequestPaneModel) handleUpdateParam() {
	cursor, key, value, err := getKeyValueFromTableCursor(&m.table)
	if err != nil {
		return
	}
	m.textInputDialog.SetCmdFunc(updateParamCmdFunc(cursor, key))
	m.textInputDialog.SetPrompt(focusedStyle.Render(key + "="))
	m.textInputDialog.SetValue(value)
	m.textInputDialog.Focus()
	m.dctx.SetDialog(&m.textInputDialog)
}

func (m *RequestPaneModel) handleNewParam() {
	m.doubleTextInputDialog.SetCmdFunc(newParamCmdFunc)
	m.doubleTextInputDialog.FocusUpper()
	m.dctx.SetDialog(&m.doubleTextInputDialog)
}

func (m *RequestPaneModel) handleDeleteParam() tea.Cmd {
	cursor, _, _, err := getKeyValueFromTableCursor(&m.table)
	if err != nil {
		return nil
	}
	return messages.UpdateRequestCmd(func(r *internal.Request) {
		r.RemoveParamI(cursor)
	})
}

func (m *RequestPaneModel) handleUpdateHeader() {
	cursor, key, value, err := getKeyValueFromTableCursor(&m.table)
	if err != nil {
		return
	}
	m.textInputDialog.SetCmdFunc(updateHeaderCmdFunc(cursor, key))
	m.textInputDialog.SetPrompt(focusedStyle.Render(key + "="))
	m.textInputDialog.SetValue(value)
	m.textInputDialog.Focus()
	m.dctx.SetDialog(&m.textInputDialog)
}

func (m *RequestPaneModel) handleNewHeader() {
	m.doubleTextInputDialog.SetCmdFunc(newHeaderCmdFunc)
	m.doubleTextInputDialog.FocusUpper()
	m.dctx.SetDialog(&m.doubleTextInputDialog)
}

func (m *RequestPaneModel) handleDeleteHeader() tea.Cmd {
	cursor, _, _, err := getKeyValueFromTableCursor(&m.table)
	if err != nil {
		return nil
	}
	return messages.UpdateRequestCmd(func(r *internal.Request) {
		r.RemoveHeaderI(cursor)
	})
}

func (m *RequestPaneModel) handleUpdateBody() {
	if m.rctx.Empty() {
		return
	}
	request := m.rctx.Request()
	body := string(request.Body)
	m.textAreaDialog.SetCmdFunc(updateBodyCmdFunc)
	m.textAreaDialog.SetValue(body)
	m.textAreaDialog.Focus()
	m.dctx.SetDialog(&m.textAreaDialog)
}

func (m *RequestPaneModel) handleDeleteBody() tea.Cmd {
	return messages.UpdateRequestCmd(func(r *internal.Request) {
		r.Body = []byte{}
	})
}

// Refresh refreshes the table items based on the current tab.
func (m *RequestPaneModel) Refresh() {
	rows := make([]table.Row, 0)
	if m.rctx.Empty() {
		m.viewport.SetContent("")
		m.table.SetRows(rows)
		return
	}
	switch m.tab {
	case requestParamsTab:
		for _, kv := range m.rctx.Request().Params {
			rows = append(rows, table.Row{kv.Key, kv.Value})
		}
		m.table.SetRows(rows)
	case requestHeadersTab:
		for _, kv := range m.rctx.Request().Headers {
			rows = append(rows, table.Row{kv.Key, kv.Value})
		}
		m.table.SetRows(rows)
	case requestBodyTab:
		body := string(m.rctx.Request().Body)
		body = styles.ColorizeJsonIfValid(body)
		m.viewport.SetContent(body)
	default:
		m.table.SetRows(rows)
	}
}

func (m RequestPaneModel) Update(msg tea.Msg) (RequestPaneModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.dctx.Empty() {
			switch msg.String() {
			case "esc":
				return m, messages.SetFocusCmd(views.CollectionPaneView)
			case "[", "shift+tab":
				m.switchTab(-1)
			case "]", "tab":
				m.switchTab(1)
			}

			if !m.rctx.Empty() {
				switch msg.String() {
				case "x":
					return m, messages.ExecuteRequestCmd
				case "enter":
					switch m.tab {
					case requestParamsTab:
						m.handleUpdateParam()
					case requestHeadersTab:
						m.handleUpdateHeader()
					case requestBodyTab:
						m.handleUpdateBody()
					}
				case "n":
					switch m.tab {
					case requestParamsTab:
						m.handleNewParam()
					case requestHeadersTab:
						m.handleNewHeader()
					case requestBodyTab:
						m.handleUpdateBody()
					}
				case "d":
					switch m.tab {
					case requestParamsTab:
						cmds = append(cmds, m.handleDeleteParam())
					case requestHeadersTab:
						cmds = append(cmds, m.handleDeleteHeader())
					case requestBodyTab:
						cmds = append(cmds, m.handleDeleteBody())
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		verticalMarginHeight := 9
		m.viewport.Width = msg.Width - 2
		m.viewport.Height = msg.Height - verticalMarginHeight
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *RequestPaneModel) Blur() {
	m.table.SetStyles(tableBlurStyles())
}

func (m *RequestPaneModel) Focus() {
	m.table.SetStyles(tableStyles())
}

func (m RequestPaneModel) View() string {
	var text string
	text = m.renderTabBar()
	if !m.rctx.Empty() {
		switch m.tab {
		case requestParamsTab, requestHeadersTab:
			text += renderTableWithoutHeader(&m.table)
		case requestBodyTab:
			text += m.viewport.View()
		}
	}
	return m.generateStyle().Render(text)
}
