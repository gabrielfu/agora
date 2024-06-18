package dialogs

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Dialog interface {
	SetWidth(int)
	SetHeight(int)
	Update(tea.Msg) (any, tea.Cmd)
	View() string
}
