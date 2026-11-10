package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/momin-chezgi/task-tracker/internal/task"
	"github.com/momin-chezgi/task-tracker/internal/user"
)

// jsonFile overrides the default location in tests.
var jsonFile string

func DataFilePath() (string, error) {
	if jsonFile != "" {
		return jsonFile, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "task-tracker", "data.json"), nil
}

func Store(strg Storage) error {
	path, err := DataFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(strg.Tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func Load() (Storage, error) {
	path, err := DataFilePath()
	if err != nil {
		return Storage{Tasks: []task.Task{}, Users: []user.User{}}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Storage{Tasks: []task.Task{}, Users: []user.User{}}, nil
		}
		return Storage{}, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return Storage{Tasks: []task.Task{}, Users: []user.User{}}, nil
	}
	strg := Storage{}
	err = json.Unmarshal(data, &strg)
	return strg, err
}
