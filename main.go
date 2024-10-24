package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui"
)

func Run() error {
	var collectionStore *internal.CollectionStore
	var err error
	if len(os.Args) > 1 {
		rootDir := filepath.Join(os.Args[1], ".agora")
		collectionStore, err = internal.NewCollectionStore(rootDir)
	} else {
		collectionStore, err = internal.NewDefaultCollectionStore()
	}
	if err != nil {
		return fmt.Errorf("error initializing collection store: %v", err)
	}

	collectionRequestDir := collectionStore.CurrentCollectionRequestDir()
	requestStore, err := internal.NewRequestFileStore(collectionRequestDir)
	if err != nil {
		return fmt.Errorf("error initializing collection store: %v", err)
	}
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
