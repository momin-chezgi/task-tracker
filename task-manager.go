package main

import (
	"errors"
	"fmt"
	"strconv"
)

type taskManager struct {
}

func (tm *taskManager) ProcessTheQueryAndGetTheResult(q *Query) (outErr error) {
	storage, err := Load()
	if err != nil {
		return err
	}
	defer func() {
		if outErr == nil {
			outErr = Store(storage)
		}
	}()

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
		stt, err := strToStatus(q.Arg1)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(q.Arg2)
		if err != nil {
			return err
		}
		err = storage.Mark(stt, id)
	case ListCommand:
		if q.Arg1 == "" {
			for _, t := range storage.Tasks {
				if t.Status != Cancelled {
					q.Output = append(q.Output, t.Format())
				}
			}
		} else {
			stt, err := strToStatus(q.Arg1)
			if err != nil {
				return err
			}
			for _, t := range storage.Tasks {
				if t.Status == stt {
					q.Output = append(q.Output, t.Format())
				}
			}
		}
	default:
		return errors.New("Invalid action")
	}
	return err
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
