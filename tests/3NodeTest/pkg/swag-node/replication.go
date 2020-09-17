package node

import (
	"context"

	"github.com/rs/zerolog/log"
)

// GetAllReplica returns a list of all replica that this node has stored.
func (n *Node) GetAllReplica(expectedStatusCode int, expectEmptyResponse bool) map[string]int {
	log.Debug().Str("node", n.URL).Msgf("Sending a Get for all Replicas; expecting %d", expectedStatusCode)

	data, resp, recvErr := n.Client.ReplicationApi.ReplicaGet(context.Background())

	if err := checkResponse(resp, recvErr, data, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("GetAllReplica: %s", err.Error())
		n.Errors++
		return nil
	}

	nodes := make(map[string]int)

	for _, n := range data.Nodes {
		nodes[n.Id] = int(n.Expiry)
	}

	return nodes
}

// GetReplica returns a replica.
func (n *Node) GetReplica(nodeID string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Get for Replica %s; expecting %d", nodeID, expectedStatusCode)

	data, resp, recvErr := n.Client.ReplicationApi.ReplicaNodeIdGet(context.Background(), nodeID)

	if err := checkResponse(resp, recvErr, data, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("GetReplica: %s", err.Error())
		n.Errors++
	}
}

// DeleteReplica deletes a replica.
func (n *Node) DeleteReplica(nodeID string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Delete for Replica %s; expecting %d", nodeID, expectedStatusCode)

	resp, recvErr := n.Client.ReplicationApi.ReplicaNodeIdDelete(context.Background(), nodeID)

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("DeleteReplica: %s", err.Error())
		n.Errors++
	}
}
