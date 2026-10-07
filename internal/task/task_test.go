package task

import (
	"fmt"
	"testing"
	"time"
)

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
