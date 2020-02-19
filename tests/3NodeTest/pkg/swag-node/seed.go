package node

import (
	"context"

	"github.com/rs/zerolog/log"

	client "gitlab.tu-berlin.de/mcc-fred/fred/ext/go-client"
)

// SeedNode makes the current node a seed node.
func (n *Node) SeedNode(nodeID, nodeHost string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Post to Seed; expecting %d", expectedStatusCode)

	resp, recvErr := n.Client.SeedApi.SeedPost(context.Background(), client.Seed{
		Id:   nodeID,
		Addr: nodeHost,
	})

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("SeedNode: %s", err.Error())
		n.Errors++
	}
}
