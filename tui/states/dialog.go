package states

import "github.com/gabrielfu/tipi/tui/views"

type Viewable interface {
	Prev() views.View
	View() string
}

type DialogContext struct {
	dialog Viewable
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

func (d *DialogContext) Dialog() Viewable {
	return d.dialog
}

func (d *DialogContext) SetDialog(dialog Viewable) {
	d.dialog = dialog
}

func (d *DialogContext) Clear() {
	d.dialog = nil
}
