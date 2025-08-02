package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const URL = "https://charm.sh/"

type model struct {
	status   int
	err      error
	checking bool
}

func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}

	res, err := c.Get(URL)
	if err != nil {
		return errMsg{err}
	}

	return statusMsg(res.StatusCode)
}

type (
	statusMsg int
	errMsg    struct{ err error }
)

func (m model) Init() tea.Cmd {
	return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.status = int(msg)
		m.checking = false
		return m, nil

	case errMsg:
		m.err = msg.err
		m.checking = false
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("We have some trouble: %v\n", m.err)
	}

	s := fmt.Sprintf("Checking: %s", URL)
	if m.status > 0 {
		s += fmt.Sprintf("\n%d: %s", m.status, http.StatusText(m.status))
	} else {
		s += "\nLoading..."
	}

	return s
}

func main() {
	// Initialize model with checking state
	m := model{checking: true}
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
