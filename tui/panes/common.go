package panes

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielfu/agora/tui/styles"
)

type kvItem struct {
	key, value string
}

func (i kvItem) Title() string       { return i.key }
func (i kvItem) Description() string { return i.value }
func (i kvItem) FilterValue() string { return i.key }

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
