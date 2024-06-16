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

// fix support for multiple params with same key
func updateParamCmdFunc(key string) dialogs.TextInputCmdFunc {
	return func(value string) tea.Cmd {
		return messages.UpdateRequestCmd(func(r *internal.Request) {
			r.Params[key] = value
		})
	}
}

type defaultItem struct {
	key, value string
}

func (i defaultItem) Title() string       { return i.key }
func (i defaultItem) Description() string { return i.value }
func (i defaultItem) FilterValue() string { return i.key }

type RequestPaneModel struct {
	width       int
	height      int
	borderColor string

	rctx *states.RequestContext
	dctx *states.DialogContext

	tab             requestPaneTab
	editParamDialog dialogs.TextInputDialog
	list            list.Model
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
		editParamDialog: dialogs.NewTextInputDialog(
			64,
			[]string{"Param"},
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
	item, ok := m.list.SelectedItem().(defaultItem)
	if !ok {
		return
	}
	key, value := item.key, item.value
	m.editParamDialog.SetCmdFunc(updateParamCmdFunc(key))
	m.editParamDialog.SetPrompt(focusedStyle.Render(key + "="))
	m.editParamDialog.SetValue(value)
	m.editParamDialog.Focus()
	m.dctx.SetDialog(&m.editParamDialog)
}

func (m *RequestPaneModel) Refresh() {
	// refresh param list
	if !m.rctx.Empty() {
		switch m.tab {
		case paramsTab:
			var items []list.Item
			for k, v := range m.rctx.Request().Params {
				items = append(items, defaultItem{key: k, value: v})
			}
			m.list.SetItems(items)
		case headersTab:
			var items []list.Item
			for k, v := range m.rctx.Request().Headers {
				items = append(items, defaultItem{key: k, value: v})
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
