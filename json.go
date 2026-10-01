package main

import (
	"encoding/json"
	"os"
)

var (
	jsonFile = "data.json"
)

func Store(strg Storage) error {
	data, err := json.Marshal(strg.Tasks)
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
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return Storage{Tasks: tasks}, err
}
