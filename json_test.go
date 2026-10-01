package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
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
			},
		},
	}
	for _, tt := range storageTestCases {
		t.Run("Test saving a storage", func(t *testing.T) {
			err := Store(tt)
			if err != nil {
				t.Fatalf("An error occurred while calling Store(): %w", err)
			}
			data, err := os.ReadFile(jsonFile)
			if err != nil {
				t.Fatalf("An error occurred while reading the file: %w", err)
			}
			var savedStorage Storage
			err = json.Unmarshal(data, &savedStorage)
			if err != nil {
				t.Fatalf("An error occurred in Unmarshal(): %w", err)
			}
			if !reflect.DeepEqual(tt, savedStorage) {
				t.Errorf("The saved and written structs are not equal!")
			}
		})
	}
}
