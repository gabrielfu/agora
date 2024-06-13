package messages

import tea "github.com/charmbracelet/bubbletea"

var (
	ExecuteRequestCmd = func() tea.Msg { return ExecuteRequestMsg{} }
)
