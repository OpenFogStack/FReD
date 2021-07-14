package main

import (
	"git.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg/grpcclient"
)

// TestSuite represents a group of tests that should be run together
type TestSuite interface {
	// RunTests runs all tests of this testsuite
	RunTests()
	Name() string
}

type Config struct {
	waitUser bool

	nodeAhost     string
	nodeAhttpPort string

	nodeBhost     string
	nodeBhttpPort string

	nodeChost     string
	nodeChttpPort string

	triggerNodeHost   string
	triggerNodeWSHost string
	triggerNodeID     string

	certFile string
	keyFile  string

	littleCertFile string
	littleKeyFile  string

	nodeA        *grpcclient.Node
	nodeB        *grpcclient.Node
	nodeC        *grpcclient.Node
	littleClient *grpcclient.Node
}
