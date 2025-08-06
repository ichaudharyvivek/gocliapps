package domain

import (
	"bufio"
	"fmt"
	"os"

	"gorm.io/gorm"
)

// Project is the main entity of the program. It holds entries.
// Implements the list.Item for bubbletea TUI
type Project struct {
	gorm.Model
	Name string
}

func NewProject(name string) *Project {
	return &Project{
		Name: name,
	}
}

func NewProjectPrompt() string {
	fmt.Println("Welcome!")
	fmt.Println("Let's create a new project.")
	fmt.Print("> ")

	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()

	name := sc.Text()
	return name
}

// Title the project title to display in a list
func (p Project) Title() string {
	return p.Name
}

// Description the project description to display in a list
func (p Project) Description() string {
	return fmt.Sprintf("%d", p.ID)
}

// FilterValue choose what field to use for filtering in a Bubbletea list component
func (p Project) FilterValue() string {
	return p.Name
}
