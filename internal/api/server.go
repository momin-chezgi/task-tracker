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
