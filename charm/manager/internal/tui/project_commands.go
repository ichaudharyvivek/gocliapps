package tui

import (
	p "manager/internal/domain/project"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Returns a command that creates a new project and updates the list.
func createProjectCmd(name string, repo p.Repository) tea.Cmd {
	return func() tea.Msg {
		project := p.NewProject(name)
		if err := repo.Create(project); err != nil {
			return errMsg{err}
		}

		return updateProjectListMsg{}
	}
}

// Returns a command that renames a project with id and updates the list
func renameProjectCmd(id uint, name string, repo p.Repository) tea.Cmd {
	return func() tea.Msg {
		found, err := repo.GetByID(id)
		if err != nil {
			return errMsg{err}
		}

		found.Name = name
		if err = repo.UpdateByID(id, found); err != nil {
			return errMsg{err}
		}

		projects, _ := repo.GetAll()

		items := make([]list.Item, len(projects))
		for i, project := range projects {
			items[i] = list.Item(project)
		}

		return renameProjectMsg(items)
	}
}

// Returns a command that deletes a project with id and updates the list
func deleteProjectCmd(id uint, repo p.Repository) tea.Cmd {
	return func() tea.Msg {
		err := repo.DeleteByID(id)
		if err != nil {
			return errMsg{err}
		}

		return updateProjectListMsg{}
	}
}
