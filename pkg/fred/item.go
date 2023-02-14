package fred

import (
	"git.tu-berlin.de/mcc-fred/vclock"
)

// Item is an item in the key-value store.
type Item struct {
	Keygroup   KeygroupName
	ID         string
	Val        string
	Version    vclock.VClock
	Tombstoned bool
}
