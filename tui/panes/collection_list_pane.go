package panes

import (
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/tui/dialogs"
	"github.com/gabrielfu/agora/tui/messages"
	"github.com/gabrielfu/agora/tui/states"
	"github.com/gabrielfu/agora/tui/styles"
	"github.com/gabrielfu/agora/tui/views"
)

var newCollectionCmdFunc dialogs.TextInputCmdFunc = func(collection string) tea.Cmd {
	return messages.CreateCollectionCmd(collection)
}

func updateCollectionCmdFunc(collection string) dialogs.TextInputCmdFunc {
	return func(newName string) tea.Cmd {
		return messages.UpdateCollectionCmd(collection, newName)
	}
}

type CollectionListPaneModel struct {
	width       int
	height      int
	borderColor string

	dctx           *states.DialogContext
	list           list.Model
	editNameDialog dialogs.TextInputDialog
}

func NewCollectionListPaneModel(dctx *states.DialogContext) CollectionListPaneModel {
	l := list.New([]list.Item{}, simpleItemDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowFilter(false)

	return CollectionListPaneModel{
		dctx: dctx,
		list: l,
		editNameDialog: dialogs.NewTextInputDialog(
			64,
			[]string{"Name"},
			nil,
			updateNameCmd,
			views.CollectionPaneView,
		),
	}
}

func (m *CollectionListPaneModel) SetWidth(width int) {
	m.width = width
	m.list.SetDelegate(simpleItemDelegate{Width: width - 2})
}

func (m *CollectionListPaneModel) SetHeight(height int) {
	m.height = height
}

func (m *CollectionListPaneModel) SetBorderColor(color string) {
	m.borderColor = color
}

func (m CollectionListPaneModel) footer() string {
	cursor := m.list.Cursor() + 1
	total := len(m.list.Items())
	cursorString := " -"
	if cursor <= total {
		cursorString = strconv.Itoa(cursor)
	}
	return cursorString + " / " + strconv.Itoa(total)
}

func (m CollectionListPaneModel) generateStyle() lipgloss.Style {
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{
			Title:  []string{"[0]", "Collections"},
			Footer: []string{m.footer()},
		},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
}

func (m *CollectionListPaneModel) SetCollections(collections []string) {
	var items []list.Item
	for _, collection := range collections {
		items = append(items, simpleItem{value: collection})
	}
	m.list.SetItems(items)
	m.Update(nil)
}

func (m *CollectionListPaneModel) handleSelectCollection(collection string) tea.Cmd {
	return tea.Batch(
		messages.SetCollectionCmd(collection),
		messages.SetFocusCmd(views.CollectionPaneView),
	)
}

func (m *CollectionListPaneModel) handleNewCollection() {
	m.editNameDialog.SetCmdFunc(newCollectionCmdFunc)
	m.editNameDialog.Focus()
	m.dctx.SetDialog(&m.editNameDialog)
}

func (m *CollectionListPaneModel) handleUpdateCollection() {
	item, ok := m.list.SelectedItem().(simpleItem)
	if !ok {
		return
	}
	m.editNameDialog.SetCmdFunc(updateCollectionCmdFunc(item.value))
	m.editNameDialog.SetValue(item.value)
	m.editNameDialog.Focus()
	m.dctx.SetDialog(&m.editNameDialog)
}

func (m *CollectionListPaneModel) handleDeleteCollection() tea.Cmd {
	item, ok := m.list.SelectedItem().(simpleItem)
	if !ok {
		return nil
	}
	return messages.DeleteCollectionCmd(item.value)
}

func (m CollectionListPaneModel) Update(msg tea.Msg) (CollectionListPaneModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item := m.list.SelectedItem().(simpleItem)
			cmds = append(cmds, m.handleSelectCollection(item.value))
		case "n":
			m.handleNewCollection()
		case "r":
			m.handleUpdateCollection()
		case "d":
			cmds = append(cmds, m.handleDeleteCollection())
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(m.width-2, m.height-3)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m CollectionListPaneModel) View() string {
	text := m.list.View()
	return m.generateStyle().Render(text)
}
