package panes

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui/dialogs"
	"github.com/gabrielfu/agora/tui/messages"
	"github.com/gabrielfu/agora/tui/states"
	"github.com/gabrielfu/agora/tui/styles"
	"github.com/gabrielfu/agora/tui/views"
)

type CollectionPaneModel struct {
	width       int
	height      int
	borderColor string

	collection     string
	requests       []internal.Request
	table          table.Model
	cursor         int
	rctx           *states.RequestContext
	dctx           *states.DialogContext
	editNameDialog dialogs.TextInputDialog
}

func NewCollectionPaneModel(rctx *states.RequestContext, dctx *states.DialogContext, collection string) CollectionPaneModel {
	t := table.New(
		table.WithColumns(makeCollectionColumns(0)),
		table.WithRows(make([]table.Row, 0)),
		table.WithFocused(true),
		table.WithStyles(tableStyles()),
		// // Wait until StyleFunc is supported on bubbles/table
		// table.WithStyleFunc(func(row, col int, value string) lipgloss.Style {
		// 	if col == 0 { // is method column
		// 		color := styles.GetMethodColor(value)
		// 		return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		// 	}
		// 	return lipgloss.NewStyle()
		// }),
	)

	// disable "u" and "d"
	t.KeyMap.HalfPageUp.SetEnabled(false)
	t.KeyMap.HalfPageDown.SetEnabled(false)

	return CollectionPaneModel{
		table:      t,
		collection: collection,
		rctx:       rctx,
		dctx:       dctx,
		editNameDialog: dialogs.NewTextInputDialog(
			64,
			[]string{"Name"},
			nil,
			updateNameCmd,
			views.CollectionPaneView,
		),
	}
}

func makeCollectionColumns(width int) []table.Column {
	return []table.Column{
		{Title: "Method", Width: 6},
		{Title: "URL", Width: max(0, width-6)},
	}
}

func (m *CollectionPaneModel) SetWidth(width int) {
	m.width = width
	m.table.SetWidth(width)
	m.table.SetColumns(makeCollectionColumns(width))
}

func (m *CollectionPaneModel) SetHeight(height int) {
	m.height = height
	m.table.SetHeight(height)
}

func (m *CollectionPaneModel) SetBorderColor(color string) {
	m.borderColor = color
}

func (m *CollectionPaneModel) SetCollection(collection string) {
	m.collection = collection
}

func (m *CollectionPaneModel) SetRequests(requests []internal.Request) {
	m.requests = requests
	var rows []table.Row
	for _, request := range requests {
		// TODO: cell level color doesn't work yet for bubbles table
		method := styles.RenderMethod(request.Method)
		var display string
		if request.Name != "" {
			display = request.Name
		} else if request.URL != "" {
			display = styles.RenderURL(request.URL)
		} else {
			display = "untitled"
		}
		rows = append(rows, table.Row{method, display})
	}
	m.table.SetRows(rows)
	m.Update(nil)
}

func (m CollectionPaneModel) footer() string {
	cursor := m.table.Cursor() + 1
	total := len(m.table.Rows())
	cursorString := " -"
	if cursor <= total {
		cursorString = strconv.Itoa(cursor)
	}
	return cursorString + " / " + strconv.Itoa(total)
}

func (m CollectionPaneModel) generateStyle() lipgloss.Style {
	title := []string{"[1]", "Collection", "(" + m.collection + ")"}
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{
			Title:  title,
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

func (m *CollectionPaneModel) Blur() {
	m.table.SetStyles(tableBlurStyles())
}

func (m *CollectionPaneModel) Focus() {
	m.table.SetStyles(tableStyles())
}

func (m CollectionPaneModel) Update(msg tea.Msg) (CollectionPaneModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "x":
			return m, messages.ExecuteRequestCmd
		case "enter":
			return m, messages.SetFocusCmd(views.UrlPaneView)
		case "n":
			return m, messages.CreateRequestCmd(*internal.NewRequest("GET", ""))
		case "r":
			if !m.rctx.Empty() {
				m.editNameDialog.SetValue(m.rctx.Request().Name)
				m.editNameDialog.Focus()
				m.dctx.SetDialog(&m.editNameDialog)
			}
		case "d":
			if !m.rctx.Empty() {
				return m, messages.DeleteRequestCmd(m.rctx.Request().ID)
			}
		case "c":
			if !m.rctx.Empty() {
				return m, messages.CopyRequestCmd(*m.rctx.Request())
			}
		}
	}

	// process key messages to the table model
	m.table, cmd = m.table.Update(msg)
	// retrieve the request object and set the context
	cursor := m.table.Cursor()
	if m.cursor != cursor {
		m.rctx.Clear()
		m.cursor = cursor
	}
	if cursor >= 0 && cursor < len(m.requests) {
		m.rctx.SetRequest(&m.requests[cursor])
	}
	return m, cmd
}

func (m CollectionPaneModel) View() string {
	text := renderTableWithoutHeader(&m.table)
	return m.generateStyle().Render(text)
}
