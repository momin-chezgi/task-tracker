// To build a release:
// go build -ldflags="-X main.version=1.1.0" -o ..\..\task-tracker
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/momin-chezgi/task-tracker/internal/app"
	"github.com/momin-chezgi/task-tracker/internal/query"
)

var (
	version     = "dev"
	helpMessage = "\n\nTask Tracker (CLI)\n This is an application to create tasks, track them, change their status or even cancel them.\n\n commands are:\n$ add \"[Description of the task]\"\n$ update [id] \"[The new description]\"\n$ cancel [id]\n$ mark-in-progress [id]\n$ mark-done [id]\n$ mark-todo [id]\n$ list\n$ list in-progress\n$ list done\n$ list todo\n$ show [id]\n$ path\n"
)

var (
	lessThanOneErr             = errors.New("Please give at least one parameter!")
	invalidOptionErr           = errors.New("Please enter a valid option!")
	incorrectNumberOfArguments = errors.New("Incorrect number of arguments!")
)

func main() {
	var q query.Query

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stdout, helpMessage)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "help":
		fmt.Fprintln(os.Stdout, helpMessage)
		os.Exit(0)
	case "version":
		fmt.Fprintf(os.Stdout, "Version: %v\n", version)
		os.Exit(0)
	case "path":
		filePath, err := store.DataFilePath()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Where the data saved: %v", filePath)
		os.Exit(0)
	case "add":
		q.Action = query.AddCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
	case "update":
		q.Action = query.UpdateCommand
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "%v, want:3, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
		q.Arg2 = os.Args[3]
	case "cancel":
		q.Action = query.CancelCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
	case "mark-in-progress":
		q.Action = query.MarkCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[1]
		q.Arg2 = os.Args[2]
	case "mark-done":
		q.Action = query.MarkCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[1]
		q.Arg2 = os.Args[2]
	case "mark-todo":
		q.Action = query.MarkCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[1]
		q.Arg2 = os.Args[2]
	case "list":
		q.Action = query.ListCommand
		if len(os.Args) > 3 {
			fmt.Fprintf(os.Stderr, "%v, want:1 or 2 , got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		if len(os.Args) == 3 {
			q.Arg1 = os.Args[2]
		}
	case "show":
		q.Action = query.ShowCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
	default:
		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", invalidOptionErr)
		os.Exit(1)
	}

	var tm app.TaskManager

	err := tm.ProcessTheQueryAndGetTheResult(&q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
		os.Exit(1)
	}

	for _, line := range q.Output {
		fmt.Fprintln(os.Stdout, line)
	}
}
