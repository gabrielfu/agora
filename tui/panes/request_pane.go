package panes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
	"github.com/gabrielfu/tipi/tui/dialogs"
	"github.com/gabrielfu/tipi/tui/messages"
	"github.com/gabrielfu/tipi/tui/states"
	"github.com/gabrielfu/tipi/tui/styles"
	"github.com/gabrielfu/tipi/tui/views"
)

var focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.FocusBorderColor))

type requestPaneTab int

const (
	paramsTab requestPaneTab = iota
	headersTab
	bodyTab
)

func updateParamCmdFunc(key, originalValue string) dialogs.TextInputCmdFunc {
	return func(value string) tea.Cmd {
		return messages.UpdateRequestCmd(func(r *internal.Request) {
			r.RemoveParam(key, originalValue)
			r.WithParam(key, value)
		})
	}
}

func newParamCmdFunc() dialogs.DoubleTextInputCmdFunc {
	return func(key, value string) tea.Cmd {
		return messages.UpdateRequestCmd(func(r *internal.Request) {
			r.WithParam(key, value)
		})
	}
}

type kvItem struct {
	key, value string
}

func (i kvItem) Title() string       { return i.key }
func (i kvItem) Description() string { return i.value }
func (i kvItem) FilterValue() string { return i.key }

type RequestPaneModel struct {
	width       int
	height      int
	borderColor string

	rctx *states.RequestContext
	dctx *states.DialogContext

	tab                   requestPaneTab
	textInputDialog       dialogs.TextInputDialog
	doubleTextInputDialog dialogs.DoubleTextInputDialog
	list                  list.Model
}

func NewRequestPaneModel(rctx *states.RequestContext, dctx *states.DialogContext) RequestPaneModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	// l.SetShowPagination(false)
	l.SetShowFilter(false)
	return RequestPaneModel{
		rctx: rctx,
		dctx: dctx,
		tab:  paramsTab,
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
		list: l,
	}
}

func (m *RequestPaneModel) SetWidth(width int) {
	m.width = width
}

func (m *RequestPaneModel) SetHeight(height int) {
	m.height = height
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
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: []string{"[3]", "Request"}},
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
	item, ok := m.list.SelectedItem().(kvItem)
	if !ok {
		return
	}
	key, value := item.key, item.value
	m.textInputDialog.SetCmdFunc(updateParamCmdFunc(key, value))
	m.textInputDialog.SetPrompt(focusedStyle.Render(key + "="))
	m.textInputDialog.SetValue(value)
	m.textInputDialog.Focus()
	m.dctx.SetDialog(&m.textInputDialog)
}

func (m *RequestPaneModel) handleNewParam() {
	m.doubleTextInputDialog.SetCmdFunc(newParamCmdFunc())
	m.doubleTextInputDialog.FocusUpper()
	m.dctx.SetDialog(&m.doubleTextInputDialog)
}

func (m *RequestPaneModel) Refresh() {
	// refresh param list
	if !m.rctx.Empty() {
		switch m.tab {
		case paramsTab:
			var items []list.Item
			for _, kv := range m.rctx.Request().Params {
				items = append(items, kvItem{key: kv.Key, value: kv.Value})
			}
			m.list.SetItems(items)
		case headersTab:
			var items []list.Item
			for _, kv := range m.rctx.Request().Headers {
				items = append(items, kvItem{key: kv.Key, value: kv.Value})
			}
			m.list.SetItems(items)
		default:
			m.list.SetItems([]list.Item{})
		}
	}
}

func (m RequestPaneModel) Update(msg tea.Msg) (RequestPaneModel, tea.Cmd) {
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
			case "enter":
				switch m.tab {
				case paramsTab:
					m.handleUpdateParam()
				}
			case "n":
				m.handleNewParam()
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(m.width-2, m.height-3)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m RequestPaneModel) View() string {
	var text string
	text = m.renderTabBar()
	if !m.rctx.Empty() {
		switch m.tab {
		case paramsTab, headersTab:
			text += m.list.View()
		case bodyTab:
			text += fmt.Sprintf("%v", m.rctx.Request().Body)
		}
	}
	return m.generateStyle().Render(text)
}
