package task

import (
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"Title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TskStatus TskStatus `json:"TskStatus"`
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
	return fmt.Sprintf(taskFormattingTemplate, t.Title, t.TskStatus, t.ID, t.CreatedAt, t.UpdatedAt)
}

type TskStatus int

const (
	Cancelled TskStatus = iota - 1
	ToDo
	InProgress
	Done
)

func (stt TskStatus) String() string {
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

func StrToTskStatus(s string) (TskStatus, error) {
	switch s {
	case "mark-in-progress":
		return InProgress, nil
	case "in-progress":
		return InProgress, nil
	case "mark-todo":
		return ToDo, nil
	case "todo":
		return ToDo, nil
	case "mark-done":
		return Done, nil
	case "done":
		return Done, nil
	default:
		return ToDo, errors.New("Invalid TskStatus given")
	}
}
