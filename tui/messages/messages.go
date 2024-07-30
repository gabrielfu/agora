package messages

import (
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui/views"
)

type SetFocusMsg struct {
	View views.View
}

type ExitDialogMsg struct {
	Dest views.View
}

type ExecuteRequestMsg struct{}

type UpdateRequestMsg struct {
	Func func(*internal.Request)
}

type CreateRequestMsg struct {
	Req internal.Request
}

type DeleteRequestMsg struct {
	ID string
}

type CopyRequestMsg struct {
	Req internal.Request
}
