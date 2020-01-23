package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg"
)

var (
	// Wait for the user to press enter to continue
	waitUser = true
	reader = bufio.NewReader(os.Stdin)
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

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	nodeAhost := flag.String("nodeAhost", "172.26.0.10", "host of nodeA (e.g. localhost)")
=======
	nodeAhost := flag.String("nodeAhost", "localhost", "host of nodeA (e.g. localhost)")
>>>>>>> Add Mac Support, do not fatal on invalid json responses
	nodeAhttpPort := flag.String("nodeAhttp", "9001", "port of nodeA (e.g. 9001)")
	nodeAzmqhost := flag.String("nodeAzmqhost", "172.26.0.10", "host of nodeA (e.g. localhost) that can be reached by the other nodes")
	nodeAzmqPort := flag.Int("nodeAzmqPort", 5555, "ZMQ Port of nodeA")
	nodeAzmqID := flag.String("nodeAzmqID", "nodeA", "ZMQ Id of nodeA")

	nodeBhost := flag.String("nodeBhost", "localhost", "host of nodeB (e.g. localhost)")
	nodeBhttpPort := flag.String("nodeBhttp", "9002", "port of nodeB (e.g. 9001)")
	nodeBzmqhost := flag.String("nodeBzmqhost", "172.26.0.11", "host of nodeB (e.g. localhost) that can be reached by the other nodes")
	nodeBzmqPort := flag.Int("nodeBzmqPort", 5555, "ZMQ Port of nodeB")
	nodeBzmqID := flag.String("nodeBzmqID", "nodeB", "ZMQ Id of nodeB")

	nodeChost := flag.String("nodeChost", "localhost", "host of nodeC (e.g. localhost)")
	nodeChttpPort := flag.String("nodeChttp", "9003", "port of nodeC (e.g. 9001)")
	nodeCzmqhost := flag.String("nodeCzmqhost", "172.26.0.12", "host of nodeC (e.g. localhost) that can be reached by the other nodes")
	nodeCzmqPort := flag.Int("nodeCzmqPort", 5555, "ZMQ Port of nodeC")
	nodeCzmqID := flag.String("nodeCzmqID", "nodeC", "ZMQ Id of nodeC")

<<<<<<< HEAD

=======
	nodeCurl := fmt.Sprintf("http://%s:%s/%s/", *nodeChost, *nodeChttpPort, *apiVersion)
	log.Debug().Msg(string(*nodeAzmqPort) + "would be a unused var if not for this message")
>>>>>>> Count the number of errors, add run configuration
=======
	nodeAhost := flag.String("nodeAhost", "localhost", "host of nodeA (e.g. localhost)")
	nodeAhttpPort := flag.String("nodeAhttp", "80", "port of nodeA (e.g. 9001)")
=======
	nodeAhost := flag.String("nodeAhost", "172.26.0.10", "host of nodeA (e.g. localhost)")
	nodeAhttpPort := flag.String("nodeAhttp", "9001", "port of nodeA (e.g. 9001)")
>>>>>>> Finalize 3NodeTest, run tests in docker containers
	nodeAzmqPort := flag.Int("nodeAzmqPort", 5555, "ZMQ Port of nodeA")
	nodeAzmqID := flag.String("nodeAzmqID", "nodeA", "ZMQ Id of nodeA")

	nodeBhost := flag.String("nodeBhost", "172.26.0.11", "host of nodeB (e.g. localhost)")
	nodeBhttpPort := flag.String("nodeBhttp", "9001", "port of nodeB (e.g. 9001)")
	nodeBzmqPort := flag.Int("nodeBzmqPort", 5555, "ZMQ Port of nodeB")
	nodeBzmqID := flag.String("nodeBzmqID", "nodeB", "ZMQ Id of nodeB")

	nodeChost := flag.String("nodeChost", "172.26.0.12", "host of nodeC (e.g. localhost)")
	nodeChttpPort := flag.String("nodeChttp", "9001", "port of nodeC (e.g. 9001)")
	nodeCzmqPort := flag.Int("nodeCzmqPort", 5555, "ZMQ Port of nodeC")
	nodeCzmqID := flag.String("nodeCzmqID", "nodeC", "ZMQ Id of nodeC")


>>>>>>> change some bugs in trevers code
	flag.Parse()

	nodeAurl := fmt.Sprintf("http://%s:%s/%s/", *nodeAhost, *nodeAhttpPort, *apiVersion)
	log.Debug().Msgf("Node A: %s with ZMQ Port %d and ID %s", nodeAurl, *nodeAzmqPort, *nodeAzmqID)

	nodeBurl := fmt.Sprintf("http://%s:%s/%s/", *nodeBhost, *nodeBhttpPort, *apiVersion)
	log.Debug().Msgf("Node B: %s with ZMQ Port %d and ID %s", nodeBurl, *nodeBzmqPort, *nodeBzmqID)

	nodeCurl := fmt.Sprintf("http://%s:%s/%s/", *nodeChost, *nodeChttpPort, *apiVersion)
	log.Debug().Msgf("Node C: %s with ZMQ Port %d and ID %s", nodeCurl, *nodeCzmqPort, *nodeCzmqID)

	nodeA := pkg.NewNode(nodeAurl)
	nodeB := pkg.NewNode(nodeBurl)
	nodeC := pkg.NewNode(nodeCurl)

	var resp map[string]string

	// Seed NodeA
	logNodeAction(nodeA, "Seeding nodeA")
	nodeA.SeedNode(*nodeAzmqID, *nodeAzmqhost, 200, true)


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
	nodeA.PutItem("nonexistentkeygroup", "item","data", 409, false)

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
	logNodeAction(nodeA, "Telling nodeA about nodeB")
	nodeA.RegisterReplica(*nodeBzmqID, *nodeBzmqhost, *nodeBzmqPort, 200, true)

	logNodeAction(nodeA, "Telling nodeA about nodeC")
	nodeA.RegisterReplica(*nodeCzmqID, *nodeCzmqhost, *nodeCzmqPort, 200, true)

	logNodeAction(nodeA, "Getting all Replicas that nodeA has")
	parsed := nodeA.GetAllReplica(200, false)
	// Example Response: [{"Addr":{"Addr":"localhost","IsIP":false},"ID":"nodeB","Port":5556}]
	// Test for nodeA
	var nodeBnumber, nodeCnumber string
	if parsed.Path("0.ID").Data().(string) == *nodeBzmqID {
		nodeBnumber = "0"
		nodeCnumber = "1"
	} else {
		nodeBnumber = "1"
		nodeCnumber = "0"
	}
	if parsed.Path(nodeBnumber+".ID").Data().(string) != *nodeBzmqID {
		logNodeFaliure(nodeA, nodeBnumber+".ID == nodeB", parsed.Path("0.ID").String())
	}
	if int(parsed.Path(nodeBnumber+".Port").Data().(float64)) != *nodeBzmqPort {
		logNodeFaliure(nodeA, nodeBnumber+".Port == nodeBZmqPort", parsed.Path("0.Port").String())
	}
	if parsed.Path(nodeBnumber+".Addr.Addr").Data().(string) != *nodeBhost {
		logNodeFaliure(nodeA, nodeBnumber+".Addr.Addr == nodeBhost", parsed.Path("0.Addr.Addr").String())
	}
	// Test for nodeC
	if parsed.Path(nodeCnumber+".ID").Data().(string) != *nodeCzmqID {
		logNodeFaliure(nodeA, nodeCnumber+".ID == nodeC", parsed.Path("1.ID").String())
	}
	if int(parsed.Path(nodeCnumber+".Port").Data().(float64)) != *nodeCzmqPort {
		logNodeFaliure(nodeA, nodeCnumber+".Port == nodeCZmqPort", parsed.Path("1.Port").String())
	}
	if parsed.Path(nodeCnumber+".Addr.Addr").Data().(string) != *nodeChost {
		logNodeFaliure(nodeA, nodeCnumber+".Addr.Addr == nodeChost", parsed.Path("1.Addr.Addr").String())
	}

	// Fun with replicas
	logNodeAction(nodeA, "Adding nodeB as Replica node for KG1")
	nodeA.AddReplicaNode("KG1", *nodeBzmqID, 200, true)

	logNodeAction(nodeA, "Adding a replica for a nonexisting Keygroup")
	nodeA.AddReplicaNode("trololololo", *nodeBzmqID, 409, false)

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
	nodeA.AddReplicaNode("KGB", *nodeCzmqID, 409, false)
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

	if nodeA.Errors != 0 || nodeB.Errors != 0 || nodeC.Errors != 0 {
		log.Error().Msgf("Total Errors: %d", nodeA.Errors + nodeB.Errors + nodeC.Errors)
		os.Exit(1)
	}
}

func logNodeAction(node *pkg.Node, action string) {
	log.Info().Str("node", node.URL).Msg(action)
	if waitUser {
		log.Info().Msg("Please press enter to execute this step:")
		_,_,_ = reader.ReadLine()
	}
}

func logNodeFaliure(node *pkg.Node, expected, result string) {
	log.Warn().Str("node", node.URL).Msgf("expected: %s, but got: %#v", expected, result)
	if waitUser {
		log.Info().Msg("Please press enter to continue:")
		_,_,_ = reader.ReadLine()
	}
}

func logNodeFaliureMap(node *pkg.Node, expected string, result map[string]string) {
	log.Warn().Str("node", node.URL).Msgf("expected: %s, but got: %#v", expected, result)
	if waitUser {
		log.Info().Msg("Please press enter to continue:")
		_,_,_ = reader.ReadLine()
	}
}

func logDebugInfo(node *pkg.Node, info string) {
	log.Debug().Str("node", node.URL).Msg(info)
}
