package messages

import (
	"fmt"

	"github.com/gabrielfu/tipi/tui/views"
)

type SetFocusMsg struct {
	View views.View
}

func (m SetFocusMsg) String() string {
	return fmt.Sprintf("SetFocusMsg(%v)", m.View)
}

type ExitDialogMsg struct {
	Dest views.View
}
