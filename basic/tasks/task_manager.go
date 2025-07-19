package tasks

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrEmptyDescription = errors.New("task description cannot be empty")
	ErrTaskNotFound     = errors.New("task not found")
)

type Manager interface {
	GetAll() []*Task
	GetById(id int) *Task
	Create(desc string) (*Task, error)
	CreateAll(tasks []*Task) error
	Delete(id int) error
	MarkComplete(id int) error
}

type TasksManager struct {
	tasks  []*Task
	nextID int
	mu     sync.RWMutex
}

func NewTaskManager() *TasksManager {
	return &TasksManager{
		tasks:  make([]*Task, 0),
		nextID: 1,
	}
}

func (tm *TasksManager) GetAll() []*Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tasks := make([]*Task, len(tm.tasks))
	copy(tasks, tm.tasks)
	return tasks
}

func (tm *TasksManager) GetById(id int) *Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	for _, task := range tm.tasks {
		if task.ID == id {
			return task
		}
	}

	return nil
}

func (tm *TasksManager) Create(desc string) (*Task, error) {
	if desc == "" {
		return nil, ErrEmptyDescription
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	task := NewTask(tm.nextID, desc)
	tm.tasks = append(tm.tasks, task)
	tm.nextID += 1

	return task, nil
}

func (tm *TasksManager) CreateAll(tasks []*Task) error {
	if len(tasks) == 0 {
		return nil
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.tasks = append(tm.tasks, tasks...)
	maxID := 0
	for _, task := range tm.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	tm.nextID = maxID + 1
	return nil
}

func (tm *TasksManager) Delete(id int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for i, t := range tm.tasks {
		if t.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("tasks with id '%d' not found", id)
}

func (tm *TasksManager) MarkComplete(id int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, t := range tm.tasks {
		if t.ID == id {
			t.Status = Completed
			return nil
		}
	}

	return fmt.Errorf("tasks with id '%d' not found", id)
}
