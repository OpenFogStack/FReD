package data

// Service provides methods to manipulate items in the key-value store.
type Service interface {
	Read(i Item) (Item, error)
	Update(i Item) error
	Delete(i Item) error
	CreateKeygroup(i Item) error
	DeleteKeygroup(i Item) error
}
