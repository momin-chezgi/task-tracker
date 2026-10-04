package main

import (
	"fmt"
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

var outOfRangeIDs = []int{-1, 0, 10000000, -2938973132}

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

	for _, id := range outOfRangeIDs {
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

func TestDelete(t *testing.T) {
	for _, id := range outOfRangeIDs {
		t.Run(fmt.Sprintf("Out_range_id:%v", id), func(t *testing.T) {
			if err := storage.Delete(id); err != invalidIDErr {
				t.Errorf("An out of range number can't be the ID of a task!")
			}
		})
	}
	t.Run("Delete_the_first_task", func(t *testing.T) {
		err := storage.Delete(1)
		if err != nil {
			t.Fatalf("An error occurred while deleting a task: %v", err)
		}
		if storage.Tasks[0].Status != Cancelled {
			t.Errorf("The task #1 was not deleted!")
		}
	})
}

func TestMark(t *testing.T) {
	for _, id := range outOfRangeIDs {
		t.Run(fmt.Sprintf("Out_of_range_id:%v", id), func(t *testing.T) {
			if err := storage.Mark(Done, id); err != invalidIDErr {
				t.Errorf("An out of range number can't be the ID of a task!")
			}
		})
	}

	t.Run("Mark_as_cancelled:#1", func(t *testing.T) {
		if err := storage.Mark(Cancelled, 1); err != nil {
			t.Errorf("An error occurred while marking task #1 to the Cancelled flag")
		}
		if storage.Tasks[0].Status != Cancelled {
			t.Errorf("The task wasn't marked successfully; Got: %v, want: Cancelled", storage.Tasks[0].Status.String())
		}
	})
	t.Run("Mark_as_ToDo:#1", func(t *testing.T) {
		if err := storage.Mark(ToDo, 1); err != nil {
			t.Errorf("An error occurred while marking task #1 to the ToDo flag")
		}
		if storage.Tasks[0].Status != ToDo {
			t.Errorf("The task wasn't marked successfully; Got: %v, want: ToDo", storage.Tasks[0].Status.String())
		}
	})
	t.Run("Mark_as_inProgress:#1", func(t *testing.T) {
		if err := storage.Mark(InProgress, 1); err != nil {
			t.Errorf("An error occurred while marking task #1 to the InProgres flag")
		}
		if storage.Tasks[0].Status != InProgress {
			t.Errorf("The task wasn't marked successfully; Got: %v, want: InProgress", storage.Tasks[0].Status.String())
		}
	})
	t.Run("Mark_as_done:#1", func(t *testing.T) {
		if err := storage.Mark(Done, 1); err != nil {
			t.Errorf("An error occurred while marking task #1 to the Done flag")
		}
		if storage.Tasks[0].Status != Done {
			t.Errorf("The task wasn't marked successfully; Got: %v, want: Done", storage.Tasks[0].Status.String())
		}
	})

}

func TestFormat(t *testing.T) {
	t.Run("Task #1", func(t *testing.T) {
		tm := time.Now().Round(0)
		task := Task{
			ID:          100,
			Description: "Come up with a jogging routine",
			CreatedAt:   tm,
			UpdatedAt:   tm,
			Status:      Cancelled,
		}
		want := fmt.Sprintf("---------------------\nCome up with a jogging routine	(Cancelled)\nID: 100\nCreated at: %v\nUpdated at: %v\n---------------------\n", tm, tm)
		if task.Format() != want {
			t.Errorf("The task formatting is not in the correct way. Got:%s, want:%s", task.Format(), want)
		}
	})

}
