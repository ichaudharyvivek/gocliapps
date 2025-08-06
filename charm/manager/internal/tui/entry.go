package tui

import (
	"fmt"
	"os"

	e "manager/internal/domain/entry"
	"manager/pkg/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type (
	updateEntries   []*e.Entry
	editorFinishMsg struct {
		err  error
		file *os.File
	}
)

type EntryModel struct {
	quitting  bool
	activePID uint
	error     string
	entries   []*e.Entry
	viewport  viewport.Model
	paginator paginator.Model
}

func InitEntry(pID uint, repo e.Repository) *EntryModel {
	m := &EntryModel{
		activePID: pID,
	}

	// Set viewport
	top, right, bottom, left := DocStyle.GetMargin()
	m.viewport = viewport.New(WindowSize.Width-left-right, WindowSize.Height-top-bottom-1)
	m.viewport.Style = lipgloss.NewStyle().Align(lipgloss.Bottom)

	// Init paginator
	m.paginator = paginator.New()
	m.paginator.Type = paginator.Dots
	m.paginator.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	m.paginator.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	// Set viewport and content
	m.entries = m.setupEntries(pID, repo).(updateEntries)
	m.paginator.SetTotalPages(len(m.entries))
	m.setViewportContent()
	return m
}

func (m *EntryModel) setupEntries(pID uint, repo e.Repository) tea.Msg {
	entries, err := repo.GetAllByProjectID(pID)
	if err != nil {
		return fmt.Errorf("cannot get entries with project id %d: %v", pID, err)
	}

	utils.Reverse(entries)
	return updateEntries(entries)
}

func (m *EntryModel) setViewportContent() {
	var content string
	if len(m.entries) == 0 {
		content = "There are no entries for this project :)"
	} else {
		content = e.FormattedEntry(m.entries[m.paginator.Page])
	}

	str, err := glamour.Render(content, "dark")
	if err != nil {
		m.error = "could not render content with glamour"
	}

	m.viewport.SetContent(str)
}

func (m *EntryModel) helpView() string {
	return HelpStyle("\n ↑/↓: navigate  • esc: back • c: create entry • d: delete entry • q: quit\n")
}

func (m EntryModel) Init() tea.Cmd {
	return nil
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		top, right, bottom, left := DocStyle.GetMargin()
		m.viewport = viewport.New(WindowSize.Width-left-right, WindowSize.Height-top-bottom-6)

	case errMsg:
		m.error = msg.Error()

	case editorFinishMsg:
		m.quitting = false
		if msg.err != nil {
			return m, tea.Quit
		}

		cmds = append(cmds, createEntryCmd(m.activePID, msg.file, AppStore.Entry))

	case updateEntries:
		m.entries = []*e.Entry(msg)
		m.paginator.SetTotalPages(len(m.entries))
		m.setViewportContent()

	case tea.KeyMsg:

		switch {

		case key.Matches(msg, Keymap.Create):
			m.quitting = false
			return m, openEditorCmd()

		case key.Matches(msg, Keymap.Back):
			return InitProject()

		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit

		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.paginator, cmd = m.paginator.Update(msg)
	cmds = append(cmds, cmd)

	m.setViewportContent() // refresh the content on every Update call
	return m, tea.Batch(cmds...)
}

func (m EntryModel) View() string {
	if m.quitting {
		return ""
	}

	formatted := lipgloss.JoinVertical(lipgloss.Left, "\n", m.viewport.View(), m.helpView(), ErrStyle(m.error), m.paginator.View())
	return DocStyle.Render(formatted)
}
