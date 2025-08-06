package db

import (
	e "manager/internal/domain/entry"
	p "manager/internal/domain/project"

	"gorm.io/gorm"
)

type Repositories struct {
	Project p.Repository
	Entry   e.Repository
}

func NewRepository(db *gorm.DB) *Repositories {
	return &Repositories{
		Project: NewProjectRepository(db),
		Entry:   NewEntryRepository(db),
	}
}
