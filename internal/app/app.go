package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/momin-chezgi/task-tracker/internal/query"
	"github.com/momin-chezgi/task-tracker/internal/store"
	"github.com/momin-chezgi/task-tracker/internal/task"
)

func ProcessTheQueryAndGetTheResult(q *query.Query) (outErr error) {
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
		usrID, _ := strconv.Atoi(q.Arg2)
		tskID, err := storage.AddTsk(q.Arg1, usrID)
		if err != nil {
			return err
		}
		q.Output = []string{fmt.Sprintf("Task added successfully (ID: %v)", tskID)}
	case query.UpdateCommand:
		tskID, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.UpdateTsk(tskID, q.Arg2)
		return err
	case query.CancelCommand:
		tskID, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.CancelTsk(tskID)
		return err
	case query.MarkCommand:
		stt, err := task.StrToTskStatus(q.Arg1)
		if err != nil {
			return err
		}
		tskID, err := strconv.Atoi(q.Arg2)
		if err != nil {
			return err
		}
		err = storage.Mark(stt, tskID)
		return err
	case query.ListCommand:
		if q.Arg1 == "" {
			for _, t := range storage.Tasks {
				if t.TskStatus != task.Cancelled {
					q.Output = append(q.Output, t.Format())
				}
			}
		} else {
			stt, err := task.StrToTskStatus(q.Arg1)
			if err != nil {
				return err
			}
			for _, t := range storage.Tasks {
				if t.TskStatus == stt {
					q.Output = append(q.Output, t.Format())
				}
			}
		}
		if len(q.Output) == 0 {
			q.Output = append(q.Output, "Nothing here!")
		}
		return err
	case query.ShowCommand:
		tskID, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		if tskID < 0 || tskID >= len(storage.Tasks) {
			return store.InvalidIDErr
		}
		q.Output = []string{storage.Tasks[tskID].Format()}
	default:
		return errors.New("Invalid action")
	}
	return err
}
