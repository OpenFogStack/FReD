package replication

// Store is an interface that encapsulates methods to store replication configuration.
type Store interface {
	CreateNode(n Node) error
	DeleteNode(n Node) error
	GetNode(n Node) (Node, error)
	ExistsNode(n Node) bool
	CreateKeygroup(k Keygroup) error
	DeleteKeygroup(k Keygroup) error
	GetKeygroup(k Keygroup) (Keygroup, error)
	ExistsKeygroup(k Keygroup) bool
	AddReplica(k Keygroup, n Node) error
	RemoveReplica(k Keygroup, n Node) error
}
