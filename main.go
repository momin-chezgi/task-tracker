package main

import (
	"errors"
	"fmt"
	"os"
)

type Query struct {
	Action Action
	Arg1   string
	Arg2   string
	Output []string
}

type Action int

const (
	AddCommand Action = iota + 1
	UpdateCommand
	CancelCommand
	MarkCommand
	ListCommand
)

var (
	lessThanOneErr             = errors.New("Please give at least one parameter!")
	invalidOptionErr           = errors.New("Please enter a valid option!")
	incorrectNumberOfArguments = errors.New("Incorrect number of arguments!")
)

func main() {
	var query Query

	if len(os.Args) < 2 {

		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", lessThanOneErr)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add":
		query.Action = AddCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
		}
		query.Arg1 = os.Args[2]
	case "update":
		query.Action = UpdateCommand
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "%v, want:3, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
		}
		query.Arg1 = os.Args[2]
		query.Arg2 = os.Args[3]
	case "cancel":
		query.Action = CancelCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
		}
		query.Arg1 = os.Args[2]
	case "mark-in-progress":
		query.Action = MarkCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
		}
		query.Arg1 = os.Args[1]
		query.Arg2 = os.Args[2]
	case "mark-done":
		query.Action = MarkCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
		}
		query.Arg1 = os.Args[1]
		query.Arg2 = os.Args[2]
	case "mark-todo":
		query.Action = MarkCommand
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
		}
		query.Arg1 = os.Args[1]
		query.Arg1 = os.Args[2]
	case "list":
		query.Action = ListCommand
		if len(os.Args) > 3 {
			fmt.Fprintf(os.Stderr, "%v, want:1 or 2 , got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
		}
		if len(os.Args) == 3 {
			query.Arg1 = os.Args[2]
		}
	default:
		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", invalidOptionErr)
		os.Exit(0)
	}

	var tm taskManager

	err := tm.ProcessTheQueryAndGetTheResult(&query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
		os.Exit(0)
	}

	for _, line := range query.Output {
		fmt.Println(line)
	}
}
