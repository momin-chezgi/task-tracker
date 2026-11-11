package store

import (
	"testing"

	"github.com/momin-chezgi/task-tracker/internal/user"
)

func resetUserIDCounter(t *testing.T) {
	t.Helper()
	original := user.UserIDCounter
	user.UserIDCounter = 1
	t.Cleanup(func() { user.UserIDCounter = original })
}

func TestAddUsr(t *testing.T) {
	resetUserIDCounter(t)
	storage := Storage{Users: []user.User{{ID: 0, Name: "root"}}}

	id, err := storage.AddUsr("Ada", user.Email("ada@example.com"), user.Password("secret"))
	if err != nil {
		t.Fatalf("AddUsr() returned unexpected error: %v", err)
	}
	if id != 1 {
		t.Errorf("AddUsr() ID = %d, want 1", id)
	}
	if len(storage.Users) != 2 {
		t.Fatalf("AddUsr() left %d users, want 2", len(storage.Users))
	}
	got := storage.Users[1]
	if got.ID != id || got.Name != "Ada" || got.Email != "ada@example.com" || got.Password != "secret" {
		t.Errorf("AddUsr() stored %#v, want user ID %d with supplied name, email, and password", got, id)
	}
}

func TestAddUsrRejectsInvalidFields(t *testing.T) {
	tests := []struct {
		name     string
		username string
		email    user.Email
		password user.Password
		wantErr  error
	}{
		{"blank name", " \t", "ada@example.com", "secret", EmptyTitleErr},
		{"blank email", "Ada", " ", "secret", InvalidEmailAddr},
		{"blank password", "Ada", "ada@example.com", " ", InvalidPassword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := Storage{}
			if _, err := storage.AddUsr(tt.username, tt.email, tt.password); err != tt.wantErr {
				t.Errorf("AddUsr() error = %v, want %v", err, tt.wantErr)
			}
			if len(storage.Users) != 0 {
				t.Errorf("AddUsr() stored users after invalid input: %#v", storage.Users)
			}
		})
	}
}

func TestUpdateUsrFields(t *testing.T) {
	resetUserIDCounter(t)
	user.UserIDCounter = 2
	storage := Storage{Users: []user.User{
		{ID: 0, Name: "root", Email: "root@example.com", Password: "root"},
		{ID: 1, Name: "Ada", Email: "old@example.com", Password: "old"},
	}}

	tests := []struct {
		name string
		call func() error
		want string
	}{
		{"name", func() error { return storage.UpdateUsrName(1, "Ada Lovelace") }, "Ada Lovelace"},
		{"email", func() error { return storage.UpdateUsrEmail(1, "ada@example.com") }, "ada@example.com"},
		{"password", func() error { return storage.UpdateUsrPass(1, "new-secret") }, "new-secret"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); err != nil {
				t.Fatalf("update returned unexpected error: %v", err)
			}
			got := storage.Users[1]
			switch tt.name {
			case "name":
				if got.Name != tt.want {
					t.Errorf("updated name = %q, want %q", got.Name, tt.want)
				}
			case "email":
				if got.Email.String() != tt.want {
					t.Errorf("updated email = %q, want %q", got.Email, tt.want)
				}
			case "password":
				if got.Password.String() != tt.want {
					t.Errorf("updated password = %q, want %q", got.Password, tt.want)
				}
			}
		})
	}
}

func TestUpdateUsrRejectsInvalidValues(t *testing.T) {
	storage := Storage{Users: []user.User{{ID: 0}, {ID: 1, Name: "Ada", Email: "ada@example.com", Password: "secret"}}}
	if err := storage.UpdateUsrName(1, " \t"); err != EmptyTitleErr {
		t.Errorf("UpdateUsrName() error = %v, want %v", err, EmptyTitleErr)
	}
	if err := storage.UpdateUsrEmail(1, " "); err != InvalidEmailAddr {
		t.Errorf("UpdateUsrEmail() error = %v, want %v", err, InvalidEmailAddr)
	}
	if err := storage.UpdateUsrPass(1, " "); err != InvalidPassword {
		t.Errorf("UpdateUsrPass() error = %v, want %v", err, InvalidPassword)
	}
	if err := storage.UpdateUsrName(99, "Grace"); err != InvalidIDErr {
		t.Errorf("UpdateUsrName() with missing user error = %v, want %v", err, InvalidIDErr)
	}
}
