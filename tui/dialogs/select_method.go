package dialogs

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/tipi/tui/messages"
	"github.com/gabrielfu/tipi/tui/styles"
	"github.com/gabrielfu/tipi/tui/views"
)

var methodColors = map[string]string{
	"GET":     "#68D696",
	"POST":    "#EED577",
	"PUT":     "#74AEF6",
	"PATCH":   "#C0A8E1",
	"DELETE":  "#EF968A",
	"HEAD":    "#68D696",
	"OPTIONS": "#E55AA8",
}

func GetMethodColor(method string) string {
	if color, ok := methodColors[method]; ok {
		return color
	}
	return methodColors["GET"]
}

var (
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().Background(lipgloss.Color(styles.SelectedBackgroundColor))
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(item)
	if !ok {
		return
	}
	method := fmt.Sprintf("%-8s", string(it))
	color := GetMethodColor(string(it))
	fn := itemStyle.Foreground(lipgloss.Color(color)).Render
	if index == m.Index() {
		fn = selectedItemStyle.Render
	}
	fmt.Fprint(w, fn(method))
}

var methods = []list.Item{
	item("GET"),
	item("POST"),
	item("PUT"),
	item("PATCH"),
	item("DELETE"),
	item("HEAD"),
	item("OPTIONS"),
}

type SelectMethodDialog struct {
	width int
	list  list.Model
}

func NewSelectMethodDialog() SelectMethodDialog {
	width := 10
	l := list.New(methods, itemDelegate{}, width, 7)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowFilter(false)
	return SelectMethodDialog{list: l, width: width}
}

func (m *SelectMethodDialog) SetWidth(width int) {
	m.width = width
}

func (m SelectMethodDialog) generateStyle() lipgloss.Style {
	border := styles.GenerateBorder(
		lipgloss.RoundedBorder(),
		styles.GenerateBorderOption{Title: []string{"Method"}},
		m.width,
	)
	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(lipgloss.Color(styles.FocusBorderColor)).
		Width(m.width).
		Padding(0, 1)
}

func (m SelectMethodDialog) Update(msg tea.Msg) (any, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd := func() tea.Msg {
				return messages.ExitDialogMsg{Dest: views.UrlPaneView}
			}
			return m, cmd
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SelectMethodDialog) View() string {
	return m.generateStyle().Render(m.list.View())
}