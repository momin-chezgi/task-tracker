package store

import (
	"github.com/momin-chezgi/task-tracker/internal/task"
	"github.com/momin-chezgi/task-tracker/internal/user"
)

type Storage struct {
	Tasks []task.Task
	Users []user.User
}

func (s Storage) ContainsUsr(usrID int) bool {
	if user.UserIDCounter <= usrID || usrID <= 0 {
		return false
	}
	return true
}

func (s Storage) ContainsTsk(tskID int) bool {
	if task.TskIdCounter <= tskID || tskID <= 0 {
		return false
	}
	return true
}
