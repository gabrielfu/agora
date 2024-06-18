package panes

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
	"github.com/gabrielfu/tipi/tui/dialogs"
	"github.com/gabrielfu/tipi/tui/messages"
	"github.com/gabrielfu/tipi/tui/states"
	"github.com/gabrielfu/tipi/tui/styles"
	"github.com/gabrielfu/tipi/tui/views"
)

var (
	tableSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color(styles.SelectedBackgroundColor))
	tableBlurSelectedStyle = lipgloss.NewStyle()
)

func tableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Selected = tableSelectedStyle
	s.Cell = lipgloss.NewStyle()
	return s
}

func tableBlurStyles() table.Styles {
	s := table.DefaultStyles()
	s.Selected = tableBlurSelectedStyle
	s.Cell = lipgloss.NewStyle()
	return s
}

type CollectionPaneModel struct {
	width       int
	height      int
	borderColor string

	requests       []internal.Request
	table          table.Model
	cursor         int
	rctx           *states.RequestContext
	dctx           *states.DialogContext
	editNameDialog dialogs.TextInputDialog
}

func NewCollectionPaneModel(rctx *states.RequestContext, dctx *states.DialogContext) CollectionPaneModel {
	t := table.New(
		table.WithColumns(makeColumns(0)),
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

	return CollectionPaneModel{table: t, rctx: rctx, dctx: dctx,
		editNameDialog: dialogs.NewTextInputDialog(
			64,
			[]string{"Name"},
			nil,
			updateNameCmd,
			views.CollectionPaneView,
		),
	}
}

func makeColumns(width int) []table.Column {
	return []table.Column{
		{Title: "Method", Width: 6},
		{Title: "URL", Width: max(0, width-6)},
	}
}

func (m *CollectionPaneModel) SetWidth(width int) {
	m.width = width
	m.table.SetWidth(width)
	m.table.SetColumns(makeColumns(width))
}

func (m *CollectionPaneModel) SetHeight(height int) {
	m.height = height
	m.table.SetHeight(height)
}

func (m *CollectionPaneModel) SetBorderColor(color string) {
	m.borderColor = color
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
		} else {
			display = styles.RenderURL(request.URL)
		}
		rows = append(rows, table.Row{method, display})
	}
	m.table.SetRows(rows)
	m.Update(nil)
}

func (m CollectionPaneModel) footer() string {
	return strconv.Itoa(m.table.Cursor()+1) + " / " + strconv.Itoa(len(m.table.Rows()))
}

func (m CollectionPaneModel) generateStyle() lipgloss.Style {
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{
			Title:  []string{"[1]", "Collection"},
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

func (m CollectionPaneModel) renderTableWithoutHeader() string {
	t := m.table.View()
	ts := strings.SplitN(t, "\n", 2)
	return ts[len(ts)-1]
}

func (m CollectionPaneModel) Update(msg tea.Msg) (CollectionPaneModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			return m, messages.ExecuteRequestCmd
		case "enter":
			return m, messages.SetFocusCmd(views.UrlPaneView)
		case "r":
			if !m.rctx.Empty() {
				m.editNameDialog.SetValue(m.rctx.Request().Name)
				m.editNameDialog.Focus()
				m.dctx.SetDialog(&m.editNameDialog)
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
	text := m.renderTableWithoutHeader()
	return m.generateStyle().Render(text)
}
