package replicationhandler

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

// Client is an interface to send replication messages across nodes.
type Client interface {
	SendCreateKeygroup(addr replication.Address, port int, kgname commons.KeygroupName) error
	SendDeleteKeygroup(addr replication.Address, port int, kgname commons.KeygroupName) error
	SendUpdate(addr replication.Address, port int, kgname commons.KeygroupName, id string, value string) error
	SendDelete(addr replication.Address, port int, kgname commons.KeygroupName, id string) error
	SendAddNode(addr replication.Address, port int, node replication.Node) error
	SendRemoveNode(addr replication.Address, port int, node replication.Node) error
	SendAddReplica(addr replication.Address, port int, kgname commons.KeygroupName, node replication.Node) error
	SendRemoveReplica(addr replication.Address, port int, kgname commons.KeygroupName, node replication.Node) error
}
