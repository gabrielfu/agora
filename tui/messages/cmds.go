package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielfu/tipi/internal"
	"github.com/gabrielfu/tipi/tui/views"
)

var (
	SetFocusCmd = func(v views.View) tea.Cmd {
		return func() tea.Msg { return SetFocusMsg{View: v} }
	}
	ExitDialogCmd = func(v views.View) tea.Cmd {
		return func() tea.Msg { return ExitDialogMsg{Dest: v} }
	}
	ExecuteRequestCmd tea.Cmd = func() tea.Msg { return ExecuteRequestMsg{} }
	UpdateRequestCmd          = func(f func(*internal.Request)) tea.Cmd {
		return func() tea.Msg { return UpdateRequestMsg{Func: f} }
	}
)