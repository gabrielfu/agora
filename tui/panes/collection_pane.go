package panes

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
	"github.com/gabrielfu/tipi/tui/states"
)

type CollectionPaneModel struct {
	width       int
	height      int
	borderColor string

	requests []internal.Request
	table    table.Model
	ctx      *states.RequestContext
}

func NewCollectionPaneModel(ctx *states.RequestContext) CollectionPaneModel {
	s := table.DefaultStyles()
	s.Selected = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#61AFEF")).
		Bold(false)
	s.Cell = lipgloss.NewStyle()

	t := table.New(
		table.WithColumns(makeColumns(0)),
		table.WithRows(make([]table.Row, 0)),
		table.WithFocused(true),
		table.WithStyles(s),
	)

	return CollectionPaneModel{table: t, ctx: ctx}
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
		method := RenderMethod(request.Method)
		u := RenderURL(request.URL)
		rows = append(rows, table.Row{method, u})
	}
	m.table.SetRows(rows)
	m.Update(nil)
}

func (m CollectionPaneModel) footer() string {
	return strconv.Itoa(m.table.Cursor()+1) + " / " + strconv.Itoa(len(m.table.Rows()))
}

func (m CollectionPaneModel) generateStyle() lipgloss.Style {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{
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

func (m CollectionPaneModel) renderTableWithoutHeader() string {
	t := m.table.View()
	ts := strings.SplitN(t, "\n", 2)
	return ts[len(ts)-1]
}

func (m CollectionPaneModel) Update(msg tea.Msg) (CollectionPaneModel, tea.Cmd) {
	var cmd tea.Cmd
	// process key messages to the table model
	m.table, cmd = m.table.Update(msg)
	// retrieve the request object and set the context
	cursor := m.table.Cursor()
	if cursor >= 0 && cursor < len(m.requests) {
		m.ctx.SetRequest(&m.requests[cursor])
	}
	return m, cmd
}

func (m CollectionPaneModel) View() string {
	text := m.renderTableWithoutHeader()
	return m.generateStyle().Render(text)
}
