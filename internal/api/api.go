package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/momin-chezgi/task-tracker/internal/task"
)

var tasks = make([]task.Task, 0)
var idCounter = 0

func Serve() {

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", sayHello)

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is my first API.")
	})

	mux.HandleFunc("/goodbye", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Goodbye! See you next time!")
	})

	mux.HandleFunc("POST /tasks", addNewTask)

	mux.HandleFunc("PATCH /tasks/{id}", updateTask)

	mux.HandleFunc("DELETE /tasks/{id}", cancelTask)

	mux.HandleFunc("GET /tasks", listAllTasks)

	mux.HandleFunc("GET /tasks/{id}", showTask)

	fmt.Println("Server listening on http://localhost:8080")

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		fmt.Println("Server error:", err.Error())
	}

}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from your Go server!")
	name := r.URL.Query().Get("name")
	if strings.TrimSpace(name) == "" {
		name = "Guest"
	}
	fmt.Fprintln(w, "name: "+name)
}

func listAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	activeTasks := make([]task.Task, 0)
	for _, rt := range tasks {
		if rt.TskStatus != task.Cancelled {
			activeTasks = append(activeTasks, rt)
		}
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(activeTasks); err != nil {
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

	tasks = append(tasks, reqTask)

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

	if tasks[id].TskStatus == task.Cancelled {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	marshaled, err := json.Marshal(tasks[id])
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

	if tasks[id].TskStatus == task.Cancelled {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if tasks[id].TskStatus == task.Done {
		http.Error(w, "this task has been done, you cannot mark a done task to cancelled!", http.StatusBadRequest)
		return
	}

	tasks[id].TskStatus = task.Cancelled
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
	if tasks[id].TskStatus == task.Cancelled {
		http.Error(w, "task has been deleted", http.StatusNotFound)
		return
	}

	var patch struct {
		Title     *string `json:"title"`
		TskStatus *int    `json:"completed"`
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
		switch task.TskStatus(*patch.TskStatus) {
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
		tasks[id].Title = *patch.Title
	}

	if patch.TskStatus != nil {
		tasks[id].TskStatus = task.TskStatus(*patch.TskStatus)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks[id]); err != nil {
		fmt.Println("encoding updated task:", err)
	}

}
