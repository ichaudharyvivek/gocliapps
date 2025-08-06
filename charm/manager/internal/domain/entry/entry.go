package domain

import "gorm.io/gorm"

// Entry is the log under a specific project
type Entry struct {
	gorm.Model
	ProjectID uint `gorm:"foreignKey:Project"`
	Message   string
}

func NewEntry(pID uint, message string) *Entry {
	return &Entry{
		ProjectID: pID,
		Message:   message,
	}
}
