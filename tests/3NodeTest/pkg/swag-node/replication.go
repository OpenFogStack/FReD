package node

import (
	"context"

	"github.com/rs/zerolog/log"

	client "gitlab.tu-berlin.de/mcc-fred/fred/ext/go-client"
)

// RegisterReplica registers a new replica with this node.
func (n *Node) RegisterReplica(nodeID, nodeIP string, nodePort int, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Registering Replica %s; expecting %d", nodeID, expectedStatusCode)

	nodes := make([]client.Node, 1)

	nodes[0] = client.Node{
		Id:      nodeID,
		Addr:    nodeIP,
		ZmqPort: int32(nodePort),
	}

	resp, recvErr := n.Client.ReplicationApi.ReplicaPost(context.Background(), client.Replica{Nodes: nodes})

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("RegisterReplica: %s", err.Error())
		n.Errors++
	}

	return
}

// GetAllReplica returns a list of all replica that this node has stored.
func (n *Node) GetAllReplica(expectedStatusCode int, expectEmptyResponse bool) []string {
	log.Debug().Str("node", n.URL).Msgf("Sending a Get for all Replicas; expecting %d", expectedStatusCode)

	data, resp, recvErr := n.Client.ReplicationApi.ReplicaGet(context.Background())

	if err := checkResponse(resp, recvErr, data, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("GetAllReplica: %s", err.Error())
		n.Errors++
		return nil
	}

	return data.Nodes
}

// GetReplica returns a replica.
func (n *Node) GetReplica(nodeID string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Get for Replica %s; expecting %d", nodeID, expectedStatusCode)

	data, resp, recvErr := n.Client.ReplicationApi.ReplicaNodeIdGet(context.Background(), nodeID)

	if err := checkResponse(resp, recvErr, data, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("GetReplica: %s", err.Error())
		n.Errors++
	}

	return
}

// DeleteReplica deletes a replica.
func (n *Node) DeleteReplica(nodeID string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Delete for Replica %s; expecting %d", nodeID, expectedStatusCode)

	resp, recvErr := n.Client.ReplicationApi.ReplicaNodeIdDelete(context.Background(), nodeID)

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("DeleteReplica: %s", err.Error())
		n.Errors++
	}

	return
}
