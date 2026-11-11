package store

import (
	"strings"

	"github.com/momin-chezgi/task-tracker/internal/user"
)

func (s *Storage) AddUsr(name string, email user.Email, pass user.Password) (int, error) {
	if strings.TrimSpace(name) == "" {
		return -1, EmptyTitleErr
	}
	if !pass.Validate() {
		return -1, InvalidPassword
	}
	if !email.Validate() {
		return -1, InvalidEmailAddr
	}
	newUser := user.User{
		ID:       user.GenId(),
		Name:     name,
		Email:    email,
		Password: pass,
	}
	s.Users = append(s.Users, newUser)
	return newUser.ID, nil
}

func (s *Storage) UpdateUsrName(usrID int, name string) error {
	if strings.TrimSpace(name) == "" {
		return EmptyTitleErr
	}
	if !s.ContainsUsr(usrID) {
		return InvalidIDErr
	}
	s.Users[usrID].Name = name
	return nil
}

func (s *Storage) UpdateUsrEmail(usrID int, email user.Email) error {
	if !email.Validate() {
		return InvalidEmailAddr
	}
	if !s.ContainsUsr(usrID) {
		return InvalidIDErr
	}
	s.Users[usrID].Email = email
	return nil
}

func (s *Storage) UpdateUsrPass(usrID int, pass user.Password) error {
	if !pass.Validate() {
		return InvalidPassword
	}
	if !s.ContainsUsr(usrID) {
		return InvalidIDErr
	}
	s.Users[usrID].Password = pass
	return nil
}
