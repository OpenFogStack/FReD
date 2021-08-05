package threenodetest

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

const (
	composePath  = "../runner/"
	certBasePath = "../runner/certificates/"
)

var nodes map[string]*node

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go#31832326
// removing capital letters so we can get more conflicts
const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestMain(m *testing.M) {
	dc := testcontainers.NewLocalDockerCompose([]string{
		composePath + "fredwork.yml",
		composePath + "etcd.yml",
		composePath + "nodeA.yml",
		composePath + "nodeB.yml",
		composePath + "nodeC.yml",
		composePath + "trigger.yml",
	}, "3nodetest")

	execError := dc.Down()

	if execError.Error != nil {
		panic(fmt.Sprintf("Failed when running: %v", execError.Command))
	}

	execError = dc.WithCommand([]string{"build", "--quiet"}).Invoke()

	if execError.Error != nil {
		panic(fmt.Sprintf("Failed when running: %v", execError))
	}

	// dc.WithExposedService("nodeB", 9001, wait.NewLogStrategy("FReD Node is operational!"))

	execError = dc.WithCommand([]string{"up", "-d", "--quiet-pull", "--force-recreate", "--renew-anon-volumes", "--remove-orphans"}).Invoke()

	if execError.Error != nil {
		panic(fmt.Sprintf("Failed when running: %v", execError))
	}

	nodes = map[string]*node{
		"nodeA": newNode("localhost", 9001, "172.26.1.1:9001", "nodeA", certBasePath+"client.crt", certBasePath+"client.key", certBasePath+"ca.crt"),
		"nodeB": newNode("localhost", 9002, "172.26.2.1:9001", "nodeB", certBasePath+"client.crt", certBasePath+"client.key", certBasePath+"ca.crt"),
		"nodeC": newNode("localhost", 9003, "172.26.3.1:9001", "nodeC", certBasePath+"client.crt", certBasePath+"client.key", certBasePath+"ca.crt"),
	}

	// TODO: health check
	// port health check requires /bin/sh in the container
	time.Sleep(20 * time.Second)

	stat := m.Run()

	for _, n := range nodes {
		n.close()
	}

	execError = dc.Down()

	if execError.Error != nil {
		panic(fmt.Sprintf("Failed when running: %v", execError.Command))
	}

	os.Exit(stat)
}

func TestStandard(t *testing.T) {
	err := nodes["nodeA"].createKeygroup("testing", true, 0)
	assert.NoError(t, err)

	err = nodes["nodeA"].deleteKeygroup("testing")
	assert.NoError(t, err)

	err = nodes["nodeB"].deleteKeygroup("trololololo")
	assert.Error(t, err)

	err = nodes["nodeA"].createKeygroup("KG1", true, 0)
	assert.NoError(t, err)

	// Test Get/Put of a single node

	_, err = nodes["nodeA"].putItem("KG1", "KG1Item", "KG1Value")
	assert.NoError(t, err)

	vals, _, err := nodes["nodeA"].getItem("KG1", "KG1Item")
	assert.NoError(t, err)

	assert.Len(t, vals, 1)
	assert.Equal(t, "KG1Value", vals[0])

	_, _, err = nodes["nodeA"].getItem("trololoool", "wut")
	assert.Error(t, err)

	_, err = nodes["nodeA"].putItem("nonexistentkeygroup", "item", "data")
	assert.Error(t, err)

	_, err = nodes["nodeA"].putItem("KG1", "KG1Item", "KG1Value2")
	assert.NoError(t, err)

	vals, _, err = nodes["nodeA"].getItem("KG1", "KG1Item")
	assert.NoError(t, err)
	assert.Len(t, vals, 1)
	assert.Equal(t, "KG1Value2", vals[0])
}
func TestScan(t *testing.T) {
	err := nodes["nodeA"].createKeygroup("scantest", true, 0)
	assert.NoError(t, err)
	numItems := 20
	scanStart := 5
	scanRange := 10
	// 2. put in a bunch of items
	ids := make([]string, numItems)
	data := make([]string, numItems)

	for i := 0; i < 20; i++ {
		data[i] = "val" + strconv.Itoa(i)
		ids[i] = "id" + strconv.Itoa(i)
		_, err = nodes["nodeA"].putItem("scantest", ids[i], data[i])
		assert.NoError(t, err)
	}

	// 3. do a scan read
	// we expect [scanRange] amount of items, starting with [scanStart]
	startKey := "id" + strconv.Itoa(scanStart)

	items, err := nodes["nodeA"].scanItems("scantest", startKey, uint64(scanRange))
	assert.NoError(t, err)

	expected := scanRange - scanStart

	assert.Len(t, items, expected)

	for i := scanStart; i < scanStart+expected; i++ {
		assert.Contains(t, items, ids[i])
		assert.Equal(t, data[i], items[ids[i]])
	}
}

func TestReplica(t *testing.T) {
	parsed, err := nodes["nodeA"].getAllReplica()
	assert.NoError(t, err)
	// Example Response: map[string]string
	// {"nodeA": "1.2.3.4:5000", "nodeB": "4.5.6.7:4000"}
	// Test for nodeA
	assert.Len(t, parsed, 3)
	assert.Contains(t, parsed, "nodeA")
	assert.Contains(t, parsed, "nodeB")
	assert.Contains(t, parsed, "nodeC")
	assert.Equal(t, nodes["nodeA"].addr, parsed["nodeA"])
	assert.Equal(t, nodes["nodeB"].addr, parsed["nodeB"])
	assert.Equal(t, nodes["nodeC"].addr, parsed["nodeC"])

	// Fun with replicas
	err = nodes["nodeA"].createKeygroup("KGRep", true, 0)
	assert.NoError(t, err)

	err = nodes["nodeA"].addKeygroupReplica("KGRep", nodes["nodeB"].id, 0)
	assert.NoError(t, err)

	_, err = nodes["nodeB"].putItem("KGRep", "KGRepItem", "val")
	assert.NoError(t, err)

	err = nodes["nodeB"].deleteItem("KGRep", "KGRepItem")
	assert.NoError(t, err)

	_, _, err = nodes["nodeB"].getItem("KGRep", "KGRepItem")
	assert.Error(t, err)

	// Test sending data between nodes

	err = nodes["nodeB"].createKeygroup("KGN", true, 0)
	assert.NoError(t, err)
	err = nodes["nodeB"].addKeygroupReplica("KGN", nodes["nodeA"].id, 0)
	assert.NoError(t, err)

	_, err = nodes["nodeB"].putItem("KGN", "Item", "Value")
	assert.NoError(t, err)
	time.Sleep(1500 * time.Millisecond)
	resp, _, err := nodes["nodeA"].getItem("KGN", "Item")
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, resp[0], "Value")

	_, err = nodes["nodeA"].putItem("KGN", "Item2", "Value2")
	assert.NoError(t, err)
	time.Sleep(1500 * time.Millisecond)
	resp, _, err = nodes["nodeB"].getItem("KGN", "Item2")
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, resp[0], "Value2")

	err = nodes["nodeA"].addKeygroupReplica("trololololo", nodes["nodeB"].id, 0)
	assert.Error(t, err)

	err = nodes["nodeC"].createKeygroup("KGN", true, 0)
	assert.Error(t, err)

	err = nodes["nodeC"].addKeygroupReplica("KGN", nodes["nodeC"].id, 0)
	assert.NoError(t, err)

	err = nodes["nodeA"].createKeygroup("kgall", true, 0)
	assert.NoError(t, err)
	err = nodes["nodeA"].addKeygroupReplica("kgall", nodes["nodeB"].id, 0)
	assert.NoError(t, err)
	err = nodes["nodeB"].addKeygroupReplica("kgall", nodes["nodeC"].id, 0)
	assert.NoError(t, err)

	_, err = nodes["nodeC"].putItem("kgall", "item", "value")
	assert.NoError(t, err)
	time.Sleep(1500 * time.Millisecond)

	respA, _, err := nodes["nodeA"].getItem("kgall", "item")
	assert.NoError(t, err)

	respB, _, err := nodes["nodeB"].getItem("kgall", "item")
	assert.NoError(t, err)

	assert.Len(t, respA, 1)
	assert.Equal(t, respA[0], "value")

	assert.Len(t, respB, 1)
	assert.Equal(t, respB[0], "value")

	err = nodes["nodeB"].deleteKeygroupReplica("kgall", nodes["nodeB"].id)
	assert.NoError(t, err)

	time.Sleep(1500 * time.Millisecond)

	_, _, err = nodes["nodeB"].getItem("kgall", "item")

	assert.Error(t, err)

	err = nodes["nodeA"].addKeygroupReplica("kgall", nodes["nodeB"].id, 0)
	assert.NoError(t, err)
	time.Sleep(1500 * time.Millisecond)
	respA, _, err = nodes["nodeA"].getItem("kgall", "item")
	assert.NoError(t, err)

	assert.Len(t, respA, 1)
	assert.Equal(t, respA[0], "value")

	respB, _, err = nodes["nodeB"].getItem("kgall", "item")
	assert.NoError(t, err)

	assert.Len(t, respB, 1)
	assert.Equal(t, respB[0], "value")

	// delete the last node from a keygroup

	err = nodes["nodeA"].createKeygroup("deletetest", true, 0)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	_, err = nodes["nodeA"].putItem("deletetest", "item", "value")
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	err = nodes["nodeA"].addKeygroupReplica("deletetest", nodes["nodeB"].id, 0)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	err = nodes["nodeA"].deleteKeygroupReplica("deletetest", nodes["nodeA"].id)
	assert.NoError(t, err)
	// NodeB is the only replica left

	time.Sleep(1 * time.Second)
	err = nodes["nodeB"].deleteKeygroupReplica("deletetest", nodes["nodeB"].id)
	assert.Error(t, err)

	// checking "GetReplicas" and "getKeygroupReplica"

	_, host, err := nodes["nodeB"].getReplica(nodes["nodeA"].id)
	assert.NoError(t, err)

	assert.Equal(t, nodes["nodeA"].addr, host)

	err = nodes["nodeB"].createKeygroup("replicanodetest", true, 0)
	assert.NoError(t, err)
	err = nodes["nodeB"].addKeygroupReplica("replicanodetest", nodes["nodeA"].id, 10)
	assert.NoError(t, err)
	hosts, err := nodes["nodeA"].getKeygroupReplica("replicanodetest")
	assert.NoError(t, err)
	assert.Len(t, hosts, 2)

	assert.Contains(t, hosts, "nodeA")
	assert.Contains(t, hosts, "nodeB")
	assert.Equal(t, 10, hosts["nodeA"])
	assert.Equal(t, 0, hosts["nodeB"])
}

func checkTriggerNode(t *testing.T, triggerNodeWSHost string) {
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

	assert.NoError(t, err)

	var result []LogEntry
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)

	err = resp.Body.Close()

	assert.NoError(t, err)

	assert.Len(t, result, len(expected))

	assert.ElementsMatch(t, result, expected)
}

func TestTrigger(t *testing.T) {
	triggerNodeID := "triggerA"
	triggerNodeHost := "172.26.5.1:3333"
	triggerNodeWSHost := "localhost:8000"

	// let's test trigger nodes
	// create a new keygroup on nodeA

	err := nodes["nodeA"].createKeygroup("triggertesting", true, 0)
	assert.NoError(t, err)

	err = nodes["nodeA"].createKeygroup("nottriggertesting", true, 0)
	assert.NoError(t, err)

	//post an item1 to new keygroup

	_, err = nodes["nodeA"].putItem("triggertesting", "item1", "value1")
	assert.NoError(t, err)
	//add trigger node to nodeA

	err = nodes["nodeA"].addKeygroupTrigger("triggertesting", triggerNodeID, triggerNodeHost)
	assert.NoError(t, err)
	//post another item2 to new keygroup

	_, err = nodes["nodeA"].putItem("triggertesting", "item2", "value2")
	assert.NoError(t, err)
	//delete item1 from keygroup

	err = nodes["nodeA"].deleteItem("triggertesting", "item1")
	assert.NoError(t, err)
	// post an item3 to keygroup nottriggertesting that should not be sent to trigger node

	_, err = nodes["nodeA"].putItem("nottriggertesting", "item3", "value3")
	assert.NoError(t, err)
	//add keygroup to nodeB as well

	err = nodes["nodeA"].addKeygroupReplica("triggertesting", nodes["nodeB"].id, 0)
	assert.NoError(t, err)
	//post item4 to nodeB

	_, err = nodes["nodeB"].putItem("triggertesting", "item4", "value4")
	assert.NoError(t, err)
	//remove trigger node from nodeA

	err = nodes["nodeA"].deleteKeygroupTrigger("triggertesting", triggerNodeID)
	assert.NoError(t, err)
	//post item5 to nodeA

	_, err = nodes["nodeA"].putItem("triggertesting", "item5", "value5")
	assert.NoError(t, err)
	// check logs from trigger node
	// we should have the following logs (and nothing else):
	// put triggertesting item2 value2
	// del triggertesting item1
	// put triggertesting item4 value4

	checkTriggerNode(t, triggerNodeWSHost)

	err = nodes["nodeA"].deleteKeygroup("triggertesting")
	assert.NoError(t, err)

	err = nodes["nodeA"].deleteKeygroup("nottriggertesting")
	assert.NoError(t, err)

	_, err = nodes["nodeA"].getKeygroupTriggers("triggertesting")
	assert.Error(t, err)
}

func TestImmutable(t *testing.T) {
	// testing immutable keygroups

	err := nodes["nodeB"].createKeygroup("log", false, 0)
	assert.NoError(t, err)

	res, err := nodes["nodeB"].appendItem("log", 0, "value1")
	assert.NoError(t, err)

	assert.Equal(t, "0", res)

	_, err = nodes["nodeB"].putItem("log", res, "value-2")
	assert.Error(t, err)

	respB, _, err := nodes["nodeB"].getItem("log", res)
	assert.NoError(t, err)

	assert.Len(t, respB, 1)
	assert.Equal(t, "value1", respB[0])

	err = nodes["nodeB"].deleteItem("log", "0")
	assert.Error(t, err)

	err = nodes["nodeB"].addKeygroupReplica("log", nodes["nodeC"].id, 0)
	assert.NoError(t, err)

	_, err = nodes["nodeC"].putItem("log", "0", "value-3")
	assert.Error(t, err)

	res, err = nodes["nodeC"].appendItem("log", 1, "value-4")
	assert.NoError(t, err)

	assert.Equal(t, "1", res)

	_, err = nodes["nodeC"].appendItem("log", 1, "value-5")
	assert.Error(t, err)
}

func TestExpiry(t *testing.T) {
	// test expiring data items

	err := nodes["nodeC"].createKeygroup("expirytest", true, 0)
	assert.NoError(t, err)

	err = nodes["nodeC"].addKeygroupReplica("expirytest", nodes["nodeA"].id, 5)
	assert.NoError(t, err)

	_, err = nodes["nodeC"].putItem("expirytest", "test", "test")
	assert.NoError(t, err)

	_, _, err = nodes["nodeA"].getItem("expirytest", "test")
	assert.NoError(t, err)
	time.Sleep(5 * time.Second)
	_, _, err = nodes["nodeA"].getItem("expirytest", "test")
	assert.Error(t, err)

	_, _, err = nodes["nodeC"].getItem("expirytest", "test")
	assert.NoError(t, err)
}

func TestSelfReplica(t *testing.T) {
	// testing adding a node as a replica for a keygroup on itself

	err := nodes["nodeB"].createKeygroup("pulltest", true, 0)
	assert.NoError(t, err)
	_, err = nodes["nodeB"].putItem("pulltest", "item1", "val1")
	assert.NoError(t, err)
	_, err = nodes["nodeB"].putItem("pulltest", "item2", "val2")
	assert.NoError(t, err)

	err = nodes["nodeA"].addKeygroupReplica("pulltest", nodes["nodeA"].id, 0)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
	// check if the items exist
	res, _, err := nodes["nodeA"].getItem("pulltest", "item1")

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "val1", res[0])

	res, _, err = nodes["nodeA"].getItem("pulltest", "item2")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "val2", res[0])

	_, err = nodes["nodeA"].putItem("pulltest", "item3", "val3")
	assert.NoError(t, err)

	// check if nodeB also gets that item
	res, _, err = nodes["nodeB"].getItem("pulltest", "item3")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "val3", res[0])
}

func TestAuth(t *testing.T) {
	// test RBAC and authentication

	littleClient := newNode("localhost", 9001, "172.26.1.1:9001", "nodeA", certBasePath+"littleclient.crt", certBasePath+"littleclient.key", certBasePath+"ca.crt")

	err := nodes["nodeA"].createKeygroup("rbactest", true, 0)
	assert.NoError(t, err)

	_, err = nodes["nodeA"].putItem("rbactest", "item1", "value1")
	assert.NoError(t, err)

	err = nodes["nodeA"].addUser("littleclient", "rbactest", "ReadKeygroup")
	assert.NoError(t, err)

	val, _, err := littleClient.getItem("rbactest", "item1")

	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Len(t, val, 1)
	assert.Equal(t, "value1", val[0])

	_, err = littleClient.putItem("rbactest", "item1", "value2")
	assert.Error(t, err)

	err = nodes["nodeA"].addUser("littleclient", "rbactest", "ConfigureReplica")
	assert.NoError(t, err)

	err = littleClient.addKeygroupReplica("rbactest", nodes["nodeB"].id, 0)
	assert.NoError(t, err)

	err = nodes["nodeB"].removeUser("littleclient", "rbactest", "ReadKeygroup")
	assert.NoError(t, err)

	_, _, _ = littleClient.getItem("rbactest", "item1")
	assert.Error(t, err)
}

func concurrentUpdates(t *testing.T, nodes []*node, concurrent int, updates int, run int) {

	if len(nodes) < 1 {
		return
	}

	keygroup := fmt.Sprintf("concurrencyTest%d", run)

	err := nodes[0].createKeygroup(keygroup, true, 0)
	assert.NoError(t, err)

	for i, n := range nodes {
		if i == 0 {
			continue
		}

		err = nodes[0].addKeygroupReplica(keygroup, n.id, 0)
		assert.NoError(t, err)
	}

	expected := make([]map[string]string, concurrent)
	done := make(chan struct{})
	for i := 0; i < concurrent; i++ {
		expected[i] = make(map[string]string)
		go func(node *node, keygroup string, expected *map[string]string) {
			for j := 0; j < updates; j++ {
				key := randStringBytes(1)
				val := randStringBytes(10)
				_, err := node.putItem(keygroup, key, val)
				assert.NoError(t, err)
				(*expected)[key] = val
			}
			done <- struct{}{}
		}(nodes[i%len(nodes)], keygroup, &expected[i])
	}

	// block until all goroutines have finished
	for i := 0; i < concurrent; i++ {
		<-done
	}

	// let's check if everything worked

	for i := 0; i < concurrent; i++ {
		for key, val := range expected[i] {
			v, _, err := nodes[0].getItem(keygroup, key)
			assert.NoError(t, err)

			assert.NotEmpty(t, v)
			assert.LessOrEqual(t, concurrent, len(v))

			for j := range v {
				if v[j] == val {
					// ok!
					continue
				}

				// hm, our returned value isn't the same as it should be - let's check the other maps
				found := false
				for l := 0; l < concurrent; l++ {

					jVal, ok := expected[l][key]
					if !ok {
						continue
					}

					if jVal == v[j] {
						found = true
						break
					}
				}

				assert.True(t, found)
			}
		}
	}
}

func TestConcurrency(t *testing.T) {
	run := 0
	// Test 1: create a keygroup on a node, have it updated by two concurrent goroutines
	// expected behavior: all updates arrive
	concurrentUpdates(t, []*node{nodes["nodeA"]}, 2, 1000, run)

	// Test 2: create a keygroup on a node, have it updated by 100 concurrent goroutines
	// expected behavior: all updates arrive
	run++
	concurrentUpdates(t, []*node{nodes["nodeA"]}, 100, 100, run)

	// Test 3: create a keygroup on two nodes, have one goroutine update data at each node
	// expected behavior: all updates arrive, both nodes have the same data
	run++
	concurrentUpdates(t, []*node{nodes["nodeA"], nodes["nodeB"]}, 2, 10, run)
}

func concurrentUpdatesImmutable(t *testing.T, nodes []*node, concurrent int, updates int, run int) {

	if len(nodes) < 1 {
		return
	}

	keygroup := fmt.Sprintf("concurrencyTestImmutable%d", run)

	err := nodes[0].createKeygroup(keygroup, false, 0)
	assert.NoError(t, err)

	for i, n := range nodes {
		if i == 0 {
			continue
		}

		err = nodes[0].addKeygroupReplica(keygroup, n.id, 0)
		assert.NoError(t, err)
	}

	expected := make([]map[uint64]string, concurrent)
	done := make(chan struct{})
	for i := 0; i < concurrent; i++ {
		expected[i] = make(map[uint64]string)
		go func(i int, node *node, keygroup string, expected *map[uint64]string) {
			for j := 0; j < updates; j++ {
				val := randStringBytes(10)
				id := uint64(time.Now().UnixNano()) + uint64(i)
				_, err = node.appendItem(keygroup, id, val)
				assert.NoError(t, err)
				(*expected)[id] = val
			}
			done <- struct{}{}
		}(i, nodes[i%len(nodes)], keygroup, &expected[i])
	}

	// block until all goroutines have finished
	for i := 0; i < concurrent; i++ {
		<-done
	}

	// let's check if everything worked
	// in this case no same key should be in in different maps
	keys := make(map[uint64]string)

	for i := 0; i < concurrent; i++ {
		for key, val := range expected[i] {
			assert.NotContains(t, keys, key)

			keys[key] = val
		}
	}

	for key, val := range keys {
		k := strconv.FormatUint(key, 10)

		v, _, err := nodes[0].getItem(keygroup, k)
		assert.NoError(t, err)

		assert.Len(t, v, 1)
		assert.Equal(t, val, v[0])
	}
}

func TestConcurrencyImmutable(t *testing.T) {
	run := 0
	// Test 1: create immutable keygroup, have two goroutines append data
	// expected behavior: all updates arrive
	run++
	concurrentUpdatesImmutable(t, []*node{nodes["nodeA"]}, 2, 500, run)

	// Test 2: create immutable keygroup, have 100 goroutines append data
	// expected behavior: all updates arrive
	run++
	concurrentUpdatesImmutable(t, []*node{nodes["nodeA"]}, 100, 50, run)

	// Test 3: create immutable keygroup on two nodes, have one goroutine each append data
	// expected behavior: all updates arrive
	run++
	concurrentUpdatesImmutable(t, []*node{nodes["nodeA"], nodes["nodeB"]}, 2, 100, run)

}

func versioningUpdates(t *testing.T, nodes []*node, concurrent int, updates int, run int) {

	if len(nodes) < 1 {
		return
	}

	keygroup := fmt.Sprintf("versioningTest%d", run)

	err := nodes[0].createKeygroup(keygroup, true, 0)
	assert.NoError(t, err)

	for i, n := range nodes {
		if i == 0 {
			continue
		}

		err = nodes[0].addKeygroupReplica(keygroup, n.id, 0)
		assert.NoError(t, err)
	}

	type item struct {
		val     string
		version vclock.VClock
	}

	expected := make([]map[string]item, concurrent)
	done := make(chan struct{})
	for i := 0; i < concurrent; i++ {
		expected[i] = make(map[string]item)
		go func(node *node, keygroup string, expected *map[string]item) {
			for j := 0; j < updates; j++ {
				key := randStringBytes(2)
				val := randStringBytes(10)
				v, err := node.putItem(keygroup, key, val)
				assert.NoError(t, err)
				(*expected)[key] = item{
					val:     val,
					version: v,
				}
			}
			done <- struct{}{}
		}(nodes[i%len(nodes)], keygroup, &expected[i])
	}

	// block until all goroutines have finished
	for i := 0; i < concurrent; i++ {
		<-done
	}

	// let's check if everything worked

	for i := 0; i < concurrent; i++ {
		for key, it := range expected[i] {
			v, versions, err := nodes[0].getItem(keygroup, key)
			assert.NoError(t, err)

			assert.NotEmpty(t, v)
			assert.LessOrEqual(t, len(v), concurrent)

			for j := range v {
				if v[j] == it.val {
					// ok!
					assert.True(t, versions[j].Compare(it.version, vclock.Equal))
					continue
				}

				// hm, our returned value isn't the same as it should be - let's check the other maps
				found := false
				for l := 0; l < concurrent; l++ {

					jVal, ok := expected[l][key]
					if !ok {
						continue
					}

					if jVal.val == v[j] {
						found = true
						assert.True(t, versions[j].Compare(it.version, vclock.Equal))
						break
					}
				}

				assert.True(t, found)
			}
		}
	}
}

func TestVersioning(t *testing.T) {
	// all of this doesnt work yet because we don't have a way to force concurrent updates...
	/*
				// a quick test first
				// create a new keygroup with replicas A and B
				kg := "versioningtestKG"

				err = nodes["nodeB"].createKeygroup(kg, true, 0)
		assert.NoError(t, err)

				err = nodes["nodeB"].addKeygroupReplica(kg, err = nodes["nodeA"].ID, 0)
		assert.NoError(t, err)
				// two clients concurrently update item "Item1", both should receive a version number
				wg := sync.WaitGroup{}
				wg.Add(2)
				v := make([]vclock.VClock, 2)
				start := make(chan struct{})

				go func() {

					start <- struct{}{}
					v[0] = err = nodes["nodeA"].putItem(kg, "Item1", "val1")
		assert.NoError(t, err)
					wg.Done()
				}()

				go func() {

					<-start
					v[1] = err = nodes["nodeB"].putItem(kg, "Item1", "val2")
		assert.NoError(t, err)
					wg.Done()
				}()

				wg.Wait()

				// one will be A:1 and one will be B:1
				vA := vclock.VClock{}
				vA.Tick(err = nodes["nodeA"].ID)
				vB := vclock.VClock{}
				vB.Tick(err = nodes["nodeB"].ID)
				if !v[0].Compare(vA, vclock.Equal) {
					logNodeFailure(err = nodes["nodeA"], vA.SortedVCString(), v[0].SortedVCString())
				}
				if !v[1].Compare(vB, vclock.Equal) {
					logNodeFailure(err = nodes["nodeB"], vA.SortedVCString(), v[1].SortedVCString())
				}

				// now try to update the version A:1 on B

				v[0] = err = nodes["nodeB"].putItemVersion(kg, "Item1", "val3", vA)
		assert.NoError(t, err)

				// return should be version number A:2
				vA.Tick(err = nodes["nodeA"].ID)
				if !v[0].Compare(vA, vclock.Equal) {
					logNodeFailure(err = nodes["nodeA"], vA.SortedVCString(), v[0].SortedVCString())
				}

				// a read should return both versions A:2 and B:1

				items, versions := err = nodes["nodeB"].getItem(kg, "Item1")
		assert.NoError(t, err)
				if len(items) != len(versions) {
					logNodeFailure(err = nodes["nodeB"], "equal number of values and versions", fmt.Sprintf("%d items and %d versions", len(items), len(versions)))
				}

				if len(items) != 2 {
					logNodeFailure(err = nodes["nodeB"], "2 versions and values", fmt.Sprintf("%d items and %d versions", len(items), len(versions)))
				}

				if items[0] != "val3" || items[1] != "val2" || !v[0].Compare(vA, vclock.Equal) || !v[1].Compare(vB, vclock.Equal) {
					logNodeFailure(err = nodes["nodeB"], fmt.Sprintf("val3 at %s and val2 at %s", vA.SortedVCString(), vB.SortedVCString()), fmt.Sprintf("%s at %s and %s at %s", items[0], items[1], v[0].SortedVCString(), v[1].SortedVCString()))
				}

				// now delete version B:1

				err = nodes["nodeA"].deleteItemVersion(kg, "Item1", vB)
		assert.NoError(t, err)
				// a read should return only A:2

				items, versions = err = nodes["nodeB"].getItem(kg, "Item1")
		assert.NoError(t, err)
				if len(items) != len(versions) {
					logNodeFailure(err = nodes["nodeB"], "equal number of values and versions", fmt.Sprintf("%d items and %d versions", len(items), len(versions)))
				}

				if len(items) != 1 {
					logNodeFailure(err = nodes["nodeB"], "1 versions and values", fmt.Sprintf("%d items and %d versions", len(items), len(versions)))
				}

				if items[0] != "val3" || !v[0].Compare(vA, vclock.Equal) {
					logNodeFailure(err = nodes["nodeB"], fmt.Sprintf("val3 at %s", vA.SortedVCString()), fmt.Sprintf("%s at %s", items[0], v[0].SortedVCString()))
				}**/

	// another quick test
	// create a keygroup
	kg := "versioningtestKG2"

	err := nodes["nodeB"].createKeygroup(kg, true, 0)
	assert.NoError(t, err)
	// add a node as a replica

	err = nodes["nodeB"].addKeygroupReplica(kg, nodes["nodeA"].id, 0)
	assert.NoError(t, err)

	// update an item

	v1, err := nodes["nodeA"].putItem(kg, "Item1", "val1")
	assert.NoError(t, err)
	// update another item

	_, err = nodes["nodeA"].putItemVersion(kg, "Item1", "val2", v1)
	assert.NoError(t, err)
	// try to delete with an older version -> should not work

	_, err = nodes["nodeA"].deleteItemVersion(kg, "Item1", v1)
	assert.Error(t, err)

	run := 0
	// Test 1: create a keygroup on a node, have it updated by two concurrent goroutines
	// expected behavior: all updates arrive
	versioningUpdates(t, []*node{nodes["nodeB"]}, 2, 1000, run)

	// Test 2: create a keygroup on a node, have it updated by 100 concurrent goroutines
	// expected behavior: all updates arrive
	run++
	versioningUpdates(t, []*node{nodes["nodeB"]}, 100, 100, run)

	// Test 3: create a keygroup on two nodes, have one goroutine update data at each node
	// expected behavior: all updates arrive, both nodes have the same data
	run++
	versioningUpdates(t, []*node{nodes["nodeB"], nodes["nodeC"]}, 2, 1000, run)

	// and now the same stuff for distributed nodes (nodeA in this case)
	run++
	versioningUpdates(t, []*node{nodes["nodeA"]}, 2, 1000, run)
	run++
	versioningUpdates(t, []*node{nodes["nodeA"]}, 100, 100, run)
	run++
	versioningUpdates(t, []*node{nodes["nodeA"], nodes["nodeB"]}, 2, 1000, run)

}
