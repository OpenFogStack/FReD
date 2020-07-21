package data

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
)

// Item is an item in the key-value store.
type Item struct {
	Keygroup commons.KeygroupName
	ID       string
	Data     string
}
