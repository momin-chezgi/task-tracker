package main

import (
	"fmt"
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
