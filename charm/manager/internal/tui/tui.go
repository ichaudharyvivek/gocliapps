package tui

import (
	"fmt"
	"manager/internal/db"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// StartTea is the entry point of the UI
// It initializes the model and start the program
func StartTea(store *db.Repositories) {
	AppStore = store
	f, err := tea.LogToFile("debug.log", "help")
	if err != nil {
		fmt.Printf("could not open log file to write: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	m, _ := InitProject()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("something went wrong: %v\n", err)
		os.Exit(1)
	}
}
