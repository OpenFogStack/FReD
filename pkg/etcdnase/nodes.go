package etcdnase

import (
	"fmt"
	"strings"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// GetNodeAddress returns the ip and port of a node
func (n *NameService) GetNodeAddress(nodeID fred.NodeID) (addr string, err error) {
	resp, err := n.getExact(fmt.Sprintf(fmtNodeAdressString, string(nodeID)))

	if err != nil {
		return "", errors.New(err)
	}

	if len(resp) == 0 {
		return "", errors.Errorf("no such node %s", nodeID)
	}

	log.Debug().Msgf("NaSe: GetNodeAdress: Address of node %s is %s", nodeID, resp[0].Value)
	return string(resp[0].Value), nil
}

// GetAllNodes returns all nodes that are stored in the NaSe in the way they can be reached by other nodes
func (n *NameService) GetAllNodes() (nodes []fred.Node, err error) {
	return n.getAllNodesBySuffix(fmt.Sprintf(sep + "address"))
}

// GetAllNodesExternal returns all nodes with the port that the exthandler is running on
func (n *NameService) GetAllNodesExternal() (nodes []fred.Node, err error) {
	return n.getAllNodesBySuffix(fmt.Sprintf(sep + "extaddress"))
}

func (n *NameService) getAllNodesBySuffix(suffix string) (nodes []fred.Node, err error) {
	resp, err := n.getPrefix(nodePrefixString)

	nodes = make([]fred.Node, 0)

	for _, value := range resp {
		key := string(value.Key)

		// TODO status checks
		if strings.HasSuffix(key, "|status") {
			continue
		}

		if !strings.HasSuffix(key, suffix) {
			continue
		}

		// Now add node to return []
		nodeID := strings.Split(key, sep)[1]

		log.Debug().Msgf("NaSe: GetAllNodes: Got Response %s // %s", nodeID, string(value.Value))

		nodes = append(nodes, fred.Node{
			ID:   fred.NodeID(nodeID),
			Host: string(value.Value),
		})
	}
	return
}
