package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg/grpcclient"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Wait for the user to press enter to continue
	waitUser bool
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
	//useTLS := *flag.Bool("useTLS", false, "Use TLS (HTTPS instead of HTTP)")
	waitUser = *flag.Bool("wait-user", false, "wait for user input after each test")

	//apiVersion := *flag.String("apiVersion", "v0", "API Version (e.g. v0)")

	nodeAhost := *flag.String("nodeAhost", "172.26.0.10", "host of nodeA (e.g. localhost)") // Docker: localhost
	nodeAhttpPort := *flag.String("nodeAhttp", "9001", "port of nodeA (e.g. 9001)")         // Docker: 9002
	//nodeAzmqhost := *flag.String("nodeAzmqhost", "172.26.0.10", "host of nodeA (e.g. localhost) that can be reached by the other nodes") // Docker: 172.26.0.10
	//nodeAzmqPort := *flag.Int("nodeAzmqPort", 5555, "ZMQ Port of nodeA")
	nodeAzmqID := *flag.String("nodeAzmqID", "nodeA", "ZMQ Id of nodeA")

	nodeBhost := *flag.String("nodeBhost", "172.26.0.11", "host of nodeB (e.g. localhost)")
	nodeBhttpPort := *flag.String("nodeBhttp", "9001", "port of nodeB (e.g. 9001)")
	//nodeBzmqhost := *flag.String("nodeBzmqhost", "172.26.0.11", "host of nodeB (e.g. localhost) that can be reached by the other nodes")
	//nodeBzmqPort := *flag.Int("nodeBzmqPort", 5555, "ZMQ Port of nodeB")
	nodeBzmqID := *flag.String("nodeBzmqID", "nodeB", "ZMQ Id of nodeB")

	nodeChost := *flag.String("nodeChost", "172.26.0.12", "host of nodeC (e.g. localhost)")
	nodeChttpPort := *flag.String("nodeChttp", "9001", "port of nodeC (e.g. 9001)")
	//nodeCzmqhost := *flag.String("nodeCzmqhost", "172.26.0.12", "host of nodeC (e.g. localhost) that can be reached by the other nodes")
	//nodeCzmqPort := *flag.Int("nodeCzmqPort", 5555, "ZMQ Port of nodeC")
	nodeCzmqID := *flag.String("nodeCzmqID", "nodeC", "ZMQ Id of nodeC")

	triggerNodeHost := *flag.String("triggerNodeHost", "172.26.0.30:3333", "host of trigger node (e.g. localhost:3333)")
	triggerNodeWSHost := *flag.String("triggerNodeWSHost", "172.26.0.30:80", "host of trigger node web server (e.g. localhost:80)")
	triggerNodeID := *flag.String("triggerNodeID", "triggernode", "Id of trigger node")

	flag.Parse()

	//var protocol string
	//
	//if useTLS {
	//	protocol = "https"
	//} else {
	//	protocol = "http"
	//}

	//nodeAurl := fmt.Sprintf("%s://%s:%s/%s", protocol, nodeAhost, nodeAhttpPort, apiVersion)
	//log.Debug().Msgf("Node A: %s with ZMQ Port %d and ID %s", nodeAurl, nodeAzmqPort, nodeAzmqID)
	//
	//nodeBurl := fmt.Sprintf("%s://%s:%s/%s", protocol, nodeBhost, nodeBhttpPort, apiVersion)
	//log.Debug().Msgf("Node B: %s with ZMQ Port %d and ID %s", nodeBurl, nodeBzmqPort, nodeBzmqID)
	//
	//nodeCurl := fmt.Sprintf("%s://%s:%s/%s", protocol, nodeChost, nodeChttpPort, apiVersion)
	//log.Debug().Msgf("Node C: %s with ZMQ Port %d and ID %s", nodeCurl, nodeCzmqPort, nodeCzmqID)

	port, _ := strconv.Atoi(nodeAhttpPort)
	nodeA := grpcclient.NewNode(nodeAhost, port)
	port, _ = strconv.Atoi(nodeBhttpPort)
	nodeB := grpcclient.NewNode(nodeBhost, port)
	port, _ = strconv.Atoi(nodeChttpPort)
	nodeC := grpcclient.NewNode(nodeChost, port)

	time.Sleep(15 * time.Second)

	// Test Keygroups
	logNodeAction(nodeA, "Creating keygroup testing")
	nodeA.CreateKeygroup("testing", true, false)

	logNodeAction(nodeA, "Deleting keygroup testing")
	nodeA.DeleteKeygroup("testing", false)

	logNodeAction(nodeA, "Deleting nonexistent keygroup")
	nodeA.DeleteKeygroup("trololololo", true)

	logNodeAction(nodeA, "Creating Keygroup KG1")
	nodeA.CreateKeygroup("KG1", true, false)

	// Test Get/Put of a single node
	logNodeAction(nodeA, "Putting KG1-Item/KG1-Value into KG1")
	nodeA.PutItem("KG1", "KG1-Item", "KG1-Value", false)

	logNodeAction(nodeA, "Getting the value in KG1")

	resp := nodeA.GetItem("KG1", "KG1-Item", false)

	if resp != "KG1-Value" {
		logNodeFailure(nodeA, "resp is \"KG1-Value\"", resp)
	}

	logNodeAction(nodeA, "Getting a Value from a nonexistent keygroup")
	nodeA.GetItem("trololoool", "wut", true)

	logNodeAction(nodeA, "Putting a Value into a nonexistent keygroup")
	nodeA.PutItem("nonexistentkeygroup", "item", "data", true)

	logNodeAction(nodeA, "Putting new value KG1-Item/KG1-Value2 into KG1")
	nodeA.PutItem("KG1", "KG1-Item", "KG1-Value2", false)

	logNodeAction(nodeA, "Getting the value in KG1")
	resp = nodeA.GetItem("KG1", "KG1-Item", false)
	if resp != "KG1-Value2" {
		logNodeFailure(nodeA, "resp is \"KG1-Value2\"", resp)
	} else {
		logDebugInfo(nodeA, "Got "+resp)
	}

	logNodeAction(nodeA, "Getting all Replicas that nodeA has")
	parsed := nodeA.GetAllReplica(false)
	// Example Response: ["nodeB", "nodeC"]
	// Test for nodeA

	if len(parsed) != 3 {
		logNodeFailure(nodeA, "len(parsed) == 3", fmt.Sprintf("%d", len(parsed)))
	}

	// sorry but i still love go
	check := (len(parsed) != 3) && func() bool {
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
		logNodeFailure(nodeA, "parsed == ["+nodeBzmqID+", "+nodeCzmqID+","+nodeAzmqID+"]", fmt.Sprintf("%#v", parsed))
	}

	// Fun with replicas
	logNodeAction(nodeA, "Adding nodeB as Replica node for KG1")
	nodeA.AddKeygroupReplica("KG1", nodeBzmqID, false)

	logNodeAction(nodeB, "Deleting the value from KG1")
	nodeB.DeleteItem("KG1", "KG1-Item", false)

	logNodeAction(nodeB, "Getting the deleted value in KG1")
	_ = nodeB.GetItem("KG1", "KG1-Item", true)

	// Test sending data between nodes
	logNodeAction(nodeB, "Creating a new Keygroup (KGN) in nodeB, setting nodeA as Replica node")
	nodeB.CreateKeygroup("KGN", true, false)
	nodeB.AddKeygroupReplica("KGN", nodeAzmqID, false)

	logNodeAction(nodeB, "Putting something in KGN on nodeB, testing whether nodeA gets Replica (sleep 1.5s in between)")
	nodeB.PutItem("KGN", "Item", "Value", false)
	time.Sleep(1500 * time.Millisecond)
	resp = nodeA.GetItem("KGN", "Item", false)
	if resp != "Value" {
		logNodeFailure(nodeA, "resp is \"Value\"", resp)
	}

	logNodeAction(nodeA, "Putting something in KGN on nodeA, testing whether nodeB gets Replica (sleep 1.5s in between)")
	nodeA.PutItem("KGN", "Item2", "Value2", false)
	time.Sleep(1500 * time.Millisecond)
	resp = nodeB.GetItem("KGN", "Item2", false)
	if resp != "Value2" {
		logNodeFailure(nodeA, "resp is \"Value2\"", resp)
	}

	logNodeAction(nodeA, "Adding a replica for a nonexisting Keygroup")
	nodeA.AddKeygroupReplica("trololololo", nodeBzmqID, true)

	logNodeAction(nodeC, "Creating an already existing keygroup with another node")
	nodeC.CreateKeygroup("KGN", true, true)

	logNodeAction(nodeC, "Telling a node that is not part of the keygroup that it is now part of that keygroup")
	nodeC.AddKeygroupReplica("KGN", nodeCzmqID, false)

	logNodeAction(nodeA, "Creating a new Keygroup (kgall) with all three nodes as replica")
	nodeA.CreateKeygroup("kgall", true, false)
	nodeA.AddKeygroupReplica("kgall", nodeBzmqID, false)
	nodeB.AddKeygroupReplica("kgall", nodeCzmqID, false)

	logNodeAction(nodeC, "... sending data to one node, checking whether all nodes get the replica (sleep 1.5s)")
	nodeC.PutItem("kgall", "item", "value", false)
	time.Sleep(1500 * time.Millisecond)
	respA := nodeA.GetItem("kgall", "item", false)
	respB := nodeB.GetItem("kgall", "item", false)

	if respA != "value" || respB != "value" {
		logNodeFailure(nodeA, "both nodes respond with 'value'", fmt.Sprintf("NodeA: %s, NodeB: %s", respA, respB))
	}

	logNodeAction(nodeB, "...removing node from the keygroup all and checking whether it still has the data (sleep 1.5s)")
	nodeB.DeleteKeygroupReplica("kgall", nodeBzmqID, false)
	time.Sleep(1500 * time.Millisecond)
	respB = nodeB.GetItem("kgall", "item", true)

	logNodeAction(nodeB, fmt.Sprintf("Got Response %s", respB))

	logNodeAction(nodeB, "...re-adding the node to the keygroup all and checking whether it now gets the data (sleep 1.5s)")
	nodeA.AddKeygroupReplica("kgall", nodeBzmqID, false)
	time.Sleep(1500 * time.Millisecond)
	respA = nodeA.GetItem("kgall", "item", false)

	if respA != "value" {
		logNodeFailure(nodeA, "resp is \"value\"", resp)
	}

	respB = nodeB.GetItem("kgall", "item", false)

	if respB != "value" {
		logNodeFailure(nodeB, "resp is \"value\"", resp)
	}

	// let's test trigger nodes
	// create a new keygroup on nodeA
	logNodeAction(nodeA, "Creating keygroup triggertesting")
	nodeA.CreateKeygroup("triggertesting", true, false)

	logNodeAction(nodeA, "Creating keygroup nottriggertesting")
	nodeA.CreateKeygroup("nottriggertesting", true, false)

	//post an item1 to new keygroup
	logNodeAction(nodeA, "post an item1 to new keygroup triggertesting")
	nodeA.PutItem("triggertesting", "item1", "value1", false)
	//add trigger node to nodeA
	logNodeAction(nodeA, "add trigger node to nodeA for keygroup triggertesting")
	nodeA.AddKeygroupTrigger("triggertesting", triggerNodeID, triggerNodeHost, false)
	//post another item2 to new keygroup
	logNodeAction(nodeA, "post another item2 to new keygroup triggertesting")
	nodeA.PutItem("triggertesting", "item2", "value2", false)
	//delete item1 from keygroup
	logNodeAction(nodeA, "delete item1 from keygroup triggertesting")
	nodeA.DeleteItem("triggertesting", "item1", false)
	// post an item3 to keygroup nottriggertesting that should not be sent to trigger node
	logNodeAction(nodeA, "post an item3 to keygroup nottriggertesting that should not be sent to trigger node")
	nodeA.PutItem("nottriggertesting", "item3", "value3", false)
	//add keygroup to nodeB as well
	logNodeAction(nodeA, "add keygroup triggertesting to nodeB as well")
	nodeA.AddKeygroupReplica("triggertesting", nodeBzmqID, false)
	//post item4 to nodeB
	logNodeAction(nodeB, "post item4 to nodeB in keygroup triggertesting")
	nodeB.PutItem("triggertesting", "item4", "value4", false)
	//remove trigger node from nodeA
	logNodeAction(nodeA, "remove trigger node from nodeA in keygroup triggertesting")
	nodeA.DeleteKeygroupTrigger("triggertesting", triggerNodeID, false)
	//post item5 to nodeA
	logNodeAction(nodeA, "post item5 to nodeA in keygroup triggertesting")
	nodeA.PutItem("triggertesting", "item5", "value5", false)
	// check logs from trigger node
	// we should have the following logs (and nothing else):
	// put triggertesting item2 value2
	// del triggertesting item1
	// put triggertesting item4 value4
	logNodeAction(nodeA, "Checking if triggers have been executed correctly")
	checkTriggerNode(triggerNodeID, triggerNodeWSHost)
	//finally delete the keygroups again
	logNodeAction(nodeA, "deleting keygroup triggertesting")
	nodeA.DeleteKeygroup("triggertesting", false)
	logNodeAction(nodeA, "deleting keygroup nottriggertesting")
	nodeA.DeleteKeygroup("nottriggertesting", false)
	// try to get the trigger nodes now
	logNodeAction(nodeA, "try to get the trigger nodes for keygroup triggertesting after deletion")
	nodeA.GetKeygroupTriggers("triggertesting", true)

	// testing immutable keygroups
	// create immutable keygroup on nodeB
	nodeB.CreateKeygroup("log", false, false)
	// create item to keygroup -> should work
	nodeB.PutItem("log", "test-item", "value-1", false)
	// update item in keygroup -> should not work
	nodeB.PutItem("log", "test-item", "value-2", true)
	// get item from keygroup -> should return appended item
	respB = nodeB.GetItem("log", "test-item", false)

	if respB != "value-1" {
		logNodeFailure(nodeB, "resp is not value-1", respB)
	}
	// delete item from keygroup -> should not work
	nodeB.DeleteItem("log", "test-item", true)
	// add nodeC as replica node
	nodeB.AddKeygroupReplica("log", nodeCzmqID, false)
	// update item on nodeC -> should not work
	nodeB.PutItem("log", "test-item", "value-3", true)

	// testing immutable keygroups
	// create immutable keygroup on nodeB
	nodeB.CreateKeygroup("log", false, 200, true)
	// create item to keygroup -> should work
	nodeB.PutItem("log", "test-item", "value-1", 200, true)
	// update item in keygroup -> should not work
	nodeB.PutItem("log", "test-item", "value-2", 400, false)
	// get item from keygroup -> should return appended item
	respB = nodeB.GetItem("log", "test-item", 200, false)
	if respB != "value-1" {
		logNodeFailure(nodeB, "resp is not value-1 but %s", respB)
	}
	// delete item from keygroup -> should not work
	nodeB.DeleteItem("log", "test-item", 400, false)
	// add nodeC as replica node
	nodeB.AddKeygroupReplica("log", nodeCzmqID, 200, false)
	// update item on nodeC -> should not work
	nodeB.PutItem("log", "test-item", "value-3", 400, false)

	totalerrors := nodeA.Errors + nodeB.Errors + nodeC.Errors

	if totalerrors > 0 {
		log.Error().Msgf("Total Errors: %d", totalerrors)
	}

	os.Exit(totalerrors)
}

func logNodeAction(node *grpcclient.Node, action string) {
	wait()
	log.Info().Str("node", node.Addr).Msg(action)
}

func logNodeFailure(node *grpcclient.Node, expected, result string) {
	wait()
	log.Warn().Str("node", node.Addr).Msgf("expected: %s, but got: %#v", expected, result)
}

func logDebugInfo(node *grpcclient.Node, info string) {
	wait()
	log.Debug().Str("node", node.Addr).Msg(info)
}

func checkTriggerNode(triggerNodeID, triggerNodeWSHost string) {
	log.Info().Str("trigger node", triggerNodeWSHost).Msg("Checking Trigger Node logs")

	type LogEntry struct {
		Op  string `json:"op"`
		Kg  string `json:"kg"`
		ID  string `json:"id"`
		Val string `json:"val"`
	}

	// put triggertesting item2 value2
	// del triggertesting item1
	// put triggertesting item4 value4

	expected := make([]LogEntry, 3)
	expected[0] = LogEntry{
		Op:  "put",
		Kg:  "triggertesting",
		ID:  "item2",
		Val: "value2",
	}

	expected[1] = LogEntry{
		Op: "del",
		Kg: "triggertesting",
		ID: "item1",
	}

	expected[2] = LogEntry{
		Op:  "put",
		Kg:  "triggertesting",
		ID:  "item4",
		Val: "value4",
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/", triggerNodeWSHost))

	if err != nil {
		log.Warn().Str("trigger node", triggerNodeWSHost).Msgf("%#v", err)
		return
	}

	var result []LogEntry
	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Warn().Str("trigger node", triggerNodeWSHost).Msgf("%#v", err)
		return
	}

	err = resp.Body.Close()

	if err != nil {
		log.Warn().Str("trigger node", triggerNodeWSHost).Msgf("%#v", err)
		return
	}

	if len(result) != len(expected) {
		log.Warn().Str("trigger node", triggerNodeID).Msgf("expected: %s, but got: %#v", expected, result)
		return
	}

	for i := range expected {
		if expected[i] != result[i] {
			log.Warn().Str("trigger node", triggerNodeID).Msgf("expected: %s, but got: %#v", expected[i], result[i])
			return
		}
	}
}

func wait() {
	if waitUser {
		log.Info().Msg("Please press enter to continue:")
		_, _, _ = reader.ReadLine()
	} else {
		time.Sleep(1 * time.Second)
	}
}
