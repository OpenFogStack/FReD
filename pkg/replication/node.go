package replication

import (
	"net"
)

// ID is an identifier of a replica node.
type ID string

// Node is a replica node.
type Node struct {
	ID   ID
	IP   net.IP
	Port int
}
