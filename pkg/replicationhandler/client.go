package replicationhandler

import (
	"net"
)

// Client is an interface to send replication messages across nodes.
type Client interface {
	SendCreateKeygroup(ip net.IP, port int, kgname string) error
	SendDeleteKeygroup(ip net.IP, port int, kgname string) error
	SendUpdate(ip net.IP, port int, kgname string, kgid string , value string) error
	SendDelete(ip net.IP, port int, kgname string, kgid string) error
}