package main

import (
	"errors"
	"fmt"
	"time"
)

type Storage struct {
	Tasks []Task
	// CanceledTasks   []*Task
	TodoTasks       []*Task
	InprogressTasks []*Task
	DoneTasks       []*Task
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
	s.TodoTasks = append(s.TodoTasks, &s.Tasks[len(s.Tasks)-1])
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
	switch s.Tasks[id].Status {
	case ToDo:
		s.TodoTasks = append(s.TodoTasks[:id], s.TodoTasks[id+1:]...)
	case InProgress:
		s.InprogressTasks = append(s.InprogressTasks[:id], s.InprogressTasks[id+1:]...)
	case Done:
		s.DoneTasks = append(s.DoneTasks[:id], s.DoneTasks[id+1:]...)
	}

	s.Tasks = append(s.Tasks[:id], s.Tasks[id+1:]...)
	return nil
}

func (s *Storage) Mark(stt Status, id int) error {
	if id < 0 || id >= len(s.Tasks) {
		return invalidIDErr
	}
	prevStt := s.Tasks[id].Status
	switch prevStt {
	case ToDo:
		s.TodoTasks = append(s.TodoTasks[:id], s.TodoTasks[id+1:]...)
	case InProgress:
		s.InprogressTasks = append(s.InprogressTasks[:id], s.InprogressTasks[id+1:]...)
	case Done:
		s.DoneTasks = append(s.DoneTasks[:id], s.DoneTasks[id+1:]...)
	}
	s.Tasks[id].Status = stt
	switch stt {
	case ToDo:
		s.TodoTasks = append(s.TodoTasks, &s.Tasks[id])
	case InProgress:
		s.InprogressTasks = append(s.InprogressTasks, &s.Tasks[id])
	case Done:
		s.DoneTasks = append(s.DoneTasks, &s.Tasks[id])
	}

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
	return fmt.Sprintf(taskFormattingTemplate, t.Description, t.ID, t.CreatedAt, t.UpdatedAt)
}
