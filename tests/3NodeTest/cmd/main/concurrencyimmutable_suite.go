package main

import (
	"fmt"
	"strconv"

	"git.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg/grpcclient"
)

type ConcurrencyImmutableSuite struct {
	c *Config
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
		v, _ := nodes[0].GetItem(keygroup, key, false)
		found := 0

		for i := 0; i < concurrent; i++ {
			val, ok := expected[i][key]

			if !ok {
				continue
			}

			found++

			if val != v[0] {
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

func (t *ConcurrencyImmutableSuite) Name() string {
	return "ConcurrencyImmutable"
}

func (t *ConcurrencyImmutableSuite) RunTests() {
	run := 0
	// Test 1: create immutable keygroup, have two goroutines append data
	// expected behavior: all updates arrive
	run++
	concurrentUpdatesImmutable([]*grpcclient.Node{t.c.nodeA}, 2, 500, run)

	// Test 2: create immutable keygroup, have 100 goroutines append data
	// expected behavior: all updates arrive
	run++
	concurrentUpdatesImmutable([]*grpcclient.Node{t.c.nodeA}, 100, 50, run)

	// Test 3: create immutable keygroup on two nodes, have one goroutine each append data
	// expected behavior: all updates arrive
	run++
	concurrentUpdatesImmutable([]*grpcclient.Node{t.c.nodeA, t.c.nodeB}, 2, 100, run)
}

func NewConcurrencyImmutableSuite(c *Config) *ConcurrencyImmutableSuite {
	return &ConcurrencyImmutableSuite{
		c: c,
	}
}
