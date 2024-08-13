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
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".agora", "data", "requests")
	store, err := internal.NewRquestFileStore(dir)
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
