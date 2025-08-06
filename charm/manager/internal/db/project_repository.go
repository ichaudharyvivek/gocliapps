package db

import (
	"fmt"
	"log"
	p "manager/internal/domain/project"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	DB *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

func (r *ProjectRepository) Print() {
	projects, err := r.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, project := range projects {
		fmt.Printf("%d: %s", project.ID, project.Name)
	}
}

func (r *ProjectRepository) IsEmpty() bool {
	if projects, _ := r.GetAll(); len(projects) == 0 {
		return true
	}

	return false
}

func (r *ProjectRepository) GetAll() ([]*p.Project, error) {
	var projects []*p.Project
	if err := r.DB.Find(&projects).Error; err != nil {
		return projects, fmt.Errorf("table is empty: %w", err)
	}

	return projects, nil
}

func (r *ProjectRepository) GetByID(id uint) (*p.Project, error) {
	var project *p.Project
	err := r.DB.Where("id = ?", id).Find(&project).Error
	if err != nil {
		return project, fmt.Errorf("cannot find project: %w", err)
	}

	return project, nil
}

func (r *ProjectRepository) Create(project *p.Project) error {
	if err := r.DB.Create(&project).Error; err != nil {
		return fmt.Errorf("cannot create project: %w", err)
	}

	return nil
}

func (r *ProjectRepository) DeleteByID(id uint) error {
	if err := r.DB.Delete(&p.Project{}, id).Error; err != nil {
		return fmt.Errorf("cannot delete project with id %d\n%w", id, err)
	}

	return nil
}

func (r *ProjectRepository) UpdateByID(id uint, project *p.Project) error {
	result := r.DB.Model(&p.Project{}).Where("id = ?", id).Updates(project)

	if result.Error != nil {
		return fmt.Errorf("failed to update project: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("project with id %d not found", id)
	}

	return nil
}
