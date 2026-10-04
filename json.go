package main

import (
	"bytes"
	"encoding/json"
	"os"
)

var (
	jsonFile = "data.json"
)

func Store(strg Storage) error {
	data, err := json.MarshalIndent(strg.Tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(jsonFile, data, 0664)
}

func Load() (Storage, error) {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return Storage{}, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return Storage{Tasks: []Task{}}, nil
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return Storage{Tasks: tasks}, err
}
