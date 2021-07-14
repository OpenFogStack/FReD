package etcdnase

import (
	"fmt"
	"math"
	"strings"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"github.com/rs/zerolog/log"
)

// ReportFailedNode saves that a node has missed an update to a keygroup that it will get another time
func (n *NameService) ReportFailedNode(nodeID fred.NodeID, kg fred.KeygroupName, id string) error {
	log.Warn().Msgf("Nase: ReportFailedNode: Reporting that nodeId %#v has missed kg %#v id %s", nodeID, kg, id)
	item := fmt.Sprintf(fmtFailedNodeKgString, nodeID, kg, id)
	log.Debug().Msgf("NaSe.ReportFailedNode: Putting %s into NaSe", item)
	err := n.put(item, "1")

	if err != nil {
		log.Err(err).Msgf("NaSe: ReportFailedNode: Node was not able to reach NaSe." +
			"So maybe this node has no Internet connection?")
		return err
	}
	return nil
}

// RequestNodeStatus request a list of items that the node has missed while it was offline
func (n *NameService) RequestNodeStatus(nodeID fred.NodeID) (kgs []fred.Item) {
	resp, err := n.getPrefix(fmt.Sprintf(fmtFailedNodePrefix, nodeID))

	if err != nil {
		log.Err(err).Msgf("NaSe: RequestNodeStatus: Failed to get Prefix")
		return nil
	}

	log.Debug().Msgf("Nase: RequestNodeStatus: found %d items that were missed", len(resp))

	for k := range resp {
		split := strings.Split(k, sep)
		kgname := split[3]
		id := split[4]
		log.Debug().Msgf("NaSe: RequestNodeStatus: missed item has kgname %s and id %s (from key in NaSe: %s)", kgname, id, k)
		kgs = append(kgs, fred.Item{
			Keygroup: fred.KeygroupName(kgname),
			ID:       id,
		})
		err = n.delete(k)

		if err != nil {
			log.Err(err).Msgf("Could not remove missed data entry %s for node %v", id, nodeID)
		}
	}
	return
}

// GetNodeWithBiggerExpiry if this node has to get an item because it has missed it, it has to get it from a node with a bigger expiry
// if there is no node with a bigger expiry then it returns the node with the highest expiry
func (n *NameService) GetNodeWithBiggerExpiry(kg fred.KeygroupName) (nodeID fred.NodeID, addr string) {
	log.Debug().Msgf("Nase: GetNodeWithBiggerExpiry finding node that replicates %s with expiry bigger than own node", string(kg))
	expiry, err := n.GetExpiry(kg)
	if err != nil {
		return "", ""
	}

	if expiry == 0 {
		// For easier comparisons
		expiry = math.MaxInt32
	}

	nodes, err := n.GetKeygroupMembers(kg, true)
	if err != nil || len(nodes) == 0 {
		// Error or no nodes found
		return "", ""
	}

	currentBestExpiry := 0
	var currentBestNodeID fred.NodeID

	for node, exp := range nodes {
		if exp == 0 {
			exp = math.MaxInt32
		}
		if exp > currentBestExpiry {
			currentBestExpiry = exp
			currentBestNodeID = node
			log.Debug().Msgf("Found node %s with Expiry %d", string(node), exp)
			if exp == math.MaxInt32 {
				break
			}
		}
	}

	log.Debug().Msgf("Returning node %s with Expiry %d", string(currentBestNodeID), currentBestExpiry)
	if currentBestExpiry < expiry {
		log.Warn().Msgf("NaSe: GetNodeWithBiggerExpiry: Was not able to find node with bigger expiry than %d, using node %s with expiry %d instead", expiry, currentBestNodeID, currentBestExpiry)
	}
	addr, _ = n.GetNodeAddress(currentBestNodeID)
	return currentBestNodeID, addr

}
