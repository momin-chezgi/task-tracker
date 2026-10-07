package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/momin-chezgi/task-tracker/internal/query"
	"github.com/momin-chezgi/task-tracker/internal/store"
	"github.com/momin-chezgi/task-tracker/internal/task"
)

type TaskManager struct {
}

func (tm *TaskManager) ProcessTheQueryAndGetTheResult(q *query.Query) (outErr error) {
	storage, err := store.Load()
	if err != nil {
		return err
	}
	defer func() {
		if outErr == nil {
			outErr = store.Store(storage)
		}
	}()

	switch q.Action {
	case query.AddCommand:
		id, err := storage.Add(q.Arg1)
		if err != nil {
			return err
		}
		q.Output = []string{fmt.Sprintf("Task added successfully (ID: %v)", id)}
	case query.UpdateCommand:
		id, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.Update(id, q.Arg2)
		return err
	case query.CancelCommand:
		id, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.Cancel(id)
		return err
	case query.MarkCommand:
		stt, err := task.StrToStatus(q.Arg1)
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(q.Arg2)
		if err != nil {
			return err
		}
		err = storage.Mark(stt, id)
		return err
	case query.ListCommand:
		if q.Arg1 == "" {
			for _, t := range storage.Tasks {
				if t.Status != task.Cancelled {
					q.Output = append(q.Output, t.Format())
				}
			}
		} else {
			stt, err := task.StrToStatus(q.Arg1)
			if err != nil {
				return err
			}
			for _, t := range storage.Tasks {
				if t.Status == stt {
					q.Output = append(q.Output, t.Format())
				}
			}
		}
		if len(q.Output) == 0 {
			q.Output = append(q.Output, "Nothing here!")
		}
		return err
	case query.ShowCommand:
		id, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		if id <= 0 || id > len(storage.Tasks) {
			return store.InvalidIDErr
		}
		q.Output = []string{storage.Tasks[id-1].Format()}
	default:
		return errors.New("Invalid action")
	}
	return err
}
