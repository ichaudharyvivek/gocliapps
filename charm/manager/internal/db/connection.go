package db

import (
	"fmt"
	e "manager/internal/domain/entry"
	p "manager/internal/domain/project"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("manager.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	err = db.AutoMigrate(&e.Entry{}, &p.Project{})
	if err != nil {
		return db, fmt.Errorf("unable to migrate data: %w", err)
	}

	return db, nil
}
