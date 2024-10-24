package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielfu/agora/internal"
	"github.com/gabrielfu/agora/tui"
)

const DEFAULT_COLLECTION_NAME = "default"

func Run() error {
	rootStore, err := internal.NewDefaultRootStore()
	if err != nil {
		return fmt.Errorf("error initializing root store: %v", err)
	}
	// TODO: handled renamed collection
	if !rootStore.CollectionExists(DEFAULT_COLLECTION_NAME) {
		if err = rootStore.CreateCollection(DEFAULT_COLLECTION_NAME); err != nil {
			return fmt.Errorf("error creating default collection: %v", err)
		}
	}
	requestStore := internal.NewRequestFileStore(rootStore)
	model := tui.NewRootModel(requestStore, tui.WithCollectionPaneWidth(0.33))
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
