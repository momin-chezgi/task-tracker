package main

import (
	"errors"
	"fmt"
	"time"
)

type Storage struct {
	Tasks []Task
}

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      Status    `json:"status"`
}

type Status int

const (
	Cancelled Status = iota + 1
	ToDo
	InProgress
	Done
)

func (stt Status) String() string {
	switch stt {
	case Cancelled:
		return "Cancelled"
	case ToDo:
		return "To-Do"
	case InProgress:
		return "In progress"
	case Done:
		return "Done"
	}
	return "None of them"
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

func (s *Storage) Delete(id int) error {
	if id <= 0 || id >= len(s.Tasks) {
		return invalidIDErr
	}

	s.Tasks[id-1].Status = Cancelled
	return nil
}

func (s *Storage) Mark(stt Status, id int) error {
	if id <= 0 || id > len(s.Tasks) {
		return invalidIDErr
	}
	s.Tasks[id-1].Status = stt
	return nil
}

var (
	taskFormattingTemplate string = "---------------------\n%s	(%v)\nID: %v\nCreated at: %v\nUpdated at: %v\n---------------------\n"
	// For instance:
	// ---------------------
	// Do the dishes	(in-progress)
	// ID: 12
	// Created at: Sep 28 15:15
	// Updated at: sep 28 20:20
	// ---------------------
)

func (t Task) Format() string {
	return fmt.Sprintf(taskFormattingTemplate, t.Description, t.Status, t.ID, t.CreatedAt, t.UpdatedAt)
}
