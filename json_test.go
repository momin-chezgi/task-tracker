package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestStoreAndLoad(t *testing.T) {
	storageTestCases := []Storage{
		{
			Tasks: []Task{
				{
					ID:          1,
					Description: "Nothin here!",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Status:      InProgress,
				},
				{
					ID:          2,
					Description: "Get the driving license",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now().Local().AddDate(1, 1, 0),
					Status:      Done,
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
