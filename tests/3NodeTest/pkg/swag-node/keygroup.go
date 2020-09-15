package node

import (
	"context"
	client "gitlab.tu-berlin.de/mcc-fred/fred/ext/go-client"

	"github.com/rs/zerolog/log"
)

// CreateKeygroup creates a new keygroup with the node. The Response should be empty if everything is correct.
func (n *Node) CreateKeygroup(kgname string, mutable bool, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Create Keygroup for group %s; expecting %d", kgname, expectedStatusCode)

	resp, recvErr := n.Client.KeygroupApi.KeygroupGroupIdPost(context.Background(), kgname, client.Body{Mutable: mutable})

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("CreateKeygroup: %s", err.Error())
		n.Errors++
	}
}

// DeleteKeygroup deletes the specified keygroup.
func (n *Node) DeleteKeygroup(kgname string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Delete Keygroup for group %s; expecting %d", kgname, expectedStatusCode)

	resp, recvErr := n.Client.KeygroupApi.KeygroupGroupIdDelete(context.Background(), kgname)

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("DeleteKeygroup: %s", err.Error())
		n.Errors++
	}
}

// GetKeygroupReplica gets all replica nodes for the specified keygroup.
func (n *Node) GetKeygroupReplica(kgname string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Get Keygroup Replica for group %s; expecting %d", kgname, expectedStatusCode)

	data, resp, recvErr := n.Client.KeygroupApi.KeygroupGroupIdReplicaGet(context.Background(), kgname)

	if err := checkResponse(resp, recvErr, data, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("GetKeygroupReplica: %s", err.Error())
		n.Errors++
	}
}

// AddKeygroupReplica adds a new Replica node to the provided keygroup.
func (n *Node) AddKeygroupReplica(kgname, replicaNodeID string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Add Keygroup Replica for group %s; expecting %d", kgname, expectedStatusCode)

	resp, recvErr := n.Client.KeygroupApi.KeygroupGroupIdReplicaNodeIdPost(context.Background(), kgname, replicaNodeID)

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("AddReplica: %s", err.Error())
		n.Errors++
	}
}

// DeleteKeygroupReplica deletes the specified keygroup.
func (n *Node) DeleteKeygroupReplica(kgname, replicaNodeID string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Delete Keygroup Replica for group %s; expecting %d", kgname, expectedStatusCode)

	resp, recvErr := n.Client.KeygroupApi.KeygroupGroupIdReplicaNodeIdDelete(context.Background(), kgname, replicaNodeID)

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("DeleteKeygroupReplica: %s", err.Error())
		n.Errors++
	}
}

// GetKeygroupTriggers gets all trigger nodes for the specified keygroup.
func (n *Node) GetKeygroupTriggers(kgname string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Get Keygroup Trigger for group %s; expecting %d", kgname, expectedStatusCode)

	data, resp, recvErr := n.Client.TriggersApi.KeygroupGroupIdTriggersGet(context.Background(), kgname)

	if err := checkResponse(resp, recvErr, data, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("GetKeygroupTriggers: %s", err.Error())
		n.Errors++
	}
}

// AddKeygroupTrigger adds a new Replica node to the provided keygroup.
func (n *Node) AddKeygroupTrigger(kgname, triggerNodeID, triggerNodeHost string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Add Keygroup Trigger for group %s; expecting %d", kgname, expectedStatusCode)

	node := client.TriggerNode{Host: triggerNodeHost}

	resp, recvErr := n.Client.TriggersApi.KeygroupGroupIdTriggersTriggerNodeIdPost(context.Background(), kgname, triggerNodeID, node)

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("AddKeygroupTrigger: %s", err.Error())
		n.Errors++
	}
}

// DeleteKeygroupTrigger deletes the specified keygroup.
func (n *Node) DeleteKeygroupTrigger(kgname, triggerNodeID string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Delete Keygroup Trigger for group %s; expecting %d", kgname, expectedStatusCode)

	resp, recvErr := n.Client.TriggersApi.KeygroupGroupIdTriggersTriggerNodeIdDelete(context.Background(), kgname, triggerNodeID)

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("DeleteKeygroupTrigger: %s", err.Error())
		n.Errors++
	}
}
