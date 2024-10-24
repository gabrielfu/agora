package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui"
)

func Run() error {
	collectionStore, err := internal.NewDefaultCollectionStore()
	if err != nil {
		return fmt.Errorf("error initializing collection store: %v", err)
	}
	dir := collectionStore.CurrentCollectionRequestDir()
	requestStore := internal.NewRequestFileStore(dir)
	model := tui.NewRootModel(collectionStore, requestStore, tui.WithCollectionPaneWidth(0.33))
	program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err = program.Run()
	return err
}

func main() {
	if err := Run(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
