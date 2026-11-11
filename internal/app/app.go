package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/momin-chezgi/task-tracker/internal/query"
	"github.com/momin-chezgi/task-tracker/internal/store"
	"github.com/momin-chezgi/task-tracker/internal/task"
	"github.com/momin-chezgi/task-tracker/internal/user"
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
	case query.AddTsk:
		usrID, _ := strconv.Atoi(q.Arg2)
		tskID, err := storage.AddTsk(q.Arg1, usrID)
		if err != nil {
			return err
		}
		q.Output = []string{fmt.Sprintf("Task added successfully (ID: %v)", tskID)}
	case query.UpdateTsk:
		tskID, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.UpdateTsk(tskID, q.Arg2)
		return err
	case query.CancelTsk:
		tskID, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.CancelTsk(tskID)
		return err
	case query.MarkTsk:
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
	case query.ListTsk:
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
	case query.ShowTsk:
		tskID, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		if !storage.ContainsTsk(tskID) {
			return store.InvalidIDErr
		}
		q.Output = []string{storage.Tasks[tskID].Format()}
	case query.AddUsr:
		usrID, err := storage.AddUsr(q.Arg1, user.Email(q.Arg2), user.Password(q.Arg3))
		if err != nil {
			return err
		}
		q.Output = []string{fmt.Sprintf("User added successfully (ID: %v)", usrID)}
	case query.UpdateUsrName:
		usrID, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.UpdateUsrName(usrID, q.Arg2)
		return err
	case query.UpdateUsrEmail:
		usrID, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.UpdateUsrEmail(usrID, user.Email(q.Arg2))
		return err
	case query.UpdateUsrPass:
		usrID, err := strconv.Atoi(q.Arg1)
		if err != nil {
			return err
		}
		err = storage.UpdateUsrPass(usrID, user.Password(q.Arg2))
		return err
	case query.ListUsr:
		for _, u := range storage.Users {
			q.Output = append(q.Output, u.Format())
		}
	case query.ShowUsr:

	default:
		return errors.New("Invalid action")
	}
	return err
}
