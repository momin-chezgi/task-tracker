package main

import "encoding/json"

func Store(strg Storage) error {
	jsTasks, err := json.Marshal(strg.Tasks)

}

func Load() (Storage, error) {

}
