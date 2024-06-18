package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielfu/tipi/internal"
	"github.com/gabrielfu/tipi/tui"
)

func Run() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".tipi")
	if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	dbFile := filepath.Join(dir, "data.sqlite")
	db, err := internal.NewRquestDatabase(dbFile)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	model := tui.NewRootModel(db, tui.WithCollectionPaneWidth(0.33))
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
