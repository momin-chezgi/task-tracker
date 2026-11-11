package user

import (
	"fmt"
	"strings"
)

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

var (
	userFormatingTemplate string = "---------------------\n%s	(%v)\nID: %v\n---------------------\n"
)

func (u User) Format() string {
	return fmt.Sprintf(userFormatingTemplate, u.Name, u.Email, u.ID)
}

type Email string

func (e Email) Validate() bool {
	if strings.TrimSpace(e.String()) == "" {
		return false
	}
	return true
}

func (e Email) String() string {
	return string(e)
}

type Password string

func (pw Password) Validate() bool {
	if strings.TrimSpace(pw.String()) == "" {
		return false
	}
	return true
}

func (pw Password) String() string {
	return string(pw)
}

var (
	UserIDCounter = 1
	rootUser      = User{ID: 0, Name: "root", Email: "root@root.com", Password: "root"}
)

func GenId() (id int) {
	id = UserIDCounter
	UserIDCounter++
	return id
}
