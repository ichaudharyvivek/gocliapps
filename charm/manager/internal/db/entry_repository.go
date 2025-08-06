package db

import (
	"fmt"
	e "manager/internal/domain/entry"

	"gorm.io/gorm"
)

type EntryRepository struct {
	DB *gorm.DB
}

func NewEntryRepository(db *gorm.DB) *EntryRepository {
	return &EntryRepository{DB: db}
}

func (r *EntryRepository) GetAllByProjectID(pID uint) ([]*e.Entry, error) {
	var entries []*e.Entry
	if err := r.DB.Where("project_id = ?", pID).Find(&entries).Error; err != nil {
		return entries, fmt.Errorf("entries not found with project id %d", pID)
	}

	return entries, nil
}

func (r *EntryRepository) Create(entry *e.Entry) error {
	if err := r.DB.Create(&entry).Error; err != nil {
		return fmt.Errorf("cannot create entry %w", err)
	}

	return nil
}

func (r *EntryRepository) DeleteByID(eID uint) error {
	if err := r.DB.Delete(&e.Entry{}, eID); err != nil {
		return fmt.Errorf("cannot delete entry with id %d", eID)
	}
	return nil
}

func (r *EntryRepository) DeleteAllByProjectID(pID uint) error {
	if err := r.DB.Where("project_id = ?", pID).Delete(&e.Entry{}); err != nil {
		return fmt.Errorf("cannot delete entries with project id %d", pID)
	}

	return nil
}
