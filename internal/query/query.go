package query

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
	ShowCommand
)
