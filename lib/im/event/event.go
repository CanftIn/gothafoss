package event

import "github.com/gocraft/dbr/v2"

type Type int

const (
	None Type = iota
	Message
	Command
)

func (t Type) Int() int {
	return int(t)
}

type Status int

const (
	Wait Type = iota
	Success
	Fail
)

func (s Status) Int() int {
	return int(s)
}

type Data struct {
	Event string
	Type  Type
	Data  interface{}
}

type Event interface {
	Begin(data *Data, tx *dbr.Tx) (int64, error)
	Commit(eventId int64)
}
