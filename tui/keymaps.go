package tui

import "github.com/elliotchance/orderedmap/v2"

type Keymap = orderedmap.OrderedMap[string, string]

var (
	EmptyKeymap              *Keymap = orderedmap.NewOrderedMap[string, string]()
	CollectionPaneKeymap     *Keymap = orderedmap.NewOrderedMap[string, string]()
	SelectMethodDialogKeymap *Keymap = orderedmap.NewOrderedMap[string, string]()
)

func init() {
	CollectionPaneKeymap.Set("<space>", "Select")
	CollectionPaneKeymap.Set("e", "Execute")
	CollectionPaneKeymap.Set("n", "New")
	CollectionPaneKeymap.Set("d", "Delete")

	SelectMethodDialogKeymap.Set("<space>", "Select")
	SelectMethodDialogKeymap.Set("<esc>", "Cancel")
}
