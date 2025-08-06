package tui

import (
	"manager/internal/db"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg struct{ error }
)

var (
	AppStore   *db.Repositories
	WindowSize tea.WindowSizeMsg
)

/* STYLING */
var (
	DocStyle   = lipgloss.NewStyle().Margin(0, 2)
	HelpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	ErrStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd534b")).Render
	AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render
)

// keymap reusable key mappings shared across models
type keymap struct {
	Create key.Binding
	Enter  key.Binding
	Rename key.Binding
	Delete key.Binding
	Back   key.Binding
	Quit   key.Binding
}

var Keymap = &keymap{
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
