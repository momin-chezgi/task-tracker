package main

import (
	"errors"
	"fmt"
	"strconv"
)

type taskManager struct {
}

func (tm *taskManager) ProcessTheQueryAndGetTheResult(q *Query, buf []string) error {
	var storage Storage
	switch q.Action {
	case AddCommand:
		id, err := storage.Add(q.Arg1)
		if err != nil {
			return err
		}
		q.Output = []string{fmt.Sprintf("Task added successfully (ID: %v)", id)}
		return nil
	case UpdateCommand:
		id, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.Update(id, q.Arg2)
		if err != nil {
			return err
		}
	case DeleteCommand:
		id, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.Delete(id)
		if err != nil {
			return err
		}
	case MarkCommand:
		st, err := strToStatus(q.Arg1)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(q.Arg2)
		if err != nil {
			return err
		}
		err = storage.Mark(st, id)
	case ListCommand:
		if q.Arg1 == "" {
			for _, t := range storage.Tasks {
				// Design a way to handle the buffer considering the length of it, idk
			}
		}
	default:
		return errors.New("Invalid action")
	}
	return nil
}

func strToStatus(s string) (Status, error) {
	switch s {
	case "mark-in-progress":
		return InProgress, nil
	case "mark-todo":
		return ToDo, nil
	case "mark-done":
		return Done, nil
	default:
		return ToDo, errors.New("Invalid status given")
	}
}
