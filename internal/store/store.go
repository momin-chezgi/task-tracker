package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/momin-chezgi/task-tracker/internal/task"
)

type Storage struct {
	Tasks []task.Task
}

var (
	EmptyDescriptionErr error = errors.New("An empty string can't be the name of a task!")
	InvalidIDErr        error = errors.New("The ID of the task is out of range!")
)

func (s *Storage) Add(dsc string) (int, error) {
	if dsc == "" {
		return -1, EmptyDescriptionErr
	}
	newTask := task.Task{
		ID:          len(s.Tasks) + 1,
		Description: dsc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status:      task.ToDo,
	}
	s.Tasks = append(s.Tasks, newTask)
	return newTask.ID, nil
}

func (s *Storage) Update(id int, dsc string) error {
	if dsc == "" {
		return EmptyDescriptionErr
	}
	if id <= 0 || id > len(s.Tasks) {
		return InvalidIDErr
	}
	s.Tasks[id-1].Description = dsc
	s.Tasks[id-1].UpdatedAt = time.Now()
	return nil
}

func (s *Storage) Cancel(id int) error {
	if id <= 0 || id > len(s.Tasks) {
		return InvalidIDErr
	}

	s.Tasks[id-1].Status = task.Cancelled
	return nil
}

func (s *Storage) Mark(stt task.Status, id int) error {
	if id <= 0 || id > len(s.Tasks) {
		return InvalidIDErr
	}
	if s.Tasks[id-1].Status == task.Cancelled {
		return errors.New("This task has been cancelled!")
	}
	s.Tasks[id-1].Status = stt
	return nil
}

// jsonFile overrides the default location in tests.
var jsonFile string

func dataFilePath() (string, error) {
	if jsonFile != "" {
		return jsonFile, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "task-tracker", "data.json"), nil
}

func Store(strg Storage) error {
	path, err := dataFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(strg.Tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func Load() (Storage, error) {
	path, err := dataFilePath()
	if err != nil {
		return Storage{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Storage{Tasks: []task.Task{}}, nil
		}
		return Storage{}, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return Storage{Tasks: []task.Task{}}, nil
	}
	var tasks []task.Task
	err = json.Unmarshal(data, &tasks)
	return Storage{Tasks: tasks}, err
}
