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
	helpMessage = `\n\nTask Tracker (CLI)\n
	This is an application to create tasks, track them, change their status or even cancel them.\n
	+ Task-related commands are:\n
	$ add \"[Description]\"\n
	$ update [ID] \"[Description]\"\n
	$ cancel [ID]\n
	$ mark-in-progress [ID]\n
	$ mark-done [ID]\n
	$ mark-todo [ID]\n
	$ list\n
	$ list in-progress\n
	$ list done\n
	$ list todo\n
	$ show [ID]\n
	+ User-related commands are:\n
	$ user_add \"[Name]\" \"[Email]\" \"[Password]\"\n
	$ user_update_name \"[Name]\"\n
	$ user_update_email \"[Email]\"\n
	$ user_update_pass \"[Password]\"\n
	$ user_list \n
	$ user_show \"[ID]\"\n\n`
)

var (
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
	case "add":
		q.Action = query.AddTsk
		if len(os.Args) != 5 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
		q.Arg2 = os.Args[3]
		q.Arg3 = os.Args[4]
	case "update":
		q.Action = query.UpdateTsk
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "%v, want:3, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
		q.Arg2 = os.Args[3]
	case "cancel":
		q.Action = query.CancelTsk
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
	case "mark-in-progress":
		q.Action = query.MarkTsk
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[1]
		q.Arg2 = os.Args[2]
	case "mark-done":
		q.Action = query.MarkTsk
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[1]
		q.Arg2 = os.Args[2]
	case "mark-todo":
		q.Action = query.MarkTsk
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[1]
		q.Arg2 = os.Args[2]
	case "list":
		q.Action = query.ListTsk
		if len(os.Args) > 3 {
			fmt.Fprintf(os.Stderr, "%v, want:1 or 2 , got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		if len(os.Args) == 3 {
			q.Arg1 = os.Args[2]
		}
	case "show":
		q.Action = query.ShowTsk
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
	case "user_add":
		q.Action = query.AddUsr
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
		q.Arg2 = os.Args[3]
	case "user_update_name":
		q.Action = query.UpdateUsrName
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "%v, want:3, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
		q.Arg2 = os.Args[3]
	case "user_update_email":
		q.Action = query.UpdateUsrEmail
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "%v, want:3, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
		q.Arg2 = os.Args[3]
	case "user_update_pass":
		q.Action = query.UpdateUsrPass
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "%v, want:3, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
		q.Arg2 = os.Args[3]
	case "user_list":
		q.Action = query.ListUsr
		if len(os.Args) > 3 {
			fmt.Fprintf(os.Stderr, "%v, want:1 or 2 , got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
	case "user_show":
		q.Action = query.ShowTsk
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%v, want:2, got:%v\n", incorrectNumberOfArguments, len(os.Args)-1)
			os.Exit(1)
		}
		q.Arg1 = os.Args[2]
	default:
		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", invalidOptionErr)
		os.Exit(1)
	}

	err := app.ProcessTheQueryAndGetTheResult(&q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
		os.Exit(1)
	}

	for _, line := range q.Output {
		fmt.Fprintln(os.Stdout, line)
	}
}
