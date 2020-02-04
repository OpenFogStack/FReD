package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg/swag-node"
)

var (
	// Wait for the user to press enter to continue
	waitUser = false
	reader   = bufio.NewReader(os.Stdin)
)

func main() {
	// Logging Setup
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	// Parse Flags
	apiVersion := *flag.String("apiVersion", "v0", "API Version (e.g. v0)")

	nodeAhost := *flag.String("nodeAhost", "localhost", "host of nodeA (e.g. localhost)")                                                // Docker: localhost
	nodeAhttpPort := *flag.String("nodeAhttp", "9002", "port of nodeA (e.g. 9001)")                                                      // Docker: 9002
	nodeAzmqhost := *flag.String("nodeAzmqhost", "172.26.0.10", "host of nodeA (e.g. localhost) that can be reached by the other nodes") // Docker: 172.26.0.10
	nodeAzmqPort := *flag.Int("nodeAzmqPort", 5555, "ZMQ Port of nodeA")
	nodeAzmqID := *flag.String("nodeAzmqID", "nodeA", "ZMQ Id of nodeA")

	nodeBhost := *flag.String("nodeBhost", "localhost", "host of nodeB (e.g. localhost)")
	nodeBhttpPort := *flag.String("nodeBhttp", "9003", "port of nodeB (e.g. 9001)")
	nodeBzmqhost := *flag.String("nodeBzmqhost", "172.26.0.11", "host of nodeB (e.g. localhost) that can be reached by the other nodes")
	nodeBzmqPort := *flag.Int("nodeBzmqPort", 5555, "ZMQ Port of nodeB")
	nodeBzmqID := *flag.String("nodeBzmqID", "nodeB", "ZMQ Id of nodeB")

	nodeChost := *flag.String("nodeChost", "localhost", "host of nodeC (e.g. localhost)")
	nodeChttpPort := *flag.String("nodeChttp", "9004", "port of nodeC (e.g. 9001)")
	nodeCzmqhost := *flag.String("nodeCzmqhost", "172.26.0.12", "host of nodeC (e.g. localhost) that can be reached by the other nodes")
	nodeCzmqPort := *flag.Int("nodeCzmqPort", 5555, "ZMQ Port of nodeC")
	nodeCzmqID := *flag.String("nodeCzmqID", "nodeC", "ZMQ Id of nodeC")

	flag.Parse()

	nodeAurl := fmt.Sprintf("http://%s:%s/%s", nodeAhost, nodeAhttpPort, apiVersion)
	log.Debug().Msgf("Node A: %s with ZMQ Port %d and ID %s", nodeAurl, nodeAzmqPort, nodeAzmqID)

	nodeBurl := fmt.Sprintf("http://%s:%s/%s", nodeBhost, nodeBhttpPort, apiVersion)
	log.Debug().Msgf("Node B: %s with ZMQ Port %d and ID %s", nodeBurl, nodeBzmqPort, nodeBzmqID)

	nodeCurl := fmt.Sprintf("http://%s:%s/%s", nodeChost, nodeChttpPort, apiVersion)
	log.Debug().Msgf("Node C: %s with ZMQ Port %d and ID %s", nodeCurl, nodeCzmqPort, nodeCzmqID)

	nodeA := node.NewNode(nodeAurl)
	nodeB := node.NewNode(nodeBurl)
	nodeC := node.NewNode(nodeCurl)

	// Seed NodeA
	logNodeAction(nodeA, "Seeding nodeA")
	nodeA.SeedNode(nodeAzmqID, nodeAzmqhost, 200, true)

	// Test Keygroups
	logNodeAction(nodeA, "Creating keygroup testing")
	nodeA.CreateKeygroup("testing", 200, true)

	logNodeAction(nodeA, "Deleting keygroup testing")
	nodeA.DeleteKeygroup("testing", 200, true)

	logNodeAction(nodeA, "Deleting nonexistent keygroup")
	nodeA.DeleteKeygroup("trololololo", 404, false)

	logNodeAction(nodeA, "Creating Keygroup KG1")
	nodeA.CreateKeygroup("KG1", 200, true)

	// Test Get/Put of a single node
	logNodeAction(nodeA, "Putting KG1-Item/KG1-Value into KG1")
	nodeA.PutItem("KG1", "KG1-Item", "KG1-Value", 200, true)

	logNodeAction(nodeA, "Getting the value in KG1")

	resp := nodeA.GetItem("KG1", "KG1-Item", 200, false)

	if resp != "KG1-Value" {
		logNodeFailure(nodeA, "resp is \"KG1-Value\"", resp)
	}

	logNodeAction(nodeA, "Getting a Value from a nonexistent keygroup")
	nodeA.GetItem("trololoool", "wut", 404, false)

	logNodeAction(nodeA, "Putting a Value into a nonexistent keygroup")
	nodeA.PutItem("nonexistentkeygroup", "item", "data", 404, false)

	logNodeAction(nodeA, "Putting new value KG1-Item/KG1-Value2 into KG1")
	nodeA.PutItem("KG1", "KG1-Item", "KG1-Value2", 200, true)

	logNodeAction(nodeA, "Getting the value in KG1")
	resp = nodeA.GetItem("KG1", "KG1-Item", 200, false)
	if resp != "KG1-Value2" {
		logNodeFailure(nodeA, "resp is \"KG1-Value2\"", resp)
	} else {
		logDebugInfo(nodeA, "Got "+resp)
	}

	// Connect the nodes
	logNodeAction(nodeA, "Telling nodeA about nodeB")
	nodeA.RegisterReplica(nodeBzmqID, nodeBzmqhost, nodeBzmqPort, 200, true)

	logNodeAction(nodeA, "Telling nodeA about nodeC")
	nodeA.RegisterReplica(nodeCzmqID, nodeCzmqhost, nodeCzmqPort, 200, true)

	logNodeAction(nodeA, "Getting all Replicas that nodeA has")
	parsed := nodeA.GetAllReplica(200, false)
	// Example Response: ["nodeB", "nodeC"]
	// Test for nodeA

	if len(parsed) != 2 {
		logNodeFailure(nodeA, "len(parsed) == 2", fmt.Sprintf("%d", len(parsed)))
	}

	// sorry but i still love go
	check := (len(parsed) != 2) && func() bool {
		for _, elem := range parsed {
			if elem == nodeBzmqID {
				return true
			}
		}

		return false
	}() && func() bool {
		for _, elem := range parsed {
			if elem == nodeCzmqID {
				return true
			}
		}

		return false
	}()

	if check {
		logNodeFailure(nodeA, "parsed == ["+nodeBzmqID+", "+nodeCzmqID+"]", fmt.Sprintf("%#v", parsed))
	}

	// Fun with replicas
	logNodeAction(nodeA, "Adding nodeB as Replica node for KG1")
	nodeA.AddKeygroupReplica("KG1", nodeBzmqID, 200, true)

	logNodeAction(nodeB, "Deleting the value from KG1")
	nodeB.DeleteItem("KG1", "KG1-Item", 200, true)

	logNodeAction(nodeB, "Getting the deleted value in KG1")
	resp = nodeB.GetItem("KG1", "KG1-Item", 404, false)

	// Test sending data between nodes
	logNodeAction(nodeB, "Creating a new Keygroup (KGN) in nodeB, setting nodeA as Replica node")
	nodeB.CreateKeygroup("KGN", 200, true)
	nodeB.AddKeygroupReplica("KGN", nodeAzmqID, 200, true)

	logNodeAction(nodeB, "Putting something in KGN on nodeB, testing whether nodeA gets Replica (sleep 1.5s in between)")
	nodeB.PutItem("KGN", "Item", "Value", 200, true)
	time.Sleep(1500 * time.Millisecond)
	resp = nodeA.GetItem("KGN", "Item", 200, false)
	if resp != "Value" {
		logNodeFailure(nodeA, "resp is \"Value\"", resp)
	}

	logNodeAction(nodeA, "Putting something in KGN on nodeA, testing whether nodeB gets Replica (sleep 1.5s in between)")
	nodeA.PutItem("KGN", "Item2", "Value2", 200, true)
	time.Sleep(1500 * time.Millisecond)
	resp = nodeB.GetItem("KGN", "Item2", 200, false)
	if resp != "Value2" {
		logNodeFailure(nodeA, "resp is \"Value2\"", resp)
	}

	logNodeAction(nodeA, "Adding a replica for a nonexisting Keygroup")
	nodeA.AddKeygroupReplica("trololololo", nodeBzmqID, 404, false)

	logNodeAction(nodeA, "Creating a new Keygroup (KGB) in nodeC and nodeA and then telling nodeA that nodeC is a Replica node")
	nodeA.CreateKeygroup("KGB", 200, true)
	nodeC.CreateKeygroup("KGB", 200, true)
	nodeC.CreateKeygroup("KGB", 200, true)
	logNodeAction(nodeA, "...Adding the replica should throw an error since the keygroup already exists in the replica and is not already a replica")
	nodeA.AddKeygroupReplica("KGB", nodeCzmqID, 409, false)
	nodeA.DeleteKeygroup("KGB", 200, true)
	nodeC.DeleteKeygroup("KGB", 200, true)

	logNodeAction(nodeA, "Adding stuff to existing KG KG1 and then adding nodeB as replica (sleep 1.5s in between)")
	nodeA.PutItem("KG1", "KG1-Item1", "KG1-Value1", 200, true)
	nodeA.PutItem("KG1", "KG1-Item2", "KG1-Value2", 200, true)
	nodeA.PutItem("KG1", "KG1-Item3", "KG1-Value3", 200, true)
	nodeA.AddKeygroupReplica("KG1", nodeBzmqID, 200, true)
	time.Sleep(1500 * time.Millisecond)
	logNodeAction(nodeA, "...Getting values from nodeB, they should have propagated")
	resp = nodeB.GetItem("KG1", "KG-Item1", 200, false)
	if resp != "KG1-Value1" {
		logNodeFailure(nodeA, "resp is \"KG1-Value1\"", resp)
	} else {
		logDebugInfo(nodeA, "Got "+resp)
	}

	if nodeA.Errors != 0 || nodeB.Errors != 0 || nodeC.Errors != 0 {
		log.Error().Msgf("Total Errors: %d", nodeA.Errors+nodeB.Errors+nodeC.Errors)
		os.Exit(1)
	}
}

func logNodeAction(node *node.Node, action string) {
	wait()
	log.Info().Str("node", node.URL).Msg(action)
}

func logNodeFailure(node *node.Node, expected, result string) {
	wait()
	log.Warn().Str("node", node.URL).Msgf("expected: %s, but got: %#v", expected, result)
}

func logDebugInfo(node *node.Node, info string) {
	wait()
	log.Debug().Str("node", node.URL).Msg(info)
}

func wait() {
	if waitUser {
		log.Info().Msg("Please press enter to continue:")
		_, _, _ = reader.ReadLine()
	}
}
