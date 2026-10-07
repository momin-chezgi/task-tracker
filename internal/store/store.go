package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Storage struct {
	Tasks []Task
}

var (
	emptyDescriptionErr error = errors.New("An empty string can't be the name of a task!")
	invalidIDErr        error = errors.New("The ID of the task is out of range!")
)

func (s *Storage) Add(dsc string) (int, error) {
	if dsc == "" {
		return -1, emptyDescriptionErr
	}
	newTask := Task{
		ID:          len(s.Tasks) + 1,
		Description: dsc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status:      ToDo,
	}
	s.Tasks = append(s.Tasks, newTask)
	return newTask.ID, nil
}

func (s *Storage) Update(id int, dsc string) error {
	if dsc == "" {
		return emptyDescriptionErr
	}
	if id <= 0 || id > len(s.Tasks) {
		return invalidIDErr
	}
	s.Tasks[id-1].Description = dsc
	s.Tasks[id-1].UpdatedAt = time.Now()
	return nil
}

func (s *Storage) Cancel(id int) error {
	if id <= 0 || id > len(s.Tasks) {
		return invalidIDErr
	}

	s.Tasks[id-1].Status = Cancelled
	return nil
}

func (s *Storage) Mark(stt Status, id int) error {
	if id <= 0 || id > len(s.Tasks) {
		return invalidIDErr
	}
	if s.Tasks[id-1].Status == Cancelled {
		return errors.New("This task has been cancelled!")
	}
	s.Tasks[id-1].Status = stt
	return nil
}

var (
	jsonFile = "data.json"
)

func Store(strg Storage) error {
	data, err := json.MarshalIndent(strg.Tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(jsonFile, data, 0664)
}

func Load() (Storage, error) {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return Storage{}, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return Storage{Tasks: []Task{}}, nil
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return Storage{Tasks: tasks}, err
}
