package main

import (
	"time"
)

type ExpirySuite struct {
	c *Config
}

func (t *ExpirySuite) RunTests() {
	// test expiring data items
	logNodeAction(t.c.nodeC, "Create normal keygroup on nodeC without expiry")
	t.c.nodeC.CreateKeygroup("expirytest", true, 0, false)

	logNodeAction(t.c.nodeC, "Add nodeA as replica with expiry of 5s")
	t.c.nodeC.AddKeygroupReplica("expirytest", t.c.nodeA.ID, 5, false)

	logNodeAction(t.c.nodeC, "Update something on nodeC")
	t.c.nodeC.PutItem("expirytest", "test", "test", false)

	logNodeAction(t.c.nodeA, "Test whether nodeA has received the update. Wait 5s and check that it is not there anymore")
	t.c.nodeA.GetItem("expirytest", "test", false)
	time.Sleep(5 * time.Second)
	t.c.nodeA.GetItem("expirytest", "test", true)

	logNodeAction(t.c.nodeC, "....the item should still exist with nodeC")
	t.c.nodeC.GetItem("expirytest", "test", false)
}

func NewExpirySuite(c *Config) *ExpirySuite {
	return &ExpirySuite{
		c: c,
	}
}
