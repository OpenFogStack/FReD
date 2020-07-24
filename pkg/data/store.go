package data

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
)

// Store is an interface for the storage medium that the key-value data items are persisted on.
type Store interface {
	// Needs: keygroup, id, data
	Update(i Item) error
	// Needs: keygroup, id
	Delete(kg commons.KeygroupName, id string) error
	// Needs: keygroup, id; Returns: data
	Read(kg commons.KeygroupName, id string) (string, error)
	// Needs: keygroup; Returns: keygroup, id, data
	ReadAll(kg commons.KeygroupName) ([]Item, error)
	// Needs: keygroup, Returns:[] keygroup, id
	IDs(kg commons.KeygroupName) ([]Item, error)
	// Needs: keygroup, id
	Exists(kg commons.KeygroupName, id string) bool
	// Needs: keygroup, Returns: err
	// Doesnt really need to store the KG itself, that is NaSe's job.
	// But might be useful for databases using it
	CreateKeygroup(kg commons.KeygroupName) error
	// Same as with CreateKeygroup
	DeleteKeygroup(kg commons.KeygroupName) error
}
