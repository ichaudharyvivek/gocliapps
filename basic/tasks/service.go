package tasks

import "fmt"

type TaskService struct {
	store   Store
	manager Manager
}

func NewTaskService(filepath string) *TaskService {
	store := NewFileStore(filepath)
	manager := NewTaskManager()

	// Load existing tasks from a file if any
	tasks, err := store.Load()
	if err != nil {
		fmt.Printf("Warning: Could not load existing tasks: %v\n", err)
	}

	if err := manager.CreateAll(tasks); err != nil {
		fmt.Printf("Warning: Could not load existing tasks: %v\n", err)
	}

	return &TaskService{
		store:   store,
		manager: manager,
	}
}

func (ts *TaskService) AddTask(description string) (*Task, error) {
	task, err := ts.manager.Create(description)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	tasks := ts.manager.GetAll()
	if err := ts.store.Save(tasks); err != nil {
		return nil, fmt.Errorf("failed to save tasks: %w", err)
	}

	return task, nil
}

func (ts *TaskService) ListTasks() ([]*Task, error) {
	tasks := ts.manager.GetAll()
	if err := ts.store.Save(tasks); err != nil {
		return nil, fmt.Errorf("failed to save tasks: %w", err)
	}

	return tasks, nil
}

func (ts *TaskService) CompleteTask(id int) error {
	if err := ts.manager.MarkComplete(id); err != nil {
		return fmt.Errorf("failed to mark task complete: %w", err)
	}

	tasks := ts.manager.GetAll()
	if err := ts.store.Save(tasks); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	return nil
}

func (ts *TaskService) DeleteTask(id int) error {
	if err := ts.manager.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	tasks := ts.manager.GetAll()
	if err := ts.store.Save(tasks); err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}

	return nil
}
