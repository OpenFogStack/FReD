package main

import (
	"fmt"

	"git.tu-berlin.de/mcc-fred/fred/pkg/vector"
	"git.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg/grpcclient"
	"github.com/DistributedClocks/GoVector/govec/vclock"
)

type VersioningSuite struct {
	c *Config
}

func versioningUpdates(nodes []*grpcclient.Node, concurrent int, updates int, run int) {

	if len(nodes) < 1 {
		return
	}

	keygroup := fmt.Sprintf("versioningTest%d", run)
	logNodeAction(nodes[0], "Create keygroup %s", keygroup)
	nodes[0].CreateKeygroup(keygroup, true, 0, false)

	for i, n := range nodes {
		if i == 0 {
			continue
		}
		logNodeAction(n, "adding node as replica for %s", keygroup)
		nodes[0].AddKeygroupReplica(keygroup, n.ID, 0, false)
	}

	type item struct {
		val     string
		version vclock.VClock
	}

	expected := make([]map[string]item, concurrent)
	done := make(chan struct{})
	for i := 0; i < concurrent; i++ {
		expected[i] = make(map[string]item)
		go func(node *grpcclient.Node, keygroup string, expected *map[string]item) {
			for j := 0; j < updates; j++ {
				key := randStringBytes(2)
				val := randStringBytes(10)
				v := node.PutItem(keygroup, key, val, false)
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
			v, versions := nodes[0].GetItem(keygroup, key, false)

			if len(v) == 0 {
				logNodeFailure(nodes[0], fmt.Sprintf("expected value %s for %s", it.val, key), "got no return value")
				continue
			}

			if len(v) > concurrent {
				logNodeFailure(nodes[0], fmt.Sprintf("expected value %s for %s", it.val, key), fmt.Sprintf("got %d return values: %#v %#v", len(v), v, versions))
				continue
			}

			for j := range v {
				if v[j] == it.val {
					// ok!
					if !versions[j].Compare(it.version, vclock.Equal) {
						logNodeFailure(nodes[0], fmt.Sprintf("expected version %s for %s", vector.SortedVCString(it.version), key), fmt.Sprintf("got %d return values: %#v %#v", len(v), v, versions))
					}
					continue
				}

				// hm, our returned value isn't the same as it should be - let's check the other maps
				found := false
				var possibleVals []item
				for l := 0; l < concurrent; l++ {

					jVal, ok := expected[l][key]
					if !ok {
						continue
					}

					possibleVals = append(possibleVals, jVal)

					if jVal.val == v[j] {
						found = true
						if !versions[j].Compare(jVal.version, vclock.Equal) {
							logNodeFailure(nodes[0], fmt.Sprintf("expected version %s for %s", vector.SortedVCString(jVal.version), key), fmt.Sprintf("got %d return values: %#v %#v", len(v), v, versions))
						}
						break
					}
				}

				if !found {
					logNodeFailure(nodes[0], fmt.Sprintf("one of values %#v for %s", possibleVals, key), fmt.Sprintf("got wrong value(s) %#v %#v", v, versions))
				}
			}
		}
	}
}

func (t *VersioningSuite) Name() string {
	return "Versioning"
}

func (t *VersioningSuite) RunTests() {
	// all of this doesnt work yet because we don't have a way to force concurrent updates...
	/*
		// a quick test first
		// create a new keygroup with replicas A and B
		kg := "versioningtestKG"
		logNodeAction(t.c.nodeB, "create keygroup %s", kg)
		t.c.nodeB.CreateKeygroup(kg, true, 0, false)
		logNodeAction(t.c.nodeB, "adding nodeA as replica for keygroup %s", kg)
		t.c.nodeB.AddKeygroupReplica(kg, t.c.nodeA.ID, 0, false)
		// two clients concurrently update item "Item1", both should receive a version number
		wg := sync.WaitGroup{}
		wg.Add(2)
		v := make([]vclock.VClock, 2)
		start := make(chan struct{})

		go func() {
			logNodeAction(t.c.nodeA, "adding item 1 to keygroup %s", kg)
			start <- struct{}{}
			v[0] = t.c.nodeA.PutItem(kg, "Item1", "val1", false)
			wg.Done()
		}()

		go func() {
			logNodeAction(t.c.nodeB, "adding item 2 to keygroup %s", kg)
			<-start
			v[1] = t.c.nodeB.PutItem(kg, "Item1", "val2", false)
			wg.Done()
		}()

		wg.Wait()

		// one will be A:1 and one will be B:1
		vA := vclock.VClock{}
		vA.Tick(t.c.nodeA.ID)
		vB := vclock.VClock{}
		vB.Tick(t.c.nodeB.ID)
		if !v[0].Compare(vA, vclock.Equal) {
			logNodeFailure(t.c.nodeA, vA.SortedVCString(), v[0].SortedVCString())
		}
		if !v[1].Compare(vB, vclock.Equal) {
			logNodeFailure(t.c.nodeB, vA.SortedVCString(), v[1].SortedVCString())
		}

		// now try to update the version A:1 on B
		logNodeAction(t.c.nodeB, "updating item with specific version %s", vA.SortedVCString())
		v[0] = t.c.nodeB.PutItemVersion(kg, "Item1", "val3", vA, false)

		// return should be version number A:2
		vA.Tick(t.c.nodeA.ID)
		if !v[0].Compare(vA, vclock.Equal) {
			logNodeFailure(t.c.nodeA, vA.SortedVCString(), v[0].SortedVCString())
		}

		// a read should return both versions A:2 and B:1
		logNodeAction(t.c.nodeB, "reading all versions for item")
		items, versions := t.c.nodeB.GetItem(kg, "Item1", false)
		if len(items) != len(versions) {
			logNodeFailure(t.c.nodeB, "equal number of values and versions", fmt.Sprintf("%d items and %d versions", len(items), len(versions)))
		}

		if len(items) != 2 {
			logNodeFailure(t.c.nodeB, "2 versions and values", fmt.Sprintf("%d items and %d versions", len(items), len(versions)))
		}

		if items[0] != "val3" || items[1] != "val2" || !v[0].Compare(vA, vclock.Equal) || !v[1].Compare(vB, vclock.Equal) {
			logNodeFailure(t.c.nodeB, fmt.Sprintf("val3 at %s and val2 at %s", vA.SortedVCString(), vB.SortedVCString()), fmt.Sprintf("%s at %s and %s at %s", items[0], items[1], v[0].SortedVCString(), v[1].SortedVCString()))
		}

		// now delete version B:1
		logNodeAction(t.c.nodeA, "deleting a version for this item")
		t.c.nodeA.DeleteItemVersion(kg, "Item1", vB, false)
		// a read should return only A:2
		logNodeAction(t.c.nodeB, "reading all versions for item")
		items, versions = t.c.nodeB.GetItem(kg, "Item1", false)
		if len(items) != len(versions) {
			logNodeFailure(t.c.nodeB, "equal number of values and versions", fmt.Sprintf("%d items and %d versions", len(items), len(versions)))
		}

		if len(items) != 1 {
			logNodeFailure(t.c.nodeB, "1 versions and values", fmt.Sprintf("%d items and %d versions", len(items), len(versions)))
		}

		if items[0] != "val3" || !v[0].Compare(vA, vclock.Equal) {
			logNodeFailure(t.c.nodeB, fmt.Sprintf("val3 at %s", vA.SortedVCString()), fmt.Sprintf("%s at %s", items[0], v[0].SortedVCString()))
		}**/

	// another quick test
	// create a keygroup
	kg := "versioningtestKG2"
	logNodeAction(t.c.nodeB, "create keygroup %s", kg)
	t.c.nodeB.CreateKeygroup(kg, true, 0, false)
	// add a node as a replica
	logNodeAction(t.c.nodeB, "adding nodeA as replica for keygroup %s", kg)
	t.c.nodeB.AddKeygroupReplica(kg, t.c.nodeA.ID, 0, false)

	// update an item
	logNodeAction(t.c.nodeA, "adding item 1 to keygroup %s", kg)
	v1 := t.c.nodeA.PutItem(kg, "Item1", "val1", false)
	// update another item
	logNodeAction(t.c.nodeA, "updating item 1 in keygroup %s", kg)
	_ = t.c.nodeA.PutItemVersion(kg, "Item1", "val2", v1, false)
	// try to delete with an older version -> should not work
	logNodeAction(t.c.nodeA, "deleting item 1 in keygroup %s with an old version", kg)
	t.c.nodeA.DeleteItemVersion(kg, "Item1", v1, true)

	run := 0
	// Test 1: create a keygroup on a node, have it updated by two concurrent goroutines
	// expected behavior: all updates arrive
	versioningUpdates([]*grpcclient.Node{t.c.nodeB}, 2, 1000, run)

	// Test 2: create a keygroup on a node, have it updated by 100 concurrent goroutines
	// expected behavior: all updates arrive
	run++
	versioningUpdates([]*grpcclient.Node{t.c.nodeB}, 100, 100, run)

	// Test 3: create a keygroup on two nodes, have one goroutine update data at each node
	// expected behavior: all updates arrive, both nodes have the same data
	run++
	versioningUpdates([]*grpcclient.Node{t.c.nodeB, t.c.nodeC}, 2, 1000, run)

	// and now the same stuff for distributed nodes (nodeA in this case)
	run++
	versioningUpdates([]*grpcclient.Node{t.c.nodeA}, 2, 1000, run)
	run++
	versioningUpdates([]*grpcclient.Node{t.c.nodeA}, 100, 100, run)
	run++
	versioningUpdates([]*grpcclient.Node{t.c.nodeA, t.c.nodeB}, 2, 1000, run)
}

func NewVersioningSuite(c *Config) *VersioningSuite {
	return &VersioningSuite{
		c: c,
	}
}
