package main

import (
	"log"
	"manager/internal/db"
	project "manager/internal/domain/project"
	"manager/internal/tui"
)

func main() {
	dbz, err := db.NewConnection()
	if err != nil {
		log.Fatalf("cannot create db connection: %v", err)
	}

	store := db.NewRepository(dbz)
	projects, err := store.Project.GetAll()
	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}

	if len(projects) < 1 {
		name := project.NewProjectPrompt()
		project := project.NewProject(name)

		err := store.Project.Create(project)
		if err != nil {
			log.Fatalf("error creating project: %v", err)
		}
	}

	tui.StartTea(store)
}
