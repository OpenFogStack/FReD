package main

import (
	"fmt"
	"time"
)

type ReplicaSuite struct {
	c *Config
}

func (t *ReplicaSuite) Name() string {
	return "Replication"
}

func (t *ReplicaSuite) RunTests() {
	// Fun with replicas
	logNodeAction(t.c.nodeA, "Create keygroup KGRep")
	t.c.nodeA.CreateKeygroup("KGRep", true, 0, false)

	logNodeAction(t.c.nodeA, "Adding nodeB as Replica node for KGRep")
	t.c.nodeA.AddKeygroupReplica("KGRep", t.c.nodeB.ID, 0, false)

	logNodeAction(t.c.nodeB, "Putting a valuein KGRep")
	t.c.nodeB.PutItem("KGRep", "KGRepItem", "val", false)

	logNodeAction(t.c.nodeB, "Deleting the value from KGRep")
	t.c.nodeB.DeleteItem("KGRep", "KGRepItem", false)

	logNodeAction(t.c.nodeB, "Getting the deleted value in KGRep")
	_ = t.c.nodeB.GetItem("KGRep", "KGRepItem", true)

	// Test sending data between nodes
	logNodeAction(t.c.nodeB, "Creating a new Keygroup (KGN) in nodeB, setting nodeA as Replica node")
	t.c.nodeB.CreateKeygroup("KGN", true, 0, false)
	t.c.nodeB.AddKeygroupReplica("KGN", t.c.nodeA.ID, 0, false)

	logNodeAction(t.c.nodeB, "Putting something in KGN on nodeB, testing whether nodeA gets Replica (sleep 1.5s in between)")
	t.c.nodeB.PutItem("KGN", "Item", "Value", false)
	time.Sleep(1500 * time.Millisecond)
	resp := t.c.nodeA.GetItem("KGN", "Item", false)
	if resp != "Value" {
		logNodeFailure(t.c.nodeA, "resp is \"Value\"", resp)
	}

	logNodeAction(t.c.nodeA, "Putting something in KGN on nodeA, testing whether nodeB gets Replica (sleep 1.5s in between)")
	t.c.nodeA.PutItem("KGN", "Item2", "Value2", false)
	time.Sleep(1500 * time.Millisecond)
	resp = t.c.nodeB.GetItem("KGN", "Item2", false)
	if resp != "Value2" {
		logNodeFailure(t.c.nodeA, "resp is \"Value2\"", resp)
	}

	logNodeAction(t.c.nodeA, "Adding a replica for a nonexisting Keygroup")
	t.c.nodeA.AddKeygroupReplica("trololololo", t.c.nodeB.ID, 0, true)

	logNodeAction(t.c.nodeC, "Creating an already existing keygroup with another node")
	t.c.nodeC.CreateKeygroup("KGN", true, 0, true)

	logNodeAction(t.c.nodeC, "Telling a node that is not part of the keygroup that it is now part of that keygroup")
	t.c.nodeC.AddKeygroupReplica("KGN", t.c.nodeC.ID, 0, false)

	logNodeAction(t.c.nodeA, "Creating a new Keygroup (kgall) with all three nodes as replica")
	t.c.nodeA.CreateKeygroup("kgall", true, 0, false)
	t.c.nodeA.AddKeygroupReplica("kgall", t.c.nodeB.ID, 0, false)
	t.c.nodeB.AddKeygroupReplica("kgall", t.c.nodeC.ID, 0, false)

	logNodeAction(t.c.nodeC, "... sending data to one node, checking whether all nodes get the replica (sleep 1.5s)")
	t.c.nodeC.PutItem("kgall", "item", "value", false)
	time.Sleep(1500 * time.Millisecond)
	respA := t.c.nodeA.GetItem("kgall", "item", false)
	respB := t.c.nodeB.GetItem("kgall", "item", false)

	if respA != "value" || respB != "value" {
		logNodeFailure(t.c.nodeA, "both nodes respond with 'value'", fmt.Sprintf("NodeA: %s, NodeB: %s", respA, respB))
	}

	logNodeAction(t.c.nodeB, "...removing node from the keygroup all and checking whether it still has the data (sleep 1.5s)")
	t.c.nodeB.DeleteKeygroupReplica("kgall", t.c.nodeB.ID, false)
	time.Sleep(1500 * time.Millisecond)
	respB = t.c.nodeB.GetItem("kgall", "item", true)

	logNodeAction(t.c.nodeB, fmt.Sprintf("Got Response %s", respB))

	logNodeAction(t.c.nodeB, "...re-adding the node to the keygroup all and checking whether it now gets the data (sleep 1.5s)")
	t.c.nodeA.AddKeygroupReplica("kgall", t.c.nodeB.ID, 0, false)
	time.Sleep(1500 * time.Millisecond)
	respA = t.c.nodeA.GetItem("kgall", "item", false)

	if respA != "value" {
		logNodeFailure(t.c.nodeA, "resp is \"value\"", resp)
	}

	respB = t.c.nodeB.GetItem("kgall", "item", false)

	if respB != "value" {
		logNodeFailure(t.c.nodeB, "resp is \"value\"", resp)
	}

	// delete the last node from a keygroup
	logNodeAction(t.c.nodeA, "Preparing to delete all members from a keygroup...")
	t.c.nodeA.CreateKeygroup("deletetest", true, 0, false)
	t.c.nodeA.PutItem("deletetest", "item", "value", false)
	t.c.nodeA.AddKeygroupReplica("deletetest", t.c.nodeB.ID, 0, false)
	t.c.nodeA.DeleteKeygroupReplica("deletetest", t.c.nodeA.ID, false)
	// NodeB is the only replica left
	logNodeAction(t.c.nodeB, "Removing last member of a keygroup delete-test")
	t.c.nodeB.DeleteKeygroupReplica("deletetest", t.c.nodeB.ID, true)
}

func NewReplicaSuite(c *Config) *ReplicaSuite {
	return &ReplicaSuite{
		c: c,
	}
}
