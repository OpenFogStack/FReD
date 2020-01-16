package main

import (
	"flag"
	"fmt"
	"os"

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
	nodeBhttpPort := flag.String("nodeBhttp", "9001", "port of nodeB (e.g. 9001)")
	nodeBzmqPort := flag.Int("nodeBzmqPort", 5555, "ZMQ Port of nodeB")
	nodeBzmqID := flag.String("nodeBzmqID", "nodeB", "ZMQ Id of nodeB")

	nodeBurl := fmt.Sprintf("http://%s:%s/%s/", *nodeBhost, *nodeBhttpPort, *apiVersion)

	nodeChost := flag.String("nodeChost", "localhost", "host of nodeC (e.g. localhost)")
	nodeChttpPort := flag.String("nodeChttp", "9001", "port of nodeC (e.g. 9001)")
	nodeCzmqPort := flag.Int("nodeCzmqPort", 5555, "ZMQ Port of nodeC")
	nodeCzmqID := flag.String("nodeCzmqID", "nodeC", "ZMQ Id of nodeC")

	nodeCurl := fmt.Sprintf("http://%s:%s/%s/", *nodeChost, *nodeChttpPort, *apiVersion)
	log.Debug().Msgf("Here are some variables: %s", nodeAzmqPort, nodeAzmqID, nodeBurl, nodeCzmqID, nodeCurl)
	flag.Parse()

	nodeA := pkg.NewNode(nodeAurl)

	var resp map[string]string
	// Test Keygroups
	logNodeAction(nodeA, "Creating keygroup testing")
	nodeA.CreateKeygroup("testing", 200, true)

	logNodeAction(nodeA, "Deleting keygroup testing")
	nodeA.DeleteKeygroup("testing", 200, true)

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
	nodeA.RegisterReplica(*nodeBzmqID, *nodeBhost, *nodeBzmqPort, 200, true)

	logNodeAction(nodeA, "Telling nodeA about nodeC")
	nodeA.RegisterReplica(*nodeBzmqID, *nodeChost, *nodeCzmqPort, 200, true)

	logNodeAction(nodeA, "Getting all Replicas that nodeA has")
	nodeA.GetAllReplica(200, true)
	// Test sending data between nodes
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
