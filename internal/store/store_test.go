package store

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/momin-chezgi/task-tracker/internal/task"
)

var storage = Storage{
	Tasks: []task.Task{
		{
			ID:          1,
			Description: "Coding the project",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Status:      task.InProgress,
		},
		{
			ID:          2,
			Description: "Contact to Saba steel company",
			CreatedAt:   time.Now().AddDate(0, 0, -28),
			UpdatedAt:   time.Now(),
			Status:      task.ToDo,
		},
	},
}

var outOfRangeIDs = []int{-1, 0, 10000000, -2938973132}

func TestAdd(t *testing.T) {
	t.Run("Empty_description", func(t *testing.T) {
		storage := Storage{}
		_, err := storage.Add("")
		if err != EmptyDescriptionErr {
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
			if err != InvalidIDErr {
				t.Errorf("An out of range number can't be the ID of a task!")
			}
		})
	}

	t.Run("Empty_description", func(t *testing.T) {
		err := storage.Update(1, "")
		if err != EmptyDescriptionErr {
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

func TestCancel(t *testing.T) {
	for _, id := range outOfRangeIDs {
		t.Run(fmt.Sprintf("Out_range_id:%v", id), func(t *testing.T) {
			if err := storage.Cancel(id); err != InvalidIDErr {
				t.Errorf("An out of range number can't be the ID of a task!")
			}
		})
	}
	t.Run("Cancel_the_first_task", func(t *testing.T) {
		err := storage.Cancel(1)
		if err != nil {
			t.Fatalf("An error occurred while deleting a task: %v", err)
		}
		if storage.Tasks[0].Status != task.Cancelled {
			t.Errorf("The task #1 was not Canceld!")
		}
	})
}

func TestMark(t *testing.T) {
	for _, id := range outOfRangeIDs {
		t.Run(fmt.Sprintf("Out_of_range_id:%v", id), func(t *testing.T) {
			if err := storage.Mark(task.Done, id); err != InvalidIDErr {
				t.Errorf("An out of range number can't be the ID of a task!")
			}
		})
	}

	t.Run("Mark_as_ToDo:#1", func(t *testing.T) {
		if err := storage.Mark(task.ToDo, 1); err != nil {
			t.Errorf("An error occurred while marking task #1 to the ToDo flag")
		}
		if storage.Tasks[0].Status != task.ToDo {
			t.Errorf("The task wasn't marked successfully; Got: %v, want: ToDo", storage.Tasks[0].Status.String())
		}
	})
	t.Run("Mark_as_inProgress:#1", func(t *testing.T) {
		if err := storage.Mark(task.InProgress, 1); err != nil {
			t.Errorf("An error occurred while marking task #1 to the InProgres flag")
		}
		if storage.Tasks[0].Status != task.InProgress {
			t.Errorf("The task wasn't marked successfully; Got: %v, want: InProgress", storage.Tasks[0].Status.String())
		}
	})
	t.Run("Mark_as_done:#1", func(t *testing.T) {
		if err := storage.Mark(task.Done, 1); err != nil {
			t.Errorf("An error occurred while marking task #1 to the Done flag")
		}
		if storage.Tasks[0].Status != task.Done {
			t.Errorf("The task wasn't marked successfully; Got: %v, want: Done", storage.Tasks[0].Status.String())
		}
	})

	t.Run("Mark_as_cancelled:#1", func(t *testing.T) {
		if err := storage.Mark(task.Cancelled, 1); err != nil {
			t.Errorf("An error occurred while marking task #1 to the Cancelled flag")
		}
		if storage.Tasks[0].Status != task.Cancelled {
			t.Errorf("The task wasn't marked successfully; Got: %v, want: Cancelled", storage.Tasks[0].Status.String())
		}
	})

}

func TestStoreAndLoad(t *testing.T) {
	storageTestCases := []Storage{
		{
			Tasks: []task.Task{
				{
					ID:          1,
					Description: "Nothin here!",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Status:      task.InProgress,
				},
				{
					ID:          2,
					Description: "Get the driving license",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now().Local().AddDate(1, 1, 0),
					Status:      task.Done,
				},
			},
		},
	}
	for _, tt := range storageTestCases {
		t.Run("Test saving a storage", func(t *testing.T) {
			err := Store(tt)
			if err != nil {
				t.Fatalf("An error occurred while calling Store(): %v", err)
			}
			savedStorage, err := Load()
			if err != nil {
				t.Fatalf("An error occurred while calling Load(): %v", err)
			}

			for i := range tt.Tasks {
				tt.Tasks[i].CreatedAt = tt.Tasks[i].CreatedAt.Round(0)
				tt.Tasks[i].UpdatedAt = tt.Tasks[i].UpdatedAt.Round(0)
			}

			if !reflect.DeepEqual(savedStorage, tt) {
				t.Errorf("The saved tasks and the given ones are different!")
				fmt.Printf("Got: %v\n", savedStorage)
				fmt.Printf("Want: %v\n", tt)
			}
		})
	}
}

func TestLoadEmptyFile(t *testing.T) {
	originalJSONFile := jsonFile
	jsonFile = t.TempDir() + "/data.json"
	t.Cleanup(func() {
		jsonFile = originalJSONFile
	})

	if err := os.WriteFile(jsonFile, nil, 0664); err != nil {
		t.Fatalf("An error occurred while creating an empty data file: %v", err)
	}

	storage, err := Load()
	if err != nil {
		t.Fatalf("Load() should accept an empty data file: %v", err)
	}
	if len(storage.Tasks) != 0 {
		t.Errorf("Load() returned tasks for an empty data file: got %v", storage.Tasks)
	}
}
