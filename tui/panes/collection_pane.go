package panes

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/internal"
)

type CollectionPaneModel struct {
	width       int
	height      int
	borderColor string

	requests []internal.Request
	table    table.Model
}

func NewCollectionPaneModel() CollectionPaneModel {
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

	return CollectionPaneModel{table: t}
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
		method := RenderMethod(request.Method)
		u := RenderURL(request.URL)
		rows = append(rows, table.Row{method, u})
	}
	m.table.SetRows(rows)
}

func (m CollectionPaneModel) generateStyle() lipgloss.Style {
	border := generateBorder(
		lipgloss.RoundedBorder(),
		GenerateBorderOption{Title: []string{"[1]", "Collection"}},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(m.borderColor)).
		Width(m.width).
		Height(m.height)
}

func (m CollectionPaneModel) Update(msg tea.Msg) (CollectionPaneModel, tea.Cmd) {
	return m, nil
}

func (m CollectionPaneModel) renderTable() string {
	t := m.table.View()
	ts := strings.SplitN(t, "\n", 2)
	return ts[len(ts)-1]
}

func (m CollectionPaneModel) View() string {
	t := m.renderTable()
	return m.generateStyle().Render(t)
}
