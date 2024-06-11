package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielfu/tipi/internal"
	"github.com/gabrielfu/tipi/tui"
)

func main() {
	db, err := internal.NewRquestDatabase("tipi.db")
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	model := tui.NewRootModel(db, tui.WithCollectionPaneWidth(0.3))
	program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := program.Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
