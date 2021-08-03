package fred

import (
	"github.com/DistributedClocks/GoVector/govec/vclock"
)

// Item is an item in the key-value store.
type Item struct {
	Keygroup   KeygroupName
	ID         string
	Val        string
	Version    vclock.VClock
	Tombstoned bool
}
