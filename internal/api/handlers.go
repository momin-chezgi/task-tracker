package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/momin-chezgi/task-tracker/internal/store"
	"github.com/momin-chezgi/task-tracker/internal/task"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from your Go server!")
	name := r.URL.Query().Get("name")
	if strings.TrimSpace(name) == "" {
		name = "Guest"
	}
	fmt.Fprintln(w, "name: "+name)
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	activetskStorage := make([]task.Task, 0)

	var filterStt task.TskStatus

	filterStr := r.URL.Query().Get("marked_as")

	if filterStr != "" {
		var err error
		if filterStt, err = task.StrToTskStatus(filterStr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		for _, rt := range tskStorage.Tasks {
			if rt.TskStatus == filterStt {
				activetskStorage = append(activetskStorage, rt)
			}
		}
	} else {
		for _, rt := range tskStorage.Tasks {
			if rt.TskStatus != task.Cancelled {
				activetskStorage = append(activetskStorage, rt)
			}
		}
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(activetskStorage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func addNewTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var reqTask task.Task

	if err := json.NewDecoder(r.Body).Decode(&reqTask); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(reqTask.Title) == "" {
		http.Error(w, "The title can't be empty!", http.StatusBadRequest)
		return
	}

	reqTask.ID = idCounter
	idCounter++

	reqTask.TskStatus = task.ToDo

	reqTask.CreatedAt = time.Now()
	reqTask.UpdatedAt = time.Now()

	tskStorage.Tasks = append(tskStorage.Tasks, reqTask)

	if err := store.Store(*tskStorage); err != nil {
		log.Printf("storage error: %v\n", err)
	}

	marshaled, err := json.Marshal(reqTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(marshaled)
}

func showTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id < 0 || id >= idCounter {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if tskStorage.Tasks[id].TskStatus == task.Cancelled {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	marshaled, err := json.Marshal(tskStorage.Tasks[id])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(marshaled)
}

func cancelTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id < 0 || id >= idCounter {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if tskStorage.Tasks[id].TskStatus == task.Cancelled {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if tskStorage.Tasks[id].TskStatus == task.Done {
		http.Error(w, "this task has been done, you cannot mark a done task to cancelled!", http.StatusBadRequest)
		return
	}

	tskStorage.Tasks[id].TskStatus = task.Cancelled
	if err := store.Store(*tskStorage); err != nil {
		log.Printf("storage error: %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid task ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	if id < 0 || id >= idCounter {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	if tskStorage.Tasks[id].TskStatus == task.Cancelled {
		http.Error(w, "task has been deleted", http.StatusNotFound)
		return
	}

	var patch struct {
		Title     *string         `json:"title"`
		TskStatus *task.TskStatus `json:"task_status"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&patch); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if patch.Title == nil && patch.TskStatus == nil {
		http.Error(w, "provide at least one field to update", http.StatusBadRequest)
		return
	}

	if patch.TskStatus != nil {
		switch *patch.TskStatus {
		case task.Cancelled, task.ToDo, task.InProgress, task.Done:
		default:
			http.Error(w, "the given task status is invalid", http.StatusBadRequest)
			return
		}
	}

	if patch.Title != nil {
		if strings.TrimSpace(*patch.Title) == "" {
			http.Error(w, "the title can't be empty", http.StatusBadRequest)
			return
		}
		tskStorage.Tasks[id].Title = *patch.Title
		tskStorage.Tasks[id].UpdatedAt = time.Now()
	}

	if patch.TskStatus != nil {
		tskStorage.Tasks[id].TskStatus = task.TskStatus(*patch.TskStatus)
		tskStorage.Tasks[id].UpdatedAt = time.Now()
	}

	if err := store.Store(*tskStorage); err != nil {
		log.Printf("storage error: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tskStorage.Tasks[id]); err != nil {
		fmt.Println("encoding updated task:", err)
	}

}
