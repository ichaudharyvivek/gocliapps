package tui

import (
	"fmt"
	"log"
	p "manager/internal/domain/project"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	mode                 int
	updateProjectListMsg struct{}
	renameProjectMsg     []list.Item
)

const (
	nav mode = iota
	edit
	create
)

// ProjectModel represent the home screen where the projects are listed
type ProjectModel struct {
	mode     mode
	list     list.Model
	input    textinput.Model
	quitting bool
}

func newProjectsList(repo p.Repository) ([]list.Item, error) {
	projects, err := repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	items := make([]list.Item, len(projects))
	for i, proj := range projects {
		items[i] = list.Item(proj)
	}

	return items, nil
}

func (m ProjectModel) getActiveProjectID() uint {
	active := m.list.Items()[m.list.Index()]
	return active.(*p.Project).ID
}

// Initialize the ProjectModel with default values
func InitProject() (tea.Model, tea.Cmd) {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Project name...."
	input.CharLimit = 250
	input.Width = 50

	items, err := newProjectsList(AppStore.Project)
	m := &ProjectModel{
		mode:  nav,
		input: input,
		list:  list.New(items, list.NewDefaultDelegate(), 8, 8),
	}

	if WindowSize.Height > 0 {
		top, right, bottom, left := DocStyle.GetMargin()
		m.list.SetSize(WindowSize.Width-left-right, WindowSize.Height-top-bottom-1)
	}

	m.list.Title = "Projects"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			Keymap.Create,
			Keymap.Enter,
			Keymap.Rename,
		}
	}

	return m, func() tea.Msg { return errMsg{err} }
}

func (m ProjectModel) Init() tea.Cmd {
	return nil
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		WindowSize = msg
		top, right, bottom, left := DocStyle.GetMargin()
		m.list.SetSize(WindowSize.Width-left-right, WindowSize.Height-top-bottom-1)

	case updateProjectListMsg:
		items, err := newProjectsList(AppStore.Project)
		if err != nil {
			log.Fatalf("cannot update the project list: %v", err)
		}

		m.mode = nav
		m.list.SetItems(items)

	case renameProjectMsg:
		m.mode = nav
		m.list.SetItems(msg)

	case tea.KeyMsg:
		if m.input.Focused() {
			if key.Matches(msg, Keymap.Enter) {
				if m.mode == create {
					cmds = append(cmds, createProjectCmd(m.input.Value(), AppStore.Project))
				}

				if m.mode == edit {
					cmds = append(cmds, renameProjectCmd(m.getActiveProjectID(), m.input.Value(), AppStore.Project))
				}

				m.mode = nav
				m.input.Blur()
				m.input.SetValue("")
			}

			if key.Matches(msg, Keymap.Back) {
				m.mode = nav
				m.input.Blur()
				m.input.SetValue("")
			}

			// Only log key press for the input when the field is focused
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {

			case key.Matches(msg, Keymap.Create):
				m.mode = create
				m.input.Focus()
				cmd = textinput.Blink

			case key.Matches(msg, Keymap.Quit):
				m.quitting = true
				return m, tea.Quit

			case key.Matches(msg, Keymap.Enter):
				activeProject := m.list.SelectedItem().(*p.Project)
				entry := InitEntry(activeProject.ID, AppStore.Entry)
				return entry.Update(WindowSize)

			case key.Matches(msg, Keymap.Rename):
				m.mode = edit
				m.input.Focus()
				cmd = textinput.Blink

			case key.Matches(msg, Keymap.Delete):
				items := m.list.Items()
				if len(items) > 0 {
					cmd = deleteProjectCmd(m.getActiveProjectID(), AppStore.Project)
				}

			default:
				m.list, cmd = m.list.Update(msg)
			}

			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m ProjectModel) View() string {
	if m.quitting {
		return ""
	}

	if m.input.Focused() {
		return DocStyle.Render(m.list.View() + "\n" + m.input.View())
	}

	return DocStyle.Render(m.list.View() + "\n")
}
