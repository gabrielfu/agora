package messages

import (
	"github.com/gabrielfu/tipi/tui/views"
)

type SetFocusMsg struct {
	View views.View
}

type ExitDialogMsg struct {
	Dest views.View
}
