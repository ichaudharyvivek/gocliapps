package domain

// Repository the CRUD functionality for entries
type Repository interface {
	GetAllByProjectID(pID uint) ([]*Entry, error)
	Create(entry *Entry) error
	DeleteByID(eID uint) error
	DeleteAllByProjectID(pID uint) error
}
