package keygroup

// Service is an interface that abstracts the methods to manipulate keygroup metadata.
type Service interface {
	Create(k Keygroup) error
	Delete(k Keygroup) error
	Exists(k Keygroup) bool
}
