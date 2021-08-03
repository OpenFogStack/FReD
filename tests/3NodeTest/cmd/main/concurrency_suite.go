package main

import (
	"fmt"

	"git.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg/grpcclient"
)

type ConcurrencySuite struct {
	c *Config
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
				key := randStringBytes(1)
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
			v, versions := nodes[0].GetItem(keygroup, key, false)

			if len(v) == 0 {
				logNodeFailure(nodes[0], fmt.Sprintf("expected value %s for %s", val, key), "got no return value")
				continue
			}

			if len(v) > concurrent {
				logNodeFailure(nodes[0], fmt.Sprintf("expected value %s for %s", val, key), fmt.Sprintf("got %d return values: %#v %#v", len(v), v, versions))
				continue
			}

			for j := range v {
				if v[j] == val {
					// ok!
					continue
				}

				// hm, our returned value isn't the same as it should be - let's check the other maps
				found := false
				var possibleVals []string
				for l := 0; l < concurrent; l++ {

					jVal, ok := expected[l][key]
					if !ok {
						continue
					}

					possibleVals = append(possibleVals, jVal)

					if jVal == v[j] {
						found = true
						break
					}
				}

				if !found {
					logNodeFailure(nodes[0], fmt.Sprintf("one of values %#v for %s", possibleVals, key), fmt.Sprintf("got wrong value %#v", v))
				}
			}
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
	//concurrentUpdates([]*grpcclient.Node{t.c.nodeA}, 2, 1000, run)

	// Test 2: create a keygroup on a node, have it updated by 100 concurrent goroutines
	// expected behavior: all updates arrive
	run++
	//concurrentUpdates([]*grpcclient.Node{t.c.nodeA}, 100, 100, run)

	// Test 3: create a keygroup on two nodes, have one goroutine update data at each node
	// expected behavior: all updates arrive, both nodes have the same data
	run++
	concurrentUpdates([]*grpcclient.Node{t.c.nodeA, t.c.nodeB}, 2, 10, run)

}

func NewConcurrencySuite(c *Config) *ConcurrencySuite {
	return &ConcurrencySuite{
		c: c,
	}
}
