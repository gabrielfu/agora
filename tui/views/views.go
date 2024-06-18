package views

type View uint

const (
	// Pane views
	CollectionPaneView View = iota
	UrlPaneView
	RequestPaneView
	ResponsePaneView
	// Dialog views
	SelectMethodDialogView
	TextInputDialogView
	TextAreaDialogView
)

func IsPaneView(v View) bool {
	return v <= ResponsePaneView
}

func IsDialogView(v View) bool {
	return v >= SelectMethodDialogView
}
