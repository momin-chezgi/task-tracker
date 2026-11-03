package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/momin-chezgi/task-tracker/internal/store"
	"github.com/momin-chezgi/task-tracker/internal/task"
)

var idCounter = 0
var tskStorage *store.Storage

func Serve(tskPtr *store.Storage) error {
	tskStorage = tskPtr
	idCounter = len(tskStorage.Tasks)

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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	srv := &http.Server{
		Addr:              "localhost:8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		fmt.Println("Server listening on http://localhost:8080")
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal recieved, starting graceful shutdown")
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}

	shutdownctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := srv.Shutdown(shutdownctx)
	if err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("server stopped gracefully")

	return nil
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

	activetskStorage := make([]task.Task, 0)
	for _, rt := range tskStorage.Tasks {
		if rt.TskStatus != task.Cancelled {
			activetskStorage = append(activetskStorage, rt)
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

	marshaled, err := json.Marshal(reqTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(marshaled)

	store.Store(*tskStorage)
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
	w.WriteHeader(http.StatusOK)
	store.Store(*tskStorage)
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	store.Store(*tskStorage)

	if err := json.NewEncoder(w).Encode(tskStorage.Tasks[id]); err != nil {
		fmt.Println("encoding updated task:", err)
	}

}
