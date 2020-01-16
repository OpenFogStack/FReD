package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg"
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
	apiVersion := flag.String("apiVersion", "v0", "API Version (e.g. v0)")

	nodeAhost := flag.String("nodeAhost", "localhost", "host of nodeA (e.g. localhost)")
	nodeAhttpPort := flag.String("nodeAhttp", "9001", "port of nodeA (e.g. 9001)")
	nodeAzmqPort := flag.Int("nodeAzmqPort", 5555, "ZMQ Port of nodeA")
	nodeAzmqID := flag.String("nodeAzmqID", "nodeA", "ZMQ Id of nodeA")

	nodeAurl := fmt.Sprintf("http://%s:%s/%s/", *nodeAhost, *nodeAhttpPort, *apiVersion)

	nodeBhost := flag.String("nodeBhost", "localhost", "host of nodeB (e.g. localhost)")
	nodeBhttpPort := flag.String("nodeBhttp", "9002", "port of nodeB (e.g. 9001)")
	nodeBzmqPort := flag.Int("nodeBzmqPort", 5556, "ZMQ Port of nodeB")
	nodeBzmqID := flag.String("nodeBzmqID", "nodeB", "ZMQ Id of nodeB")

	nodeBurl := fmt.Sprintf("http://%s:%s/%s/", *nodeBhost, *nodeBhttpPort, *apiVersion)

	nodeChost := flag.String("nodeChost", "localhost", "host of nodeC (e.g. localhost)")
	nodeChttpPort := flag.String("nodeChttp", "9003", "port of nodeC (e.g. 9001)")
	nodeCzmqPort := flag.Int("nodeCzmqPort", 5557, "ZMQ Port of nodeC")
	nodeCzmqID := flag.String("nodeCzmqID", "nodeC", "ZMQ Id of nodeC")

	nodeCurl := fmt.Sprintf("http://%s:%s/%s/", *nodeChost, *nodeChttpPort, *apiVersion)
	log.Debug().Msgf("Here are some variables: %s", nodeAzmqPort, nodeAzmqID, nodeBurl, nodeCzmqID, nodeCurl)
	flag.Parse()

	nodeA := pkg.NewNode(nodeAurl)
	nodeB := pkg.NewNode(nodeBurl)
	nodeC := pkg.NewNode(nodeCurl)

	var resp map[string]string
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
	resp = nodeA.GetItem("KG1", "KG1-Item", 200, false)
	if resp["Data"] != "KG1-Value" {
		logNodeFaliureMap(nodeA, "resp[\"Data\"] is \"KG1-Value\"", resp)
	}

	logNodeAction(nodeA,"Getting a Value from a nonexistent keygroup")
	nodeA.GetItem("trololoool", "wut", 404, false)

	logNodeAction(nodeA,"Putting a Value into a nonexistent keygroup")
	nodeA.PutItem("nonexistentkeygroup", "item","data", 400, false)

	logNodeAction(nodeA, "Putting new value KG1-Item/KG1-Value2 into KG1")
	nodeA.PutItem("KG1", "KG1-Item", "KG1-Value2", 200, true)

	logNodeAction(nodeA, "Getting the value in KG1")
	resp = nodeA.GetItem("KG1", "KG1-Item", 200, false)
	if resp["Data"] != "KG1-Value2" {
		logNodeFaliureMap(nodeA, "resp[\"Data\"] is \"KG1-Value2\"", resp)
	} else {
		logDebugInfo(nodeA, "Got "+resp["Data"])
	}

	logNodeAction(nodeA, "Deleting the value from KG1")
	nodeA.DeleteItem("KG1", "KG1-Item", 200, true)

	logNodeAction(nodeA, "Getting the deleted value in KG1")
	resp = nodeA.GetItem("KG1", "KG1-Item", 404, true)

	// Connect the nodes
	logNodeAction(nodeA, "Seeding nodeA")
	nodeA.SeedNode(*nodeAzmqID, *nodeAhost, 200, true)

	logNodeAction(nodeA, "Telling nodeA about nodeB")
	nodeA.RegisterReplica(*nodeBzmqID, *nodeBhost, *nodeBzmqPort, 200, true)

	logNodeAction(nodeA, "Telling nodeA about nodeC")
	nodeA.RegisterReplica(*nodeCzmqID, *nodeChost, *nodeCzmqPort, 200, true)

	logNodeAction(nodeA, "Getting all Replicas that nodeA has")
	parsed := nodeA.GetAllReplica(200, false)
	// Example Response: [{"Addr":{"Addr":"localhost","IsIP":false},"ID":"nodeB","Port":5556}]
	// Test for nodeA
	if parsed.Path("0.ID").Data().(string) != *nodeBzmqID {
		logNodeFaliure(nodeA, "0.ID == nodeB", parsed.Path("0.ID").String())
	}
	if int(parsed.Path("0.Port").Data().(float64)) != *nodeBzmqPort {
		logNodeFaliure(nodeA, "0.Port == nodeBZmqPort", parsed.Path("0.Port").String())
	}
	if parsed.Path("0.Addr.Addr").Data().(string) != *nodeBhost {
		logNodeFaliure(nodeA, "0.Addr.Addr == nodeBhost", parsed.Path("0.Addr.Addr").String())
	}
	// Test for nodeC
	if parsed.Path("1.ID").Data().(string) != *nodeCzmqID {
		logNodeFaliure(nodeA, "1.ID == nodeC", parsed.Path("1.ID").String())
	}
	if int(parsed.Path("1.Port").Data().(float64)) != *nodeCzmqPort {
		logNodeFaliure(nodeA, "1.Port == nodeCZmqPort", parsed.Path("1.Port").String())
	}
	if parsed.Path("1.Addr.Addr").Data().(string) != *nodeChost {
		logNodeFaliure(nodeA, "1.Addr.Addr == nodeChost", parsed.Path("1.Addr.Addr").String())
	}

	// Fun with replicas
	logNodeAction(nodeA, "Adding nodeB as Replica node for KG1")
	nodeA.AddReplicaNode("KG1", *nodeBzmqID, 200, true)

	logNodeAction(nodeA, "Adding a replica for a nonexisting Keygroup")
	nodeA.AddReplicaNode("trololololo", *nodeBzmqID, 400, false)

	// Test sending data between nodes
	logNodeAction(nodeB, "Creating a new Keygroup (KGN), setting nodeA as Replica node")
	nodeB.CreateKeygroup("KGN",200,true)
	nodeB.AddReplicaNode("KGN",*nodeAzmqID,200,true)

	logNodeAction(nodeB, "Putting something in KGN, testing whether nodeA gets Replica")
	nodeB.PutItem("KGN","Item","Value",200,true)
	time.Sleep(1000*time.Millisecond)
	resp = nodeA.GetItem("KGN", "Item", 200, false)
	if resp["Data"] != "Value" {
		logNodeFaliureMap(nodeA, "resp[\"Data\"] is \"Value\"", resp)
	}

	logNodeAction(nodeA, "Putting something in KGN, testing whether nodeB gets Replica")
	nodeB.PutItem("KGN","Item2","Value2",200,true)
	time.Sleep(1000*time.Millisecond)
	resp = nodeA.GetItem("KGN", "Item2", 200, false)
	if resp["Data"] != "Value2" {
		logNodeFaliureMap(nodeA, "resp[\"Data\"] is \"Value2\"", resp)
	}

	logNodeAction(nodeA, "Creating a new Keygroup (KGB) in nodeC and nodeA and then telling nodeA that nodeC is a Replica node")
	nodeA.CreateKeygroup("KGB",200,true)
	nodeC.CreateKeygroup("KGB",200,true)
	nodeC.CreateKeygroup("KGB",200,true)
	logNodeAction(nodeA,"...Adding the replica should throw an error since the keygroup already exists in the replica and is not already a replica")
	nodeA.AddReplicaNode("KGB", *nodeCzmqID, 400, false)
	nodeA.DeleteKeygroup("KGB", 200, true)
	nodeC.DeleteKeygroup("KGB",200,true)

	logNodeAction(nodeA,"Adding stuff to existing KG KG1 and then adding nodeB as replica")
	nodeA.PutItem("KG1", "KG1-Item1", "KG1-Value1", 200, true)
	nodeA.PutItem("KG1", "KG1-Item2", "KG1-Value2", 200, true)
	nodeA.PutItem("KG1", "KG1-Item3", "KG1-Value3", 200, true)
	nodeA.AddReplicaNode("KG1", *nodeBzmqID, 200, true)
	time.Sleep(1000 * time.Millisecond)
	logNodeAction(nodeA,"...Getting values from nodeB, they should have propagated")
	resp = nodeB.GetItem("KG1", "KG-Item1", 200, false)
	if resp["Data"] != "KG1-Value1" {
		logNodeFaliureMap(nodeA, "resp[\"Data\"] is \"KG1-Value1\"", resp)
	} else {
		logDebugInfo(nodeA, "Got "+resp["Data"])
	}
}

func logNodeAction(node *pkg.Node, action string) {
	log.Info().Str("node", node.URL).Msg(action)
}

func logNodeFaliure(node *pkg.Node, expected, result string) {
	log.Warn().Str("node", node.URL).Msgf("expected: %s, but got: %#v", expected, result)
}

func logNodeFaliureMap(node *pkg.Node, expected string, result map[string]string) {
	log.Warn().Str("node", node.URL).Msgf("expected: %s, but got: %#v", expected, result)
}

func logDebugInfo(node *pkg.Node, info string) {
	log.Debug().Str("node", node.URL).Msg(info)
}

func expectEmptyResponse(resp map[string]string) {
	if resp != nil {
		log.Warn().Msgf("Expected empty response, but got: %#v", resp)
	}
}
