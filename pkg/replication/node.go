package replication

// ID is an identifier of a replica node.
type ID string

// Node is a replica node.
type Node struct {
	ID   ID
	Addr Address
	Port int
}
