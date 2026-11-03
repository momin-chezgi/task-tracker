// To build a release:
// go build -ldflags="-X main.version=1.1.0" -o ..\..\task-tracker
package main

import (
	"fmt"

	"github.com/momin-chezgi/task-tracker/internal/api"
	"github.com/momin-chezgi/task-tracker/internal/store"
)

// var (
// 	version     = "dev"
// 	helpMessage = "\n\nTask Tracker (CLI)\n This is an application to create tasks, track them, change their status or even cancel them.\n\n commands are:\n$ add \"[Title of the task]\"\n$ update [id] \"[The new Title]\"\n$ cancel [id]\n$ mark-in-progress [id]\n$ mark-done [id]\n$ mark-todo [id]\n$ list\n$ list in-progress\n$ list done\n$ list todo\n$ show [id]\n$ path\n"
// )

// var (
// 	invalidOptionErr           = errors.New("Please enter a valid option!")
// 	incorrectNumberOfArguments = errors.New("Incorrect number of arguments!")
// )

func main() {
	tskStorage, err := store.Load()
	if err != nil {
		fmt.Printf("An error occurred in the tasks loading:%s\n", err.Error())
		return
	}

	defer func() {
		if err == nil {
			store.Store(tskStorage)
		}
	}()

	err = api.Serve(&tskStorage)
	if err != nil {
		fmt.Println("Server error:", err.Error())
		return
	}
}
