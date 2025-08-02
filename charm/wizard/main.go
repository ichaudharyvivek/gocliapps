package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	txt "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
	Question    lipgloss.Style
}

func DefaultStyle() *Styles {
	s := &Styles{}

	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(62)
	s.Question = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	return s
}

type Survey struct {
	question string
	answer   string
}

func NewSurvey(q string) *Survey {
	return &Survey{
		question: q,
		answer:   "",
	}
}

type AppState int

const (
	StateAsking AppState = iota
	StateCompleted
)

type Model struct {
	input      txt.Model
	styles     *Styles
	surveys    []*Survey
	state      AppState
	currentIdx int
	width      int
	height     int
}

func New(questions []string) *Model {
	styles := DefaultStyle()

	surveys := make([]*Survey, len(questions))
	for i, q := range questions {
		surveys[i] = NewSurvey(q)
	}

	input := txt.New()
	input.Placeholder = "Your answer here..."
	input.CharLimit = 60
	input.Width = 60
	input.Focus()

	return &Model{
		surveys:    surveys,
		input:      input,
		styles:     styles,
		state:      StateAsking,
		currentIdx: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return txt.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch m.state {

		case StateAsking:
			return m.handleAskingState(msg)
		case StateCompleted:
			return m.handleCompletedState(msg)
		}
	}

	if m.state == StateAsking {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleAskingState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Save the current answer
		answer := strings.TrimSpace(m.input.Value())
		if answer == "" {
			return m, nil
		}

		m.surveys[m.currentIdx].answer = answer

		// Move to next question or complete survey
		if m.currentIdx < len(m.surveys)-1 {
			m.currentIdx++
			m.input.SetValue("")
			m.input.Focus()
		} else {
			m.state = StateCompleted
		}

		return m, nil

	case "ctrl+c", "esc":
		return m, tea.Quit
	}

	// Let the input handle other keys
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) handleCompletedState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading...."
	}

	switch m.state {

	case StateAsking:
		return m.renderAskingView()
	case StateCompleted:
		return m.renderCompletedView()
	}

	return ""
}

func (m Model) renderAskingView() string {
	current := m.surveys[m.currentIdx]
	progress := fmt.Sprintf("Question %d of %d", m.currentIdx+1, len(m.surveys))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.styles.Question.Render(current.question),
		m.styles.InputField.Render(m.input.View()),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(progress),
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press Enter to continue • Ctrl+C to quit"),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderCompletedView() string {
	s := ""
	for i, survey := range m.surveys {
		s += fmt.Sprintf("\n[Q%d.] %s: %s", i+1, survey.question, survey.answer)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Render(s),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("\n\nPress ctrl+c or esc to exit."),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, content)
}

func main() {
	questions := []string{
		"What is your name?",
		"What is your favorite code editor?",
		"What is your favorite programming language?",
		"What is your favorite quote?",
		"What motivates you to code?",
	}

	model := New(questions)

	// Setup debug logging (optional)
	f, err := tea.LogToFile("debug.log", "")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create and run the program
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
