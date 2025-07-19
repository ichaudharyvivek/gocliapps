package tasks

type Status string

const (
	Pending   Status = "Pending"
	Completed Status = "Completed"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

func NewTask(id int, desc string) *Task {
	return &Task{
		ID:          id,
		Description: desc,
		Status:      Pending,
	}
}
