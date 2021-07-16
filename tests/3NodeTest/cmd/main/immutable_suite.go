package main

type ImmutableSuite struct {
	c *Config
}

func (t *ImmutableSuite) Name() string {
	return "Immutable"
}

func (t *ImmutableSuite) RunTests() {
	// testing immutable keygroups
	logNodeAction(t.c.nodeB, "Testing immutable keygroups by creating a new immutable keygroup on nodeB")
	t.c.nodeB.CreateKeygroup("log", false, 0, false)

	logNodeAction(t.c.nodeB, "Creating an item in this keygroup")
	res := t.c.nodeB.AppendItem("log", "value1", false)

	if res != "0" {
		logNodeFailure(t.c.nodeB, "0", res)
	}

	logNodeAction(t.c.nodeB, "Updating an item in this keygroup")
	t.c.nodeB.PutItem("log", res, "value-2", true)

	logNodeAction(t.c.nodeB, "Getting updated item from immutable keygroup")
	respB := t.c.nodeB.GetItem("log", res, false)

	if respB != "value1" {
		logNodeFailure(t.c.nodeB, "resp is value1", respB)
	}

	logNodeAction(t.c.nodeB, "Deleting an item in immutable keygroup")
	t.c.nodeB.DeleteItem("log", "0", true)

	logNodeAction(t.c.nodeB, "Adding nodeC as replica to immutable keygroup")
	t.c.nodeB.AddKeygroupReplica("log", t.c.nodeC.ID, 0, false)

	logNodeAction(t.c.nodeC, "Updating immutable item on other nodeC")
	t.c.nodeC.PutItem("log", "0", "value-3", true)

	logNodeAction(t.c.nodeC, "Appending another item to readonly log.")
	res = t.c.nodeC.AppendItem("log", "value-4", false)

	if res != "1" {
		logNodeFailure(t.c.nodeC, "1", res)
	}
}

func NewImmutableSuite(c *Config) *ImmutableSuite {
	return &ImmutableSuite{
		c: c,
	}
}
