package tui

import (
	"fmt"
	e "manager/internal/domain/entry"
	"manager/pkg/utils"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

const defaultEditor = "vim"

func openEditorCmd() tea.Cmd {
	file, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return func() tea.Msg {
			return errMsg{err}
		}
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}

	c := exec.Command(editor, file.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishMsg{err, file}
	})
}

func createEntryCmd(pID uint, file *os.File, repo e.Repository) tea.Cmd {
	return func() tea.Msg {
		input, err := utils.ReadFile(file)
		if err != nil {
			return errMsg{fmt.Errorf("cannot read file in createEntryCmd: %w", err)}
		}

		entry := e.NewEntry(pID, string(input))
		if err := repo.Create(entry); err != nil {
			return errMsg{fmt.Errorf("cannot create entry: %v", err)}
		}
		if err := os.Remove(file.Name()); err != nil {
			return errMsg{fmt.Errorf("cannot remove file: %v", err)}
		}
		if err := file.Close(); err != nil {
			return errMsg{fmt.Errorf("unable to close file: %v", err)}
		}

		// Return updated entries
		entries, err := repo.GetAllByProjectID(pID)
		if err != nil {
			return errMsg{fmt.Errorf("cannot get entries after creation: %v", err)}
		}
		utils.Reverse(entries)
		return updateEntries(entries)
	}
}
