package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"git.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg/grpcclient"
)

type ConcurrencySuite struct {
	c *Config
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go#31832326
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func concurrentUpdates(nodes []*grpcclient.Node, concurrent int, updates int, run int) {

	if len(nodes) < 1 {
		return
	}

	keygroup := fmt.Sprintf("concurrencyTest%d", run)
	logNodeAction(nodes[0], "Create keygroup %s", keygroup)
	nodes[0].CreateKeygroup(keygroup, true, 0, false)

	for i, n := range nodes {
		if i == 0 {
			continue
		}
		logNodeAction(n, "adding node as replica for %s", keygroup)
		nodes[0].AddKeygroupReplica(keygroup, n.ID, 0, false)
	}

	expected := make([]map[string]string, concurrent)
	done := make(chan struct{})
	for i := 0; i < concurrent; i++ {
		expected[i] = make(map[string]string)
		go func(node *grpcclient.Node, keygroup string, expected *map[string]string) {
			for j := 0; j < updates; j++ {
				key := randStringBytes(2)
				val := randStringBytes(10)
				node.PutItem(keygroup, key, val, false)
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
			v := nodes[0].GetItem(keygroup, key, false)

			if v == "" {
				logNodeFailure(nodes[0], fmt.Sprintf("expected value %s for %s", val, key), "got no return value")
				continue
			}

			if v == val {
				// ok!
				continue
			}

			// hm, our returned value isn't the same as it should be - let's check the other maps
			found := false
			for j := 0; j < concurrent; j++ {
				if i == j {
					continue
				}

				jVal, ok := expected[j][key]
				if !ok {
					continue
				}

				if jVal == v {
					found = true
					break
				}
			}
			if !found {
				logNodeFailure(nodes[0], fmt.Sprintf("value for %s", key), fmt.Sprintf("got wrong value %s", v))
			}
		}
	}
}

func concurrentUpdatesImmutable(nodes []*grpcclient.Node, concurrent int, updates int, run int) {

	if len(nodes) < 1 {
		return
	}

	keygroup := fmt.Sprintf("concurrencyTestImmutable%d", run)
	logNodeAction(nodes[0], "Create keygroup %s", keygroup)
	nodes[0].CreateKeygroup(keygroup, false, 0, false)

	for i, n := range nodes {
		if i == 0 {
			continue
		}
		logNodeAction(n, "adding node as replica for %s", keygroup)
		nodes[0].AddKeygroupReplica(keygroup, n.ID, 0, false)
	}

	expected := make([]map[string]string, concurrent)
	done := make(chan struct{})
	for i := 0; i < concurrent; i++ {
		expected[i] = make(map[string]string)
		go func(node *grpcclient.Node, keygroup string, expected *map[string]string) {
			for j := 0; j < updates; j++ {
				val := randStringBytes(10)
				key := node.AppendItem(keygroup, val, false)
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
	// in this case we expect the keys to be in the range [0,updates[ and no same key to be in different maps

	for k := 0; k < updates; k++ {
		// got to convert to string first
		key := strconv.Itoa(k)
		v := nodes[0].GetItem(keygroup, key, false)
		found := 0

		for i := 0; i < concurrent; i++ {
			val, ok := expected[i][key]

			if !ok {
				continue
			}

			found++

			if val != v {
				logNodeFailure(nodes[0], fmt.Sprintf("value for %s", key), fmt.Sprintf("got wrong value %s", v))
			}
		}

		if found == 0 {
			logNodeFailure(nodes[0], fmt.Sprintf("expected value for %s", key), "got no return value")
			continue
		}

		if found > 1 {
			logNodeFailure(nodes[0], fmt.Sprintf("only one client can write to key %s", key), fmt.Sprintf("%d clients were able to write to same id for %s", found, key))
			continue
		}
	}
}

func (t *ConcurrencySuite) Name() string {
	return "Concurrency"
}

func (t *ConcurrencySuite) RunTests() {
	run := 0
	// Test 1: create a keygroup on a node, have it updated by two concurrent goroutines
	// expected behavior: all updates arrive
	concurrentUpdates([]*grpcclient.Node{t.c.nodeA}, 2, 1000, run)

	// Test 2: create a keygroup on a node, have it updated by 100 concurrent goroutines
	// expected behavior: all updates arrive
	run++
	concurrentUpdates([]*grpcclient.Node{t.c.nodeA}, 100, 100, run)

	// Test 3: create a keygroup on two nodes, have one goroutine update data at each node
	// expected behavior: all updates arrive, both nodes have the same data
	run++
	concurrentUpdates([]*grpcclient.Node{t.c.nodeA, t.c.nodeB}, 2, 1000, run)

	// Test 4: create immutable keygroup, have two goroutines append data
	// expected behavior: all updates arrive
	run++
	concurrentUpdatesImmutable([]*grpcclient.Node{t.c.nodeA}, 2, 500, run)

	// Test 5: create immutable keygroup, have 100 goroutines append data
	// expected behavior: all updates arrive
	run++
	concurrentUpdatesImmutable([]*grpcclient.Node{t.c.nodeA}, 100, 50, run)

	// Test 6: create immutable keygroup on two nodes, have one goroutine each append data
	// expected behavior: all updates arrive
	run++
	concurrentUpdatesImmutable([]*grpcclient.Node{t.c.nodeA, t.c.nodeB}, 2, 100, run)

}

func NewConcurrencySuite(c *Config) *ConcurrencySuite {
	return &ConcurrencySuite{
		c: c,
	}
}
