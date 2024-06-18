package tui

import "github.com/elliotchance/orderedmap/v2"

type Keymap = orderedmap.OrderedMap[string, string]

func NewKeymap() *Keymap {
	return orderedmap.NewOrderedMap[string, string]()
}

var (
	EmptyKeymap              = NewKeymap()
	CollectionPaneKeymap     = NewKeymap()
	UrlPaneKeymap            = NewKeymap()
	RequestPaneKeymap        = NewKeymap()
	ResponsePaneKeymap       = NewKeymap()
	SelectMethodDialogKeymap = NewKeymap()
	TextInputDialogKeymap    = NewKeymap()
)

func init() {
	CollectionPaneKeymap.Set("<enter>", "Select")
	CollectionPaneKeymap.Set("x", "Execute")
	// CollectionPaneKeymap.Set("n", "New")
	CollectionPaneKeymap.Set("r", "Rename")
	// CollectionPaneKeymap.Set("d", "Delete")

	UrlPaneKeymap.Set("x", "Execute")
	UrlPaneKeymap.Set("m", "Select method")
	UrlPaneKeymap.Set("<enter>", "Edit")
	UrlPaneKeymap.Set("r", "Rename")
	UrlPaneKeymap.Set("<esc>", "Back")

	RequestPaneKeymap.Set("<enter>", "Edit")
	RequestPaneKeymap.Set("n", "New")
	RequestPaneKeymap.Set("d", "Delete")
	RequestPaneKeymap.Set("<esc>", "Back")

	ResponsePaneKeymap.Set("<esc>", "Back")

	SelectMethodDialogKeymap.Set("<enter>", "Select")
	SelectMethodDialogKeymap.Set("<esc>", "Cancel")

	TextInputDialogKeymap.Set("<enter>", "Submit")
	TextInputDialogKeymap.Set("<esc>", "Cancel")
}
