package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/momin-chezgi/task-tracker/internal/app"
	"github.com/momin-chezgi/task-tracker/internal/query"
)

var (
	lessThanOneErr             = errors.New("Please give at least one parameter!")
	invalidOptionErr           = errors.New("Please enter a valid option!")
	incorrectNumberOfArguments = errors.New("Incorrect number of arguments!")
)

func main() {
	var q query.Query

	if len(os.Args) < 2 {

		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", lessThanOneErr)
		os.Exit(1)
	}

	switch os.Args[1] {
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
		fmt.Println(line)
	}
}
