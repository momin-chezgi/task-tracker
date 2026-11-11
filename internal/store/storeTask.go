package store

import (
	"errors"
	"time"

	"github.com/momin-chezgi/task-tracker/internal/user"

	"github.com/momin-chezgi/task-tracker/internal/task"
)

var (
	EmptyTitleErr    error = errors.New("An empty string can't be the name of a task!")
	InvalidIDErr     error = errors.New("The ID of the task is out of range!")
	InvalidEmailAddr error = errors.New("The email address is not correct!")
	InvalidPassword  error = errors.New("The password is not valid")
)

func (s *Storage) AddTsk(dsc string, usrId int) (int, error) {
	if dsc == "" {
		return -1, EmptyTitleErr
	}
	if !s.ContainsUsr(usrId) {
		return -1, InvalidIDErr
	}
	newTask := task.Task{
		ID:        task.GenId(),
		Title:     dsc,
		UserID:    user.GenId(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TskStatus: task.ToDo,
	}
	s.Tasks = append(s.Tasks, newTask)
	return newTask.ID, nil
}

func (s *Storage) UpdateTsk(tskID int, dsc string) error {
	if dsc == "" {
		return EmptyTitleErr
	}
	if !s.ContainsTsk(tskID) {
		return InvalidIDErr
	}
	s.Tasks[tskID].Title = dsc
	s.Tasks[tskID].UpdatedAt = time.Now()
	return nil
}

func (s *Storage) CancelTsk(tskID int) error {
	if !s.ContainsTsk(tskID) {
		return InvalidIDErr
	}

	s.Tasks[tskID].TskStatus = task.Cancelled
	return nil
}

func (s *Storage) Mark(stt task.TskStatus, tskID int) error {
	if !s.ContainsTsk(tskID) {
		return InvalidIDErr
	}
	if s.Tasks[tskID].TskStatus == task.Cancelled {
		return errors.New("This task has been cancelled!")
	}
	s.Tasks[tskID].TskStatus = stt
	return nil
}
