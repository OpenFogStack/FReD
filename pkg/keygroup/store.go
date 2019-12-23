package keygroup

// Store is an interface that abstracts the component that stores Keygroups metadata.
type Store interface {
	Create(k Keygroup) error
	Delete(k Keygroup) error
	Exists(k Keygroup) bool
}
