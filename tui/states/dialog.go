package states

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielfu/tipi/tui/views"
)

type Dialog interface {
	Prev() views.View
	Update(tea.Msg) (any, tea.Cmd)
	View() string
}

type DialogContext struct {
	dialog Dialog
}

func NewDialogContext() *DialogContext {
	return &DialogContext{}
}

func (d *DialogContext) Empty() bool {
	return d.dialog == nil
}

func (d *DialogContext) View() string {
	if d.dialog == nil {
		return ""
	}
	return d.dialog.View()
}

func (d *DialogContext) Dialog() Dialog {
	return d.dialog
}

func (d *DialogContext) SetDialog(dialog Dialog) {
	d.dialog = dialog
}

func (d *DialogContext) Clear() {
	d.dialog = nil
}
