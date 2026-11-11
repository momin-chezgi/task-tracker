package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/momin-chezgi/task-tracker/internal/store"
	"github.com/momin-chezgi/task-tracker/internal/task"
)

func resetTaskStore(t *testing.T) {
	t.Helper()

	strg = &store.Storage{Tasks: make([]task.Task, 0)}
	task.TskIdCounter = 0

	t.Cleanup(func() {
		strg.Tasks = make([]task.Task, 0)
		task.TskIdCounter = 0
	})
}

func addASampleTaskJustForTesting(t *testing.T) {
	addingTask := task.Task{
		Title:     "This is a constant task.",
		TskStatus: task.InProgress,
	}
	marshaled, err := json.Marshal(addingTask)
	if err != nil {
		t.Fatalf("(Preprocessing:Marshaling)")
	}
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(marshaled))
	rr := httptest.NewRecorder()

	addNewTask(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("(Preprocessing:Add) expected status: %d, got %d",
			http.StatusCreated, rr.Code)
	}
}

func TestHello(t *testing.T) {
	tcNames := []string{"Bob", "Joe", "guest", "حسن", "محمد"}
	for _, tcn := range tcNames {
		t.Run("Test Hello", func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/hello/?name=%v", tcn), nil)
			rr := httptest.NewRecorder()

			sayHello(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d",
					http.StatusOK, rr.Code)
			}

			expected := fmt.Sprintf("Hello from your Go server!\nname: %v\n", tcn)
			if rr.Body.String() != expected {
				t.Errorf("expected body %q, got %q",
					expected, rr.Body.String())
			}
		})
	}

	t.Run("Empty name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/hello/", nil)
		rr := httptest.NewRecorder()

		sayHello(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d",
				http.StatusOK, rr.Code)
		}

		expected := "Hello from your Go server!\nname: Guest\n"
		if rr.Body.String() != expected {
			t.Errorf("expected body %q, got %q",
				expected, rr.Body.String())
		}
	})
}

func TestAddNewTask(t *testing.T) {
	resetTaskStore(t)

	t.Run("Add a task with empty title", func(t *testing.T) {
		emptyTask := task.Task{
			Title: "",
		}
		marshaled, err := json.Marshal(emptyTask)
		if err != nil {
			t.Fatalf("An error occurred while marshaling JSON: %v", err.Error())
		}
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(marshaled))
		rr := httptest.NewRecorder()

		addNewTask(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d",
				http.StatusBadRequest, rr.Code)
		}

		expected := json.RawMessage(`{"error":"The title can't be empty!"}` + "\n")
		if !reflect.DeepEqual(rr.Body.Bytes(), []byte(expected)) {
			t.Errorf("expected body %q, got %q",
				expected, rr.Body.String())
		}
	})

	tcTasks := []task.Task{
		{
			Title:     "Get the book",
			TskStatus: task.InProgress,
		},
		{
			Title:     "Eat the apple",
			TskStatus: task.Done,
		},
		{
			Title:     "Implement the last phase",
			TskStatus: task.ToDo,
		},
		{
			Title:     "Repair the bicycle",
			TskStatus: task.InProgress,
		},
		{
			Title:     "بعد از ظهر نان سنگک بخر",
			TskStatus: task.Done,
		},
	}

	addedNum := 0
	for _, tct := range tcTasks {
		marshaled, err := json.Marshal(tct)
		if err != nil {
			t.Errorf("marshaling error")
		}
		t.Run("add a new task", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(marshaled))
			rr := httptest.NewRecorder()

			addNewTask(rr, req)

			if rr.Code != http.StatusCreated {
				t.Fatalf("expected request %d, got %d",
					http.StatusCreated, rr.Code)
			}

			var responseTask task.Task
			if err := json.Unmarshal(rr.Body.Bytes(), &responseTask); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if addedNum != responseTask.ID {
				t.Errorf("Expected id %d, got %d",
					addedNum, responseTask.ID)
			}

			if tct.Title != responseTask.Title {
				t.Errorf("Expected title %v, got %v",
					tct.Title, responseTask.Title)
			}

			if responseTask.TskStatus != task.ToDo {
				t.Errorf("Expected new task status %v, got %v",
					task.ToDo, responseTask.TskStatus)
			}

			if tct.TskStatus == task.Cancelled {
				t.Errorf("Expected get a non-cancelled task!")
			}

			addedNum++
		})
	}
}

func TestUpdateTask(t *testing.T) {
	t.Run("update when no tasks are there", func(t *testing.T) {
		resetTaskStore(t)

		req := httptest.NewRequest(http.MethodPatch, "/tasks", nil)
		req.SetPathValue("id", "0")
		rr := httptest.NewRecorder()

		updateTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d",
				http.StatusNotFound, rr.Code)
		}
	})

	t.Run("update out of range index", func(t *testing.T) {
		resetTaskStore(t)

		addASampleTaskJustForTesting(t)

		taskSnapshot := strg.Tasks[0]

		req := httptest.NewRequest(http.MethodPatch, "/tasks", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		updateTask(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d",
				http.StatusNotFound, rr.Code)
		}
		if taskSnapshot.TskStatus != strg.Tasks[0].TskStatus ||
			taskSnapshot.ID != strg.Tasks[0].ID ||
			taskSnapshot.Title != strg.Tasks[0].Title {
			t.Errorf("Task#0 shouldn't have changed.")
		}
	})

	t.Run("reject an unknown status without changing the task", func(t *testing.T) {
		for _, status := range []int{-2, 3, 99} {
			t.Run(fmt.Sprintf("status %d", status), func(t *testing.T) {
				resetTaskStore(t)
				addASampleTaskJustForTesting(t)
				taskSnapshot := strg.Tasks[0]

				body, err := json.Marshal(struct {
					Title     string `json:"title"`
					TskStatus int    `json:"task_status"`
				}{
					Title:     "changed title",
					TskStatus: status,
				})
				if err != nil {
					t.Fatalf("failed to marshal patch: %v", err)
				}
				req := httptest.NewRequest(http.MethodPatch, "/tasks", bytes.NewReader(body))
				req.SetPathValue("id", "0")
				rr := httptest.NewRecorder()

				updateTask(rr, req)

				if rr.Code != http.StatusBadRequest {
					t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
				}
				if strg.Tasks[0] != taskSnapshot {
					t.Errorf("task should not change when the status is invalid")
				}
			})
		}
	})

	scenarios := []struct {
		numberOfTasks int
		updatingIndx  string
	}{
		{numberOfTasks: 2, updatingIndx: "0"},
		{numberOfTasks: 2, updatingIndx: "1"},
		{numberOfTasks: 3, updatingIndx: "0"},
		{numberOfTasks: 3, updatingIndx: "1"},
		{numberOfTasks: 3, updatingIndx: "2"},
		{numberOfTasks: 4, updatingIndx: "2"},
		{numberOfTasks: 5, updatingIndx: "4"},
	}
	sampleTasks := []task.Task{
		{Title: "Update your softwares", TskStatus: task.InProgress},
		{Title: "Get advice on learning English", TskStatus: task.Done},
		{Title: "Ask the teacher for advice", TskStatus: task.ToDo},
		{Title: "Consult the professor about the new project", TskStatus: task.ToDo},
		{Title: "repair the bike", TskStatus: task.Done},
	}

	for _, scnr := range scenarios {
		t.Run(fmt.Sprintf("Updating index #%s in an array with the length of %d",
			scnr.updatingIndx, scnr.numberOfTasks), func(t *testing.T) {
			resetTaskStore(t)

			addingTasks := sampleTasks[:scnr.numberOfTasks]
			for _, tsk := range addingTasks {
				marshaled, err := json.Marshal(tsk)
				if err != nil {
					t.Fatalf("(Preprocessing:Marshaling)")
				}
				req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(marshaled))
				rr := httptest.NewRecorder()

				addNewTask(rr, req)

				if rr.Code != http.StatusCreated {
					t.Fatalf("(Preprocessing:Add) expected status %d, got %d",
						http.StatusCreated, rr.Code)
				}
			}

			idx, _ := strconv.Atoi(scnr.updatingIndx)

			patch := struct {
				Title     string         `json:"title"`
				TskStatus task.TskStatus `json:"task_status"`
			}{
				Title:     addingTasks[idx].Title,
				TskStatus: addingTasks[idx].TskStatus,
			}
			marshaled, err := json.Marshal(patch)
			if err != nil {
				t.Fatalf("(Preprocessing:Marshaling)")
			}
			req := httptest.NewRequest(http.MethodPatch, "/tasks", bytes.NewReader(marshaled))
			req.SetPathValue("id", scnr.updatingIndx)
			rr := httptest.NewRecorder()

			updateTask(rr, req)
			if rr.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d",
					http.StatusOK, rr.Code)
			}
			responseTask := &task.Task{}
			if err := json.Unmarshal(rr.Body.Bytes(), responseTask); err != nil {
				t.Errorf("The format of response body is not valid")
			} else {

				if responseTask.TskStatus != strg.Tasks[idx].TskStatus ||
					responseTask.Title != strg.Tasks[idx].Title {
					t.Errorf("The response task and updated task are not equal")
				}
				reqTask := addingTasks[idx]
				if reqTask.TskStatus != responseTask.TskStatus ||
					reqTask.Title != responseTask.Title {
					t.Errorf("The requested task and the responded task are not the same")
				}
			}

		})
	}
}

func TestCancelTask(t *testing.T) {

	t.Run("Try to cancel an out of the range task", func(t *testing.T) {
		resetTaskStore(t)

		addASampleTaskJustForTesting(t)
		taskSnapshot := strg.Tasks[0]

		req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		cancelTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d",
				http.StatusNotFound, rr.Code)
		}
		if strg.Tasks[0] != taskSnapshot {
			t.Errorf("Task#0 shouldn't have changed.")
		}
	})

	t.Run("Reject a malformed task ID", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)
		taskSnapshot := strg.Tasks[0]

		req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		req.SetPathValue("id", "invalid")
		rr := httptest.NewRecorder()

		cancelTask(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d",
				http.StatusBadRequest, rr.Code)
		}
		if strg.Tasks[0] != taskSnapshot {
			t.Errorf("Task#0 shouldn't have changed.")
		}
	})

	t.Run("Reject a negative task ID", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)
		taskSnapshot := strg.Tasks[0]

		req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		req.SetPathValue("id", "-1")
		rr := httptest.NewRecorder()

		cancelTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d",
				http.StatusNotFound, rr.Code)
		}
		if strg.Tasks[0] != taskSnapshot {
			t.Errorf("Task#0 shouldn't have changed.")
		}
	})

	t.Run("Try to cancel a task when no tasks exist", func(t *testing.T) {
		resetTaskStore(t)

		req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		req.SetPathValue("id", "0")
		rr := httptest.NewRecorder()

		cancelTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d",
				http.StatusNotFound, rr.Code)
		}
		if len(strg.Tasks) != 0 {
			t.Errorf("expected no tasks, got %d", len(strg.Tasks))
		}
	})

	t.Run("Cancel an existing task", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)

		req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		req.SetPathValue("id", "0")
		rr := httptest.NewRecorder()

		cancelTask(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d",
				http.StatusOK, rr.Code)
		}
		if strg.Tasks[0].TskStatus != task.Cancelled {
			t.Errorf("expected task to be marked deleted")
		}
	})

	t.Run("Cancel only the selected task", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)
		addASampleTaskJustForTesting(t)
		firstTaskSnapshot := strg.Tasks[0]

		req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		cancelTask(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d",
				http.StatusOK, rr.Code)
		}
		if strg.Tasks[0] != firstTaskSnapshot {
			t.Errorf("Task#0 shouldn't have changed.")
		}
		if strg.Tasks[1].TskStatus != task.Cancelled {
			t.Errorf("expected Task#1 to be marked deleted")
		}
	})

	t.Run("Reject repeated cancellation", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)

		firstReq := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		firstReq.SetPathValue("id", "0")
		firstRecorder := httptest.NewRecorder()
		cancelTask(firstRecorder, firstReq)
		if firstRecorder.Code != http.StatusOK {
			t.Fatalf("expected first status %d, got %d",
				http.StatusOK, firstRecorder.Code)
		}
		taskSnapshot := strg.Tasks[0]

		secondReq := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		secondReq.SetPathValue("id", "0")
		secondRecorder := httptest.NewRecorder()
		cancelTask(secondRecorder, secondReq)

		if secondRecorder.Code != http.StatusNotFound {
			t.Errorf("expected second status %d, got %d",
				http.StatusNotFound, secondRecorder.Code)
		}
		if strg.Tasks[0] != taskSnapshot {
			t.Errorf("Task#0 shouldn't have changed after repeated cancellation.")
		}
	})

	t.Run("Don't cancel a DONE task", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)

		req := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
		req.SetPathValue("id", "0")
		rr := httptest.NewRecorder()

		strg.Tasks[0].TskStatus = task.Done

		cancelTask(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Done tasks can't be deleted, http status got:%v, want:%v",
				rr.Code, http.StatusBadRequest)
		}
	})
}

func TestListAllTasks(t *testing.T) {
	t.Run("returns an empty array when no tasks exist", func(t *testing.T) {
		resetTaskStore(t)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rr := httptest.NewRecorder()

		listTasks(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d",
				http.StatusOK, rr.Code)
		}
		if contentType := rr.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
			t.Errorf("expected content type %q, got %q",
				"application/json; charset=utf-8", contentType)
		}

		var responseTasks []task.Task
		if err := json.Unmarshal(rr.Body.Bytes(), &responseTasks); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if responseTasks == nil || len(responseTasks) != 0 {
			t.Errorf("expected an empty array, got %v", responseTasks)
		}
	})

	t.Run("returns active tasks in insertion order", func(t *testing.T) {
		resetTaskStore(t)
		strg.Tasks = []task.Task{
			{ID: 0, Title: "first active task", TskStatus: task.ToDo},
			{ID: 1, Title: "deleted task", TskStatus: task.Cancelled},
			{ID: 2, Title: "last active task", TskStatus: task.Done},
		}

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		rr := httptest.NewRecorder()

		listTasks(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d",
				http.StatusOK, rr.Code)
		}
		if contentType := rr.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
			t.Errorf("expected content type %q, got %q",
				"application/json; charset=utf-8", contentType)
		}

		var responseTasks []task.Task
		if err := json.Unmarshal(rr.Body.Bytes(), &responseTasks); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		expectedTasks := []task.Task{
			{ID: 0, Title: "first active task", TskStatus: task.ToDo},
			{ID: 2, Title: "last active task", TskStatus: task.Done},
		}
		if len(responseTasks) != len(expectedTasks) {
			t.Fatalf("expected %d tasks, got %d",
				len(expectedTasks), len(responseTasks))
		}
		for i, expectedTask := range expectedTasks {
			if responseTasks[i] != expectedTask {
				t.Errorf("task %d: expected %+v, got %+v",
					i, expectedTask, responseTasks[i])
			}
		}
	})
}

func TestFilterTasks(t *testing.T) {
	t.Run("returns an empty array when no tasks exist", func(t *testing.T) {
		resetTaskStore(t)

		req := httptest.NewRequest(http.MethodGet, "/tasks?marked_as=todo", nil)
		rr := httptest.NewRecorder()

		listTasks(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d",
				http.StatusOK, rr.Code)
		}
		if contentType := rr.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
			t.Errorf("expected content type %q, got %q",
				"application/json; charset=utf-8", contentType)
		}

		var responseTasks []task.Task
		if err := json.Unmarshal(rr.Body.Bytes(), &responseTasks); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if responseTasks == nil || len(responseTasks) != 0 {
			t.Errorf("expected an empty array, got %v", responseTasks)
		}
	})

	t.Run("returns only tasks matching the requested filter", func(t *testing.T) {
		resetTaskStore(t)
		strg.Tasks = []task.Task{
			{ID: 0, Title: "first todo", TskStatus: task.ToDo},
			{ID: 1, Title: "second done", TskStatus: task.Done},
			{ID: 2, Title: "third in progress", TskStatus: task.InProgress},
			{ID: 3, Title: "fourth todo", TskStatus: task.ToDo},
			{ID: 4, Title: "fifth cancelled", TskStatus: task.Cancelled},
		}

		req := httptest.NewRequest(http.MethodGet, "/tasks?marked_as=todo", nil)
		rr := httptest.NewRecorder()

		listTasks(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d",
				http.StatusOK, rr.Code)
		}

		var responseTasks []task.Task
		if err := json.Unmarshal(rr.Body.Bytes(), &responseTasks); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		expectedTasks := []task.Task{
			{ID: 0, Title: "first todo", TskStatus: task.ToDo},
			{ID: 3, Title: "fourth todo", TskStatus: task.ToDo},
		}
		if len(responseTasks) != len(expectedTasks) {
			t.Fatalf("expected %d tasks, got %d", len(expectedTasks), len(responseTasks))
		}
		for i, expectedTask := range expectedTasks {
			if responseTasks[i] != expectedTask {
				t.Errorf("task %d: expected %+v, got %+v", i, expectedTask, responseTasks[i])
			}
		}
	})

	t.Run("rejects an invalid filter value", func(t *testing.T) {
		resetTaskStore(t)
		strg.Tasks = []task.Task{{ID: 0, Title: "todo task", TskStatus: task.ToDo}}

		req := httptest.NewRequest(http.MethodGet, "/tasks?marked_as=not-a-status", nil)
		rr := httptest.NewRecorder()

		listTasks(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d",
				http.StatusBadRequest, rr.Code)
		}
	})
}

func TestShowTask(t *testing.T) {
	t.Run("rejects a malformed ID", func(t *testing.T) {
		resetTaskStore(t)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.SetPathValue("id", "invalid")
		rr := httptest.NewRecorder()

		showTask(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d",
				http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("returns not found for a negative ID", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.SetPathValue("id", "-1")
		rr := httptest.NewRecorder()

		showTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d",
				http.StatusNotFound, rr.Code)
		}
	})

	t.Run("returns not found for an out-of-range ID", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		showTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d",
				http.StatusNotFound, rr.Code)
		}
	})

	t.Run("returns not found when no tasks exist", func(t *testing.T) {
		resetTaskStore(t)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.SetPathValue("id", "0")
		rr := httptest.NewRecorder()

		showTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d",
				http.StatusNotFound, rr.Code)
		}
	})

	t.Run("returns not found for a deleted task", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)
		strg.Tasks[0].TskStatus = task.Cancelled

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.SetPathValue("id", "0")
		rr := httptest.NewRecorder()

		showTask(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d",
				http.StatusNotFound, rr.Code)
		}
	})

	t.Run("returns the requested task", func(t *testing.T) {
		resetTaskStore(t)
		addASampleTaskJustForTesting(t)

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		req.SetPathValue("id", "0")
		rr := httptest.NewRecorder()

		showTask(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d",
				http.StatusOK, rr.Code)
		}
		if contentType := rr.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
			t.Errorf("expected content type %q, got %q",
				"application/json; charset=utf-8", contentType)
		}

		var responseTask task.Task
		if err := json.Unmarshal(rr.Body.Bytes(), &responseTask); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		expectedTask := task.Task{
			ID:        0,
			Title:     "This is a constant task.",
			TskStatus: task.ToDo,
		}
		if responseTask != expectedTask {
			t.Errorf("expected task %+v, got %+v",
				expectedTask, responseTask)
		}
	})
}
