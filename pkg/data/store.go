package data

// Store is an interface for the storage medium that the key-value data items are persisted on.
type Store interface {
	Update(i Item) error
	Delete(i Item) error
	Read(i Item) (Item, error)
	Exists(i Item) bool
}
