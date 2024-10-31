package panes

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/tui/styles"
)

type simpleItem struct {
	value string
}

func (i simpleItem) Title() string       { return i.value }
func (i simpleItem) Description() string { return i.value }
func (i simpleItem) FilterValue() string { return i.value }

var (
	simpleItemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedSimpleItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color(styles.SelectedBackgroundColor))
)

type simpleItemDelegate struct {
	Width         int
	SelectedStyle lipgloss.Style
}

func (d simpleItemDelegate) Height() int                             { return 1 }
func (d simpleItemDelegate) Spacing() int                            { return 0 }
func (d simpleItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d simpleItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(simpleItem)
	if !ok {
		return
	}
	fn := simpleItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			str := strings.Join(s, " ")
			// pad the item to the right
			if len(str) < d.Width {
				str = str + strings.Repeat(" ", d.Width-len(str))
			}
			return d.SelectedStyle.Render(str)
		}
	}
	fmt.Fprint(w, fn(i.value))
}

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

func renderTableWithoutHeader(table *table.Model) string {
	t := table.View()
	ts := strings.SplitN(t, "\n", 2)
	return ts[len(ts)-1]
}

func makeKeyValueColumns(width int) []table.Column {
	keyWidth := int(float64(width) * 0.4)
	return []table.Column{
		{Title: "Key", Width: keyWidth},
		{Title: "Value", Width: width - keyWidth},
	}
}

func getKeyValueFromTableCursor(table *table.Model) (
	cursor int, key string, value string, err error,
) {
	cursor = table.Cursor()
	if cursor < 0 || cursor >= len(table.Rows()) {
		err = errors.New("cursor out of bounds")
		return
	}
	row := table.Rows()[cursor]
	if len(row) != 2 {
		err = errors.New("invalid row")
		return
	}
	key = row[0]
	value = row[1]
	return
}
