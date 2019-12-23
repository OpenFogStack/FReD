package inthandler

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
)

// Handler is an interface that abstracts the methods of the handler that handles internal requests.
type Handler interface {
	HandleCreateKeygroup(k keygroup.Keygroup) error
	HandleDeleteKeygroup(k keygroup.Keygroup) error
	HandleRead(i data.Item) (data.Item, error)
	HandleUpdate(i data.Item) error
	HandleDelete(i data.Item) error
}