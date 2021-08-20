package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type TriggerSuite struct {
	c *Config
}

func (t *TriggerSuite) Name() string {
	return "Trigger"
}
func (t *TriggerSuite) RunTests() {
	// let's test trigger nodes
	// create a new keygroup on nodeA
	logNodeAction(t.c.nodeA, "Creating keygroup triggertesting")
	t.c.nodeA.CreateKeygroup("triggertesting", true, 0, false)

	logNodeAction(t.c.nodeA, "Creating keygroup nottriggertesting")
	t.c.nodeA.CreateKeygroup("nottriggertesting", true, 0, false)

	//post an item1 to new keygroup
	logNodeAction(t.c.nodeA, "post an item1 to new keygroup triggertesting")
	t.c.nodeA.PutItem("triggertesting", "item1", "value1", false)
	//add trigger node to nodeA
	logNodeAction(t.c.nodeA, "add trigger node to nodeA for keygroup triggertesting")
	t.c.nodeA.AddKeygroupTrigger("triggertesting", t.c.triggerNodeID, t.c.triggerNodeHost, false)
	//post another item2 to new keygroup
	logNodeAction(t.c.nodeA, "post another item2 to new keygroup triggertesting")
	t.c.nodeA.PutItem("triggertesting", "item2", "value2", false)
	//delete item1 from keygroup
	logNodeAction(t.c.nodeA, "delete item1 from keygroup triggertesting")
	t.c.nodeA.DeleteItem("triggertesting", "item1", false)
	// post an item3 to keygroup nottriggertesting that should not be sent to trigger node
	logNodeAction(t.c.nodeA, "post an item3 to keygroup nottriggertesting that should not be sent to trigger node")
	t.c.nodeA.PutItem("nottriggertesting", "item3", "value3", false)
	//add keygroup to nodeB as well
	logNodeAction(t.c.nodeA, "add keygroup triggertesting to nodeB as well")
	t.c.nodeA.AddKeygroupReplica("triggertesting", t.c.nodeB.ID, 0, false)
	//post item4 to nodeB
	logNodeAction(t.c.nodeB, "post item4 to nodeB in keygroup triggertesting")
	t.c.nodeB.PutItem("triggertesting", "item4", "value4", false)
	//remove trigger node from nodeA
	logNodeAction(t.c.nodeA, "remove trigger node from nodeA in keygroup triggertesting")
	t.c.nodeA.DeleteKeygroupTrigger("triggertesting", t.c.triggerNodeID, false)
	//post item5 to nodeA
	logNodeAction(t.c.nodeA, "post item5 to nodeA in keygroup triggertesting")
	t.c.nodeA.PutItem("triggertesting", "item5", "value5", false)
	// check logs from trigger node
	// we should have the following logs (and nothing else):
	// put triggertesting item2 value2
	// del triggertesting item1
	// put triggertesting item4 value4
	logNodeAction(t.c.nodeA, "Checking if triggers have been executed correctly")
	checkTriggerNode(t.c.triggerNodeID, t.c.triggerNodeWSHost)
	logNodeAction(t.c.nodeA, "deleting keygroup triggertesting")
	t.c.nodeA.DeleteKeygroup("triggertesting", false)

	logNodeAction(t.c.nodeA, "deleting keygroup nottriggertesting")
	t.c.nodeA.DeleteKeygroup("nottriggertesting", false)

	logNodeAction(t.c.nodeA, "try to get the trigger nodes for keygroup triggertesting after deletion")
	t.c.nodeA.GetKeygroupTriggers("triggertesting", true)
}

func NewTriggerSuite(c *Config) *TriggerSuite {
	return &TriggerSuite{
		c: c,
	}
}

func checkTriggerNode(triggerNodeID, triggerNodeWSHost string) {
	log.Info().Str("trigger node", triggerNodeWSHost).Msg("Checking Trigger Node logs")

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

	if err != nil {
		log.Warn().Str("trigger node", triggerNodeWSHost).Msgf("%+v", err)
		return
	}

	var result []LogEntry
	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Warn().Str("trigger node", triggerNodeWSHost).Msgf("%+v", err)
		return
	}

	err = resp.Body.Close()

	if err != nil {
		log.Warn().Str("trigger node", triggerNodeWSHost).Msgf("%+v", err)
		return
	}

	if len(result) != len(expected) {
		log.Warn().Str("trigger node", triggerNodeID).Msgf("expected: %s, but got: %+v", expected, result)
		return
	}

	for i := range expected {
		if expected[i] != result[i] {
			log.Warn().Str("trigger node", triggerNodeID).Msgf("expected: %s, but got: %+v", expected[i], result[i])
			return
		}
	}
}
