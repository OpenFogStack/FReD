package main

import (
	"fmt"
	"strconv"
	"time"

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

	expected := make([]map[uint64]string, concurrent)
	done := make(chan struct{})
	for i := 0; i < concurrent; i++ {
		expected[i] = make(map[uint64]string)
		go func(i int, node *grpcclient.Node, keygroup string, expected *map[uint64]string) {
			for j := 0; j < updates; j++ {
				val := randStringBytes(10)
				id := uint64(time.Now().UnixNano()) + uint64(i)
				_ = node.AppendItem(keygroup, id, val, false)
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
			if _, ok := keys[key]; ok {
				logNodeFailure(nodes[0], fmt.Sprintf("only one client can write to key %d", key), fmt.Sprintf("several clients were able to write to same id for %d", key))
				continue
			}

			keys[key] = val
		}
	}

	for key, val := range keys {
		k := strconv.FormatUint(key, 10)

		v, _ := nodes[0].GetItem(keygroup, k, false)

		if len(v) != 1 {
			logNodeFailure(nodes[0], fmt.Sprintf("expected one value for %d", key), fmt.Sprintf("got no %d return values", len(v)))
			continue
		}

		if v[0] != val {
			logNodeFailure(nodes[0], fmt.Sprintf("value %s for %d", val, key), fmt.Sprintf("got wrong value %s", v[0]))
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
