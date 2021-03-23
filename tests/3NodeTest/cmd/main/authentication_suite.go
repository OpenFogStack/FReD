package main

type AuthenticationSuite struct {
	c *Config
}

func (t *AuthenticationSuite) RunTests() {
	// test RBAC and authentication
	logNodeAction(t.c.nodeA, "create keygroup \"rbactest\"")
	t.c.nodeA.CreateKeygroup("rbactest", true, 0, false)

	logNodeAction(t.c.nodeA, "put item into keygroup \"rbactest\"")
	t.c.nodeA.PutItem("rbactest", "item1", "value1", false)

	logNodeAction(t.c.nodeA, "add little client as read only to rbac test")
	t.c.nodeA.AddUser("littleclient", "rbactest", "ReadKeygroup", false)

	logNodeAction(t.c.nodeA, "try to read with little client -> should work")

	if val := t.c.littleClient.GetItem("rbactest", "item1", false); val != "value1" {
		logNodeFailure(t.c.nodeA, "value1", val)
	}

	logNodeAction(t.c.nodeA, "try to write with little client -> should not work")
	t.c.littleClient.PutItem("rbactest", "item1", "value2", true)

	logNodeAction(t.c.nodeA, "add role configure replica to little client -> should work")
	t.c.nodeA.AddUser("littleclient", "rbactest", "ConfigureReplica", false)

	logNodeAction(t.c.nodeA, "add replica nodeB to keygroup with little client -> should work")
	t.c.littleClient.AddKeygroupReplica("rbactest", t.c.nodeBpeeringID, 0, false)

	logNodeAction(t.c.nodeB, "remove permission to read from keygroup -> should work")
	t.c.nodeB.RemoveUser("littleclient", "rbactest", "ReadKeygroup", false)

	logNodeAction(t.c.nodeA, "try to read from keygroup with little client -> should not work")
	if val := t.c.littleClient.GetItem("rbactest", "item1", true); val != "" {
		logNodeFailure(t.c.nodeA, "", val)
	}
}

func NewAuthenticationSuite(c *Config) *AuthenticationSuite {
	return &AuthenticationSuite{
		c: c,
	}
}
