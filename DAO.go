package main

import (
	"errors"
	"fmt"
	"time"
)

type Storage struct {
	Tasks           []Task
	CanceledTasks   []*Task
	TodoTasks       []*Task
	InprogressTasks []*Task
	DoneTasks       []*Task
}

type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      Status
}

type Status int

const (
	Canceled Status = iota + 1
	ToDo
	InProgress
	Done
)

var (
	emptyDescriptionErr error = errors.New("An empty string can't be the name of a task!")
	invalidIDErr        error = errors.New("The ID of the task is out of range!")
)

func (s *Storage) Add(dsc string) (int, error) {
	if dsc == "" {
		return -1, emptyDescriptionErr
	}
	newTask := Task{
		ID:          len(s.Tasks),
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
	if id < 0 || id >= len(s.Tasks) {
		return invalidIDErr
	}
	s.Tasks[id].Description = dsc
	return nil
}

func (s *Storage) Delete(id int) error {
	if id < 0 || id >= len(s.Tasks) {
		return invalidIDErr
	}
	s.Tasks = append(s.Tasks[:id], s.Tasks[id+1:]...)
	return nil
}

func (s *Storage) Mark(st Status, id int) error {
	if id < 0 || id >= len(s.Tasks) {
		return invalidIDErr
	}
	s.Tasks[id].Status = st
	return nil
}

var (
	taskFormattingTemplate string = "%s (%v)\nID: %v\nCreated at: %v\nUpdated at: %v\n"
	// For instance "Do the dishes(in-progress)\nID: 12\nCreated at: Sep 28 15:15\nUpdated at: sep 28 20:20\n"
)

func (t Task) Format() string {
	return fmt.Sprintf(taskFormattingTemplate, t.Description, t.ID, t.CreatedAt, t.UpdatedAt)
}
