package domain

// Project CRUD operation
type Repository interface {
	Print()
	IsEmpty() bool
	GetAll() ([]*Project, error)
	GetByID(id uint) (*Project, error)
	Create(project *Project) error
	DeleteByID(id uint) error
	UpdateByID(id uint, project *Project) error
}
