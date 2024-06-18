package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui/views"
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
	CreateRequestCmd = func(r internal.Request) tea.Cmd {
		return func() tea.Msg { return CreateRequestMsg{Req: r} }
	}
	DeleteRequestCmd = func(id string) tea.Cmd {
		return func() tea.Msg { return DeleteRequestMsg{ID: id} }
	}
)
