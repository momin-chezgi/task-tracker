package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/momin-chezgi/task-tracker/internal/user"

	"github.com/momin-chezgi/task-tracker/internal/task"
)

type Storage struct {
	Tasks []task.Task
	Users []user.User
}

var (
	EmptyTitleErr error = errors.New("An empty string can't be the name of a task!")
	InvalidIDErr  error = errors.New("The ID of the task is out of range!")
)

func (s Storage) containsUsr(usrID int) bool {
	if user.UserIDCounter <= usrID || usrID <= 0 {
		return false
	}
	return true
}

func (s Storage) containsTsk(tskID int) bool {
	if task.TskIdCounter <= tskID || tskID <= 0 {
		return false
	}
	return true
}

func (s *Storage) AddTsk(dsc string, usrId int) (int, error) {
	if dsc == "" {
		return -1, EmptyTitleErr
	}
	if !s.containsUsr(usrId) {
		return -1, InvalidIDErr
	}
	newTask := task.Task{
		ID:        task.GenId(),
		Title:     dsc,
		UserID:    user.GenId(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TskStatus: task.ToDo,
	}
	s.Tasks = append(s.Tasks, newTask)
	return newTask.ID, nil
}

func (s *Storage) UpdateTsk(tskID int, dsc string) error {
	if dsc == "" {
		return EmptyTitleErr
	}
	if !s.containsTsk(tskID) {
		return InvalidIDErr
	}
	s.Tasks[tskID].Title = dsc
	s.Tasks[tskID].UpdatedAt = time.Now()
	return nil
}

func (s *Storage) CancelTsk(tskID int) error {
	if !s.containsTsk(tskID) {
		return InvalidIDErr
	}

	s.Tasks[tskID].TskStatus = task.Cancelled
	return nil
}

func (s *Storage) Mark(stt task.TskStatus, tskID int) error {
	if !s.containsTsk(tskID) {
		return InvalidIDErr
	}
	if s.Tasks[tskID].TskStatus == task.Cancelled {
		return errors.New("This task has been cancelled!")
	}
	s.Tasks[tskID].TskStatus = stt
	return nil
}

// jsonFile overrides the default location in tests.
var jsonFile string

func DataFilePath() (string, error) {
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
	path, err := DataFilePath()
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
	path, err := DataFilePath()
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
