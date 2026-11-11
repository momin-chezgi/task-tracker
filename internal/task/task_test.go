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
			ID:        100,
			Title:     "Come up with a jogging routine",
			CreatedAt: tm,
			UpdatedAt: tm,
			TskStatus: Cancelled,
		}
		want := fmt.Sprintf("---------------------\nCome up with a jogging routine	(Cancelled)\nID: 100\nCreated at: %v\nUpdated at: %v\n---------------------\n", tm, tm)
		if task.Format() != want {
			t.Errorf("The task formatting is not in the correct way. Got:%s, want:%s", task.Format(), want)
		}
	})

}

func TestTskStatusString(t *testing.T) {
	tests := []struct {
		status TskStatus
		want   string
	}{
		{Cancelled, "Cancelled"},
		{ToDo, "To-Do"},
		{InProgress, "In progress"},
		{Done, "Done"},
		{TskStatus(99), "None of them"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("TskStatus(%d).String() = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestStrToTskStatus(t *testing.T) {
	tests := []struct {
		input string
		want  TskStatus
	}{
		{"mark-in-progress", InProgress},
		{"in-progress", InProgress},
		{"mark-todo", ToDo},
		{"todo", ToDo},
		{"mark-done", Done},
		{"done", Done},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := StrToTskStatus(tt.input)
			if err != nil {
				t.Fatalf("StrToTskStatus(%q) returned unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("StrToTskStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStrToTskStatusRejectsUnknownValue(t *testing.T) {
	if _, err := StrToTskStatus("finished"); err == nil {
		t.Fatal("StrToTskStatus(\"finished\") should return an error")
	}
}
