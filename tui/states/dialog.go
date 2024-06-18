package states

import "github.com/gabrielfu/agora/tui/dialogs"

type DialogContext struct {
	dialog dialogs.Dialog
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

func (d *DialogContext) SetDialogWidth(width int) {
	if d.dialog != nil {
		d.dialog.SetWidth(width)
	}
}

func (d *DialogContext) SetDialogHeight(height int) {
	if d.dialog != nil {
		d.dialog.SetHeight(height)
	}
}

func (d *DialogContext) Dialog() dialogs.Dialog {
	return d.dialog
}

func (d *DialogContext) SetDialog(dialog dialogs.Dialog) {
	d.dialog = dialog
}

func (d *DialogContext) Clear() {
	d.dialog = nil
}
