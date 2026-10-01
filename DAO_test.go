package main

import (
	"testing"
	"time"
)

var storage = Storage{
	Tasks: []Task{
		{
			ID:          1,
			Description: "Coding the project",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Status:      InProgress,
		},
		{
			ID:          2,
			Description: "Contact to Saba steel company",
			CreatedAt:   time.Now().AddDate(0, 0, -28),
			UpdatedAt:   time.Now(),
			Status:      ToDo,
		},
	},
}

func TestAdd(t *testing.T) {
	t.Run("Empty_description", func(t *testing.T) {
		storage := Storage{}
		_, err := storage.Add("")
		if err != emptyDescriptionErr {
			t.Errorf("An empty string can't be the description of a task!")
		}
	})

	t.Run("Task_with_an_ordinary_description", func(t *testing.T) {
		newID, err := storage.Add("Arrange a meeting with my friend")
		if err != nil {
			t.Errorf("An error happened in storage.Add(): %v", err)
		}
		if newID != len(storage.Tasks) {
			t.Errorf("The new ID is not correct. Got:%d, want:%d", newID, len(storage.Tasks))
		}
	})
}

func TestUpdate(t *testing.T) {

	for _, id := range []int{-1, 0, 10000000, -2938973132} {
		t.Run("Out_of_range_id", func(t *testing.T) {
			err := storage.Update(id, "This ID is out of range and nothing should be changed!")
			if err != invalidIDErr {
				t.Errorf("An out of range number can't be the ID of a task!")
			}
		})
	}

	t.Run("Empty_description", func(t *testing.T) {
		err := storage.Update(1, "")
		if err != emptyDescriptionErr {
			t.Errorf("An empty string can't be the description of a task!")
		}
	})

	for id := 0; id < len(storage.Tasks); id++ {
		lastUpdate := storage.Tasks[id].UpdatedAt
		err := storage.Update(id+1, "This is an updated task")
		if err != nil {
			t.Errorf("Can't update successfuly: %v", err)
		}
		if storage.Tasks[id].Description != "This is an updated task" {
			t.Errorf("Can't update the description. Got: %v, want: This is an updated task", storage.Tasks[id].Description)
		}
		if storage.Tasks[id].UpdatedAt.Equal(lastUpdate) {
			t.Errorf("The updated time of the Task has not been updated")
		}
	}
}

func TestDelete(t *testing.T)
func TestMark(t *testing.T)
func TestFormat(t *testing.T)
