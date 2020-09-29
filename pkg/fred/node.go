package fred

// Address is an IP address or a hostname of a FReD node.
type Address struct {
	Addr string
	IsIP bool
}

// NodeID is an identifier of a replica node.
type NodeID string

// Node is a replica node.
type Node struct {
	ID   NodeID
	Host string
}
