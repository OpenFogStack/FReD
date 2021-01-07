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
	waitUser = *flag.Bool("wait-user", false, "wait for user input after each test")

	nodeAhost := *flag.String("nodeAhost", "172.26.1.1", "host of nodeA (e.g. localhost)") // Docker: localhost
	nodeAhttpPort := *flag.String("nodeAhttp", "9001", "port of nodeA (e.g. 9001)")        // Docker: 9002
	nodeApeeringID := *flag.String("nodeAzmqID", "nodeA", "ZMQ Id of nodeA")

	nodeBhost := *flag.String("nodeBhost", "172.26.2.1", "host of nodeB (e.g. localhost)")
	nodeBhttpPort := *flag.String("nodeBhttp", "9001", "port of nodeB (e.g. 9001)")
	nodeBpeeringID := *flag.String("nodeBzmqID", "nodeB", "ZMQ Id of nodeB")

	nodeChost := *flag.String("nodeChost", "172.26.3.1", "host of nodeC (e.g. localhost)")
	nodeChttpPort := *flag.String("nodeChttp", "9001", "port of nodeC (e.g. 9001)")
	nodeCpeeringID := *flag.String("nodeCzmqID", "nodeC", "ZMQ Id of nodeC")

	triggerNodeHost := *flag.String("triggerNodeHost", "172.26.5.1:3333", "host of trigger node (e.g. localhost:3333)")
	triggerNodeWSHost := *flag.String("triggerNodeWSHost", "172.26.5.1:80", "host of trigger node web server (e.g. localhost:80)")
	triggerNodeID := *flag.String("triggerNodeID", "triggernode", "Id of trigger node")

	certFile := *flag.String("cert-file", "/cert/client.crt", "Certificate to talk to FReD")
	keyFile := *flag.String("key-file", "/cert/client.key", "Keyfile to talk to FReD")

	littleCertFile := *flag.String("little-cert-file", "/cert/littleclient.crt", "Certificate to talk to FReD as \"littleclient\"")
	littleKeyFile := *flag.String("little-key-file", "/cert/littleclient.key", "Keyfile to talk to FReD as \"littleclient\"")

	flag.Parse()

	port, _ := strconv.Atoi(nodeAhttpPort)
	nodeA := grpcclient.NewNode(nodeAhost, port, certFile, keyFile)
	port, _ = strconv.Atoi(nodeBhttpPort)
	nodeB := grpcclient.NewNode(nodeBhost, port, certFile, keyFile)
	port, _ = strconv.Atoi(nodeChttpPort)
	nodeC := grpcclient.NewNode(nodeChost, port, certFile, keyFile)

	time.Sleep(15 * time.Second)

	// Test Keygroups
	logNodeAction(nodeA, "Creating keygroup testing")
	nodeA.CreateKeygroup("testing", true, 0, false)

	logNodeAction(nodeA, "Deleting keygroup testing")
	nodeA.DeleteKeygroup("testing", false)

	logNodeAction(nodeA, "Deleting nonexistent keygroup")
	nodeA.DeleteKeygroup("trololololo", true)

	logNodeAction(nodeA, "Creating Keygroup KG1")
	nodeA.CreateKeygroup("KG1", true, 0, false)

	// Test Get/Put of a single node
	logNodeAction(nodeA, "Putting KG1Item/KG1Value into KG1")
	nodeA.PutItem("KG1", "KG1Item", "KG1Value", false)

	logNodeAction(nodeA, "Getting the value in KG1")

	resp := nodeA.GetItem("KG1", "KG1Item", false)

	if resp != "KG1Value" {
		logNodeFailure(nodeA, "resp is \"KG1Value\"", resp)
	}

	logNodeAction(nodeA, "Getting a Value from a nonexistent keygroup")
	nodeA.GetItem("trololoool", "wut", true)

	logNodeAction(nodeA, "Putting a Value into a nonexistent keygroup")
	nodeA.PutItem("nonexistentkeygroup", "item", "data", true)

	logNodeAction(nodeA, "Putting new value KG1Item/KG1Value2 into KG1")
	nodeA.PutItem("KG1", "KG1Item", "KG1Value2", false)

	logNodeAction(nodeA, "Getting the value in KG1")
	resp = nodeA.GetItem("KG1", "KG1Item", false)
	if resp != "KG1Value2" {
		logNodeFailure(nodeA, "resp is \"KG1Value2\"", resp)
	}

	logNodeAction(nodeA, "Getting all Replicas that nodeA has")
	parsed := nodeA.GetAllReplica(false)
	// Example Response: map[string]string
	// {"nodeA": "1.2.3.4:5000", "nodeB": "4.5.6.7:4000"}
	// Test for nodeA

	if len(parsed) != 3 {
		logNodeFailure(nodeA, "GetAllReplica returns 3 nodes", fmt.Sprintf("%d", len(parsed)))
	}

	addr, ok := parsed[nodeApeeringID]
	if !ok {
		logNodeFailure(nodeA, "GetAllReplica response contains nodeA", "nodeA not found")
	} else if addr != fmt.Sprintf("%s:%s", nodeAhost, nodeAhttpPort) {
		logNodeFailure(nodeA, "nodeA address is "+fmt.Sprintf("%s:%s", nodeAhost, nodeAhttpPort), addr)
	}

	addr, ok = parsed[nodeBpeeringID]
	if !ok {
		logNodeFailure(nodeA, "GetAllReplica response contains nodeB", "nodeB not found")
	} else if addr != fmt.Sprintf("%s:%s", nodeBhost, nodeBhttpPort) {
		logNodeFailure(nodeA, "nodeB address is "+fmt.Sprintf("%s:%s", nodeBhost, nodeBhttpPort), addr)
	}

	addr, ok = parsed[nodeCpeeringID]
	if !ok {
		logNodeFailure(nodeA, "GetAllReplica response contains nodeC", "nodeC not found")
	} else if addr != fmt.Sprintf("%s:%s", nodeChost, nodeChttpPort) {
		logNodeFailure(nodeA, "nodeC address is "+fmt.Sprintf("%s:%s", nodeChost, nodeChttpPort), addr)
	}

	// Fun with replicas
	logNodeAction(nodeA, "Adding nodeB as Replica node for KG1")
	nodeA.AddKeygroupReplica("KG1", nodeBpeeringID, 0, false)

	logNodeAction(nodeB, "Deleting the value from KG1")
	nodeB.DeleteItem("KG1", "KG1Item", false)

	logNodeAction(nodeB, "Getting the deleted value in KG1")
	_ = nodeB.GetItem("KG1", "KG1Item", true)

	// Test sending data between nodes
	logNodeAction(nodeB, "Creating a new Keygroup (KGN) in nodeB, setting nodeA as Replica node")
	nodeB.CreateKeygroup("KGN", true, 0, false)
	nodeB.AddKeygroupReplica("KGN", nodeApeeringID, 0, false)

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
	nodeA.AddKeygroupReplica("trololololo", nodeBpeeringID, 0, true)

	logNodeAction(nodeC, "Creating an already existing keygroup with another node")
	nodeC.CreateKeygroup("KGN", true, 0, true)

	logNodeAction(nodeC, "Telling a node that is not part of the keygroup that it is now part of that keygroup")
	nodeC.AddKeygroupReplica("KGN", nodeCpeeringID, 0, false)

	logNodeAction(nodeA, "Creating a new Keygroup (kgall) with all three nodes as replica")
	nodeA.CreateKeygroup("kgall", true, 0, false)
	nodeA.AddKeygroupReplica("kgall", nodeBpeeringID, 0, false)
	nodeB.AddKeygroupReplica("kgall", nodeCpeeringID, 0, false)

	logNodeAction(nodeC, "... sending data to one node, checking whether all nodes get the replica (sleep 1.5s)")
	nodeC.PutItem("kgall", "item", "value", false)
	time.Sleep(1500 * time.Millisecond)
	respA := nodeA.GetItem("kgall", "item", false)
	respB := nodeB.GetItem("kgall", "item", false)

	if respA != "value" || respB != "value" {
		logNodeFailure(nodeA, "both nodes respond with 'value'", fmt.Sprintf("NodeA: %s, NodeB: %s", respA, respB))
	}

	logNodeAction(nodeB, "...removing node from the keygroup all and checking whether it still has the data (sleep 1.5s)")
	nodeB.DeleteKeygroupReplica("kgall", nodeBpeeringID, false)
	time.Sleep(1500 * time.Millisecond)
	respB = nodeB.GetItem("kgall", "item", true)

	logNodeAction(nodeB, fmt.Sprintf("Got Response %s", respB))

	logNodeAction(nodeB, "...re-adding the node to the keygroup all and checking whether it now gets the data (sleep 1.5s)")
	nodeA.AddKeygroupReplica("kgall", nodeBpeeringID, 0, false)
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
	nodeA.CreateKeygroup("triggertesting", true, 0, false)

	logNodeAction(nodeA, "Creating keygroup nottriggertesting")
	nodeA.CreateKeygroup("nottriggertesting", true, 0, false)

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
	nodeA.AddKeygroupReplica("triggertesting", nodeBpeeringID, 0, false)
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
	logNodeAction(nodeA, "deleting keygroup triggertesting")
	nodeA.DeleteKeygroup("triggertesting", false)

	logNodeAction(nodeA, "deleting keygroup nottriggertesting")
	nodeA.DeleteKeygroup("nottriggertesting", false)

	logNodeAction(nodeA, "try to get the trigger nodes for keygroup triggertesting after deletion")
	nodeA.GetKeygroupTriggers("triggertesting", true)

	// testing immutable keygroups
	logNodeAction(nodeB, "Testing immutable keygroups by creating a new immutable keygroup on nodeB")
	nodeB.CreateKeygroup("log", false, 0, false)

	logNodeAction(nodeB, "Creating an item in this keygroup")
	res := nodeB.AppendItem("log", "value1", false)

	if res != "0" {
		logNodeFailure(nodeB, "0", res)
	}

	logNodeAction(nodeB, "Updating an item in this keygroup")
	nodeB.PutItem("log", res, "value-2", true)

	logNodeAction(nodeB, "Getting updated item from immutable keygroup")
	respB = nodeB.GetItem("log", res, false)

	if respB != "value1" {
		logNodeFailure(nodeB, "resp is value1", respB)
	}

	logNodeAction(nodeB, "Deleting an item in immutable keygroup")
	nodeB.DeleteItem("log", "testitem", true)

	logNodeAction(nodeB, "Adding nodeC as replica to immutable keygroup")
	nodeB.AddKeygroupReplica("log", nodeCpeeringID, 0, false)

	logNodeAction(nodeC, "Updating immutable item on other nodeC")
	nodeC.PutItem("log", "testitem", "value-3", true)

	// TODO is this the right place???
	logNodeAction(nodeC, "Appending another item to readonly log.")
	res = nodeC.AppendItem("log", "value-2", false)

	if res != "1" {
		logNodeFailure(nodeB, "1", res)
	}

	// test expiring data items
	logNodeAction(nodeC, "Create normal keygroup on nodeC without expiry")
	nodeC.CreateKeygroup("expirytest", true, 0, false)

	logNodeAction(nodeC, "Add nodeA as replica with expiry of 5s")
	nodeC.AddKeygroupReplica("expirytest", nodeApeeringID, 5, false)

	logNodeAction(nodeC, "Update something on nodeC")
	nodeC.PutItem("expirytest", "test", "test", false)

	logNodeAction(nodeA, "Test whether nodeA has received the update. Wait 5s and check that it is not there anymore")
	nodeA.GetItem("expirytest", "test", false)
	time.Sleep(5 * time.Second)
	nodeA.GetItem("expirytest", "test", true)

	logNodeAction(nodeC, "....the item should still exist with nodeC")
	nodeC.GetItem("expirytest", "test", false)

	// testing adding a node as a replica for a keygroup on itself
	logNodeAction(nodeB, "Create and populate a new keygroup to test pulling")
	nodeB.CreateKeygroup("pulltest", true, 0, false)
	nodeB.PutItem("pulltest", "item1", "val1", false)
	nodeB.PutItem("pulltest", "item2", "val2", false)

	logNodeAction(nodeA, "add nodeA as a replica to that keygroup and see if it pulls the needed data on its own (sleep 3s)")
	nodeA.AddKeygroupReplica("pulltest", nodeApeeringID, 0, false)
	time.Sleep(3 * time.Second)
	// check if the items exist
	if res := nodeA.GetItem("pulltest", "item1", false); res != "val1" {
		logNodeFailure(nodeA, "val1", res)
	}
	if res := nodeA.GetItem("pulltest", "item2", false); res != "val2" {
		logNodeFailure(nodeA, "val2", res)
	}

	logNodeAction(nodeA, "Add an item on nodeA, check wheter it populates to nodeB")
	nodeA.PutItem("pulltest", "item3", "val3", false)
	// check if nodeB also gets that item
	if res := nodeB.GetItem("pulltest", "item3", false); res != "val3" {
		logNodeFailure(nodeB, "val3", res)
	}

	// test RBAC and authentication
	logNodeAction(nodeA, "create keygroup \"rbactest\"")
	nodeA.CreateKeygroup("rbactest", true, 0, false)

	logNodeAction(nodeA, "put item into keygroup \"rbactest\"")
	nodeA.PutItem("rbactest", "item1", "value1", false)

	logNodeAction(nodeA, "add little client as read only to rbac test")
	nodeA.AddUser("littleclient", "rbactest", "ReadKeygroup", false)

	logNodeAction(nodeA, "try to read with little client -> should work")
	littleClient := grpcclient.NewNode(nodeAhost, port, littleCertFile, littleKeyFile)
	if val := littleClient.GetItem("rbactest", "item1", false); val != "value1" {
		logNodeFailure(nodeA, "value1", val)
	}

	logNodeAction(nodeA, "try to write with little client -> should not work")
	littleClient.PutItem("rbactest", "item1", "value2", true)

	logNodeAction(nodeA, "add role configure replica to little client -> should work")
	nodeA.AddUser("littleclient", "rbactest", "ConfigureReplica", false)

	logNodeAction(nodeA, "add replica nodeB to keygroup with little client -> should work")
	littleClient.AddKeygroupReplica("rbactest", nodeBpeeringID, 0, false)

	logNodeAction(nodeB, "remove permission to read from keygroup -> should work")
	nodeB.RemoveUser("littleclient", "rbactest", "ReadKeygroup", false)

	logNodeAction(nodeA, "try to read from keygroup with little client -> should not work")
	if val := littleClient.GetItem("rbactest", "item1", true); val != "" {
		logNodeFailure(nodeA, "", val)
	}

	totalerrors := nodeA.Errors + nodeB.Errors + nodeC.Errors + littleClient.Errors

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
	node.Errors++
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
