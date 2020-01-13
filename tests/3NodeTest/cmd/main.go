package main

import (
	"flag"
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
	nodeAurl := flag.String("nodeA", "http://localhost:9001/v0/", "ip:port/apiVersion/ of nodeA")
	nodeAzmqIp := flag.String("nodeAzmqIp", "localhost", "ZMQ IP of nodeA")
	nodeAzmqPort := flag.Int("nodeAzmqPort", 5555, "ZMQ Port of nodeA")
	nodeBurl := flag.String("nodeB", "http://localhost:9002/v0/", "ip:port/apiVersion/ of nodeB")
	nodeBzmqIp := flag.String("nodeBzmqIp", "localhost", "ZMQ IP of nodeB")
	nodeBzmqPort := flag.Int("nodeBzmqPort", 5556, "ZMQ Port of nodeB")
	nodeCurl := flag.String("nodeC", "http://localhost:9003/v0/", "ip:port/apiVersion/ of nodeC")
	nodeCzmqIp := flag.String("nodeCzmqIp", "localhost", "ZMQ IP of nodeC")
	nodeCzmqPort := flag.Int("nodeCzmqPort", 5557, "ZMQ Port of nodeC")
	flag.Parse()
	log.Debug().Str("nodeAurl", *nodeAurl).Str("nodeBurl", *nodeBurl).Str("nodeCurl", *nodeCurl).Msg("Using These URLs to connect to Node API")
	log.Debug().Str("nodeAzmq", string(*nodeAzmqPort)+":"+*nodeAzmqIp).Str("nodebzmq", string(*nodeBzmqPort)+":"+*nodeBzmqIp).Str("nodeCzmq", string(*nodeCzmqPort)+":"+*nodeCzmqIp).Msgf("Using these as ZMQ connection Points")

	nodeA := pkg.NewNode(*nodeAurl)

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
	nodeA.RegisterReplica("nodeB", *nodeBzmqIp, *nodeBzmqPort, 200, true)

	logNodeAction(nodeA, "Telling nodeA about nodeC")
	nodeA.RegisterReplica("nodeC", *nodeCzmqIp, *nodeCzmqPort, 200, true)

	logNodeAction(nodeA, "Getting all Replicas that nodeA has")
	nodeA.GetAllReplica(200, true)
	// Test sending data between nodes
}

func logNodeAction(node *pkg.Node, action string) {
	log.Info().Str("node", node.Url).Msg(action)
}

func logNodeFaliure(node *pkg.Node, expected, result string) {
	log.Warn().Str("node", node.Url).Msgf("expected: %s, but got: %#v", expected, result)
}

func logNodeFaliureMap(node *pkg.Node, expected string, result map[string]string) {
	log.Warn().Str("node", node.Url).Msgf("expected: %s, but got: %#v", expected, result)
}

func logDebugInfo(node *pkg.Node, info string) {
	log.Debug().Str("node", node.Url).Msg(info)
}

func expectEmptyResponse(resp map[string]string) {
	if resp != nil {
		log.Warn().Msgf("Expected empty response, but got: %#v", resp)
	}
}
