package replicationhandler

import (
	"net"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

// Client is an interface to send replication messages across nodes.
type Client interface {
	SendCreateKeygroup(ip net.IP, port int, kgname commons.KeygroupName) error
	SendDeleteKeygroup(ip net.IP, port int, kgname commons.KeygroupName) error
	SendUpdate(ip net.IP, port int, kgname commons.KeygroupName, id string, value string) error
	SendDelete(ip net.IP, port int, kgname commons.KeygroupName, id string) error
	SendAddNode(ip net.IP, port int, nodeID replication.ID, nodeIP net.IP, nodePort int) error
	SendRemoveNode(ip net.IP, port int, nodeID replication.ID) error
	SendAddReplica(ip net.IP, port int, kgname commons.KeygroupName, nodeID replication.ID) error
	SendRemoveReplica(ip net.IP, port int, kgname commons.KeygroupName, nodeID replication.ID) error
}
