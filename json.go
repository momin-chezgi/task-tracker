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
	err = os.WriteFile(jsonFile, data, 0664)
	if err != nil {
		return err
	}
	return nil
}

func Load() (Storage, error) {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return Storage{}, err
	}
	var storage Storage
	err = json.Unmarshal(data, &storage)
	if err != nil {
		return storage, err
	}
	return storage, nil
}
