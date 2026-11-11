package query

type Query struct {
	Action Action
	Arg1   string
	Arg2   string
	Output []string
}

type Action int

const (
	AddTsk Action = iota + 1
	UpdateTsk
	CancelTsk
	MarkTsk
	ListTsk
	ShowTsk

	AddUsr
	UpdateUsrName
	UpdateUsrEmail
	UpdateUsrPass
	// DeleteUsr
	ListUsr
	ShowUsr
)
