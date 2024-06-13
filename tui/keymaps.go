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
	SelectMethodDialogKeymap = NewKeymap()
	TextInputDialogKeymap    = NewKeymap()
)

func init() {
	CollectionPaneKeymap.Set("<enter>", "Select")
	CollectionPaneKeymap.Set("e", "Execute")
	// CollectionPaneKeymap.Set("n", "New")
	// CollectionPaneKeymap.Set("d", "Delete")

	UrlPaneKeymap.Set("m", "Select method")
	UrlPaneKeymap.Set("u", "Edit url")
	// UrlPaneKeymap.Set("e", "Execute")
	UrlPaneKeymap.Set("<esc>", "Back")

	SelectMethodDialogKeymap.Set("<enter>", "Select")
	SelectMethodDialogKeymap.Set("<esc>", "Cancel")

	TextInputDialogKeymap.Set("<enter>", "Submit")
	TextInputDialogKeymap.Set("<esc>", "Cancel")
}
