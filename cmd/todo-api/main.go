// To build a release:
// go build -ldflags="-X main.version=1.1.0" -o ..\..\task-tracker
package main

import (
	"fmt"

	"github.com/momin-chezgi/task-tracker/internal/api"
	"github.com/momin-chezgi/task-tracker/internal/store"
)

var (
	version = "dev"
)

func main() {
	tskStorage, err := store.Load()
	if err != nil {
		fmt.Printf("An error occurred in the tasks loading:%s\n", err.Error())
		return
	}

	err = api.Serve(&tskStorage)
	if err != nil {
		fmt.Println("Server error:", err.Error())
		return
	}
}
