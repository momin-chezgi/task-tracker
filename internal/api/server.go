package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/momin-chezgi/task-tracker/internal/store"
)

var (
	helpMessage = `\n\nTask Tracker (API)\n
 This is an application to create tasks, track them, change their status or even cancel them.\n\n
 You can...\n
 + add a new task [POST /tasks]\n
 + update them [PATCH /tasks/{id}]\n
 + cancel them [DELETE /tasks/{id}]\n
 + list them [GET /tasks]\n
 + filter tasks by their status(todo, in-progress, done) [GET /tasks?marked_as=[The filtering status]]
 + show a specific task by its ID [GET /tasks/{id}]\n`
)

var (
	idCounter  = 0
	tskStorage *store.Storage
)

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
	mux.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, helpMessage)
	})
	mux.HandleFunc("POST /tasks", addNewTask)
	mux.HandleFunc("PATCH /tasks/{id}", updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", cancelTask)
	mux.HandleFunc("GET /tasks", listTasks)
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
		log.Println("shutdown signal received, starting graceful shutdown")
		if err := store.Store(*tskStorage); err != nil {
			log.Printf("storage error: %v\n", err)
		}
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("server stopped gracefully")

	return nil
}
