package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui"
)

func Run() error {
	collectionStore := internal.DefaultCollectionStore
	if err := collectionStore.SetDefaultRoot(); err != nil {
		return err
	}
	dir := collectionStore.CollectionRequestDir()
	store, err := internal.NewRequestFileStore(dir)
	if err != nil {
		return fmt.Errorf("error initializing file store: %v", err)
	}

	model := tui.NewRootModel(store, tui.WithCollectionPaneWidth(0.33))
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
