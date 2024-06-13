package tui

import "github.com/elliotchance/orderedmap/v2"

type Keymap = orderedmap.OrderedMap[string, string]

var (
	EmptyKeymap              *Keymap = orderedmap.NewOrderedMap[string, string]()
	CollectionPaneKeymap     *Keymap = orderedmap.NewOrderedMap[string, string]()
	UrlPaneKeymap            *Keymap = orderedmap.NewOrderedMap[string, string]()
	SelectMethodDialogKeymap *Keymap = orderedmap.NewOrderedMap[string, string]()
	TextInputDialogKeymap    *Keymap = orderedmap.NewOrderedMap[string, string]()
)

func init() {
	// CollectionPaneKeymap.Set("<enter>", "Select")
	CollectionPaneKeymap.Set("e", "Execute")
	// CollectionPaneKeymap.Set("n", "New")
	// CollectionPaneKeymap.Set("d", "Delete")

	UrlPaneKeymap.Set("m", "Select method")
	UrlPaneKeymap.Set("u", "Edit url")
	// UrlPaneKeymap.Set("e", "Execute")

	SelectMethodDialogKeymap.Set("<enter>", "Select")
	SelectMethodDialogKeymap.Set("<esc>", "Cancel")

	TextInputDialogKeymap.Set("<enter>", "Submit")
	TextInputDialogKeymap.Set("<esc>", "Cancel")
}
