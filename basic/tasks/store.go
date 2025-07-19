package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Store interface {
	Load() ([]*Task, error)
	Save(tasks []*Task) error
}

type FileStore struct {
	filepath string
}

func NewFileStore(filepath string) *FileStore {
	return &FileStore{
		filepath: filepath,
	}
}

func (fs *FileStore) Load() ([]*Task, error) {
	data, err := os.ReadFile(fs.filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*Task{}, nil
		}

		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	var tasks []*Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks file: %w", err)
	}

	return tasks, nil
}

func (fs *FileStore) Save(tasks []*Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if err := os.WriteFile(fs.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}

	return nil
}
