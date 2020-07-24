package nameservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	frederrors "gitlab.tu-berlin.de/mcc-fred/fred/pkg/errors"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

const (
	fmtKgNodeString     = "kg-%s-node-%s"
	fmtKgStatusString   = "kg-%s-status"
	fmtKgString         = "kg-%s-"
	fmtNodeAdressString = "node-%s-address"
	fmtNodeStatusString = "node-%s-status"
	nodePrefixString    = "node-"
)

// NameService is the interface to the etcd server that serves as NaSe
// It is used by the replservice to sync updates to keygroups with other nodes and thereby makes sure that ReplicationStorage always has up to date information
type NameService struct {
	cli    *clientv3.Client
	NodeID string
}

// New creates a new NameService
func New(nodeID string, endpoints []string) *NameService {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		// Deadline Exceeded
		log.Err(err).Msg("Error starting the etcd client")
	}

	return &NameService{cli: cli, NodeID: nodeID}
}

// RegisterSelf stores information about this node
func (n *NameService) RegisterSelf(address replication.Address, port int) error {
	key := fmt.Sprintf(fmtNodeAdressString, n.NodeID)
	value := fmt.Sprintf("%s:%d", string(address.Addr), port)
	log.Debug().Msgf("NaSe: registering self as %s // %s", key, value)
	return n.put(key, value)
}

// RegisterOtherNode stores information about another node in NaSe
func (n *NameService) RegisterOtherNode(id replication.ID, adress replication.Address, port int) error {
	key := fmt.Sprintf(fmtNodeAdressString, id)
	value := fmt.Sprintf("%s:%d", string(adress.Addr), port)
	log.Debug().Msgf("NaSe: registering other node as %s // %s", key, value)
	return n.put(key, value)
}

// ExistsKeygroup checks whether a Keygroup exists by checking whether there are keys with the prefix "kg-[kgname]-
func (n *NameService) ExistsKeygroup(key commons.KeygroupName) (bool, error) {
	status, err := n.getKeygroupStatus(key)
	if err != nil {
		return false, err
	}
	return status == "created", nil
}

// CreateKeygroup created the keygroup status and joins the keygroup
func (n *NameService) CreateKeygroup(key commons.KeygroupName) error {
	exists, err := n.ExistsKeygroup(key)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("keygroup already exists in name service")
	}

	// Save the status of the keygroup
	err = n.addKgStatusEntry(key, "created")

	// If the keygroup has existed before and was deleted it still has the old members
	// Why? Because all the nodes in this keygroup should be able to know that they are in a deleted keygroup
	// This is only a problem if a node doesn't see the delete state and only the new state in which it is not a member
	// But in this case it should just delete itself from the keygroup
	n.cli.Delete(context.Background(), fmt.Sprintf(fmtKgNodeString, key, ""), clientv3.WithPrefix())

	if err != nil {
		return err
	}
	return n.addOwnKgNodeEntry(key, "ok")
}

// GetKeygroupMembers returns all IDs of the Members of a Keygroup by iterating over all saved keys that start with the keygroup name
func (n *NameService) GetKeygroupMembers(key commons.KeygroupName, excludeSelf bool) (ids []replication.ID, err error) {
	resp, err := n.getPrefix(fmt.Sprintf(fmtKgNodeString, string(key), ""))
	if err != nil {
		return nil, err
	}
	for i, value := range resp {
		log.Debug().Msgf("NaSe: GetKeygroupMembers: Got result %d, key: %s value: %s", i, value.Key, value.Value)
		// If status is OK then add to available replicas
		if bytes.Equal(value.Value, []byte("ok")) {
			// If we are to exclude ourselfes
			if excludeSelf && n.NodeID == getNodeNameFromKgNodeString(string(value.Key)) {
				log.Debug().Msg("Excluding this node from results since this is the own node")
			} else {
				ids = append(ids, replication.ID(getNodeNameFromKgNodeString(string(value.Key))))
			}

		} else {
			log.Debug().Msg("NaSe: GetKeygroupMembers: Above node has a status != OK, not returning it.")
		}
	}
	return
}

// IsKeygroupMember returns whether nodeID is in the keygroup kg
func (n *NameService) IsKeygroupMember(nodeID string, kg commons.KeygroupName) (bool, error){
	resp, err := n.getExact(fmt.Sprintf(fmtKgNodeString, kg, nodeID))
	if err != nil {
		return false, frederrors.New(400, fmt.Sprintf("Node %s is not member of keygroup %s", nodeID, kg))
	}
	status := string(resp[0].Value)

	log.Debug().Msgf("NaSe: Is node %s member of keygroup %s ? Status is %s, so %#v", nodeID, kg, status, status=="ok")

	return status == "ok", nil

}

// JoinKeygroup joins an already existing keygroup
func (n *NameService) JoinKeygroup(key commons.KeygroupName) error {
	exists, err := n.ExistsKeygroup(key)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("keygroup does not exists so it cannot be joined")
	}
	return n.addOwnKgNodeEntry(key, "ok")
}

// JoinOtherNodeIntoKeygroup joins the node into an already existing keygroup
func (n *NameService) JoinOtherNodeIntoKeygroup(key commons.KeygroupName, otherNodeID replication.ID) error {
	exists, err := n.ExistsKeygroup(key)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("keygroup does not exists so it cannot be joined by another node")
	}

	// Check whether the other node exists
	_, _, err = n.GetNodeAdress(otherNodeID)
	if err != nil {
		log.Err(err).Msgf("Cannot join other node into a keygroup because the other nodes does not exist according to NaSe. Key: %s, otherNode: %s", key, otherNodeID)
		return err
	}

	return n.addOtherKgNodeEntry(otherNodeID, key, "ok")
}

// ExitKeygroup exits the local node from the keygroup
func (n *NameService) ExitKeygroup(key commons.KeygroupName) error {
	exists, err := n.ExistsKeygroup(key)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("keygroup does not exists so it cannot be exited")
	}
	return n.addOwnKgNodeEntry(key, "removed")
}

// ExitOtherNodeFromKeygroup deletes the node from the NaSe
func (n *NameService) ExitOtherNodeFromKeygroup(key commons.KeygroupName, otherNodeID replication.ID) error {
	exists, err := n.ExistsKeygroup(key)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("keygroup does not exists so another node cannot exit it")
	}

	// Check whether the other node exists
	_, _, err = n.GetNodeAdress(otherNodeID)
	if err != nil {
		log.Err(err).Msgf("Cannot exit other node from a keygroup because the other nodes does not exist according to NaSe. Key: %s, otherNode: %s", key, otherNodeID)
		return err
	}

	return n.addOtherKgNodeEntry(otherNodeID, key, "removed")
}

// DeleteKeygroup marks the keygroup as "deleted" in the NaSe
func (n *NameService) DeleteKeygroup(key commons.KeygroupName) error {
	return n.put(fmt.Sprintf(fmtKgStatusString, key), "deleted")
}

// GetNodeAdress returns the ip and port of a node
func (n *NameService) GetNodeAdress(nodeID replication.ID) (addr replication.Address, port int, err error) {
	resp, err := n.getExact(fmt.Sprintf(fmtNodeAdressString, nodeID))
	split := strings.Split(string(resp[0].Value), ":")
	addr, _ = replication.ParseAddress(split[0])
	port, _ = strconv.Atoi(split[1])
	log.Debug().Msgf("Getting adress of node %s:: %s:%d", nodeID, addr.Addr, port)
	return
}

// GetAllNodes returns all nodes that are stored in the NaSe
func (n *NameService) GetAllNodes() (nodes []replication.Node, err error) {
	resp, err := n.getPrefix(nodePrefixString)
	for _, value := range resp {
		key := string(value.Key)
		res := strings.Split(string(value.Value), ":")
		// TODO status checks
		if strings.HasSuffix(key,"-status") {
			continue
		}
		// Now add node to return []
		nodeID := replication.ID(strings.Split(key, "-")[1])
		addr, _ := replication.ParseAddress(res[0])
		port, _ := strconv.Atoi(res[1])
		log.Debug().Msgf("NaSe: GetAllNodes: Got Response %s // %s", nodeID, res)
		node := &replication.Node{
			ID:   nodeID,
			Addr: addr,
			Port: port,
		}
		nodes = append(nodes, *node)
	}
	return
}

// getPrefix gets every key that starts(!) with the specified string
// the keys are sorted ascending by key for easier debugging
func (n *NameService) getPrefix(prefix string) (kv []*mvccpb.KeyValue, err error) {
	resp, err := n.cli.Get(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	kv = resp.Kvs
	return
}

// getExact gets the exact key
func (n *NameService) getExact(key string) (kv []*mvccpb.KeyValue, err error) {
	resp, err := n.cli.Get(context.Background(), key)
	kv = resp.Kvs
	return
}

func (n *NameService) getKeygroupStatus(key commons.KeygroupName) (string, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgStatusString, key))
	if resp == nil {
		return "", err
	}
	return string(resp[0].Value), err
}

// getCount returns the number of results getPrefix would return
// func (n *NameService) getCount(prefix string) (count int64, err error) {
// 	resp, err := n.cli.Get(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
// 	count = resp.Count
// 	return
// }

// put puts the value into etcd.
func (n *NameService) put(key, value string) (err error) {
	_, err = n.cli.Put(context.Background(), key, value)
	return
}

// add the entry for this node with a status
func (n *NameService) addOwnKgNodeEntry(keygroup commons.KeygroupName, status string) error {
	return n.put(n.fmtKgNode(keygroup), status)
}

// add the entry for a remote node with a status
func (n *NameService) addOtherKgNodeEntry(node replication.ID, keygroup commons.KeygroupName, status string) error {
	key := fmt.Sprintf(fmtKgNodeString, string(keygroup), node)
	return n.put(key, status)
}

// add the entry for a (new!) keygroup with a status
func (n *NameService) addKgStatusEntry(keygroup commons.KeygroupName, status string) error {
	return n.put(fmt.Sprintf(fmtKgStatusString, string(keygroup)), status)
}

// fmtKgNodeString turns a keygroup name into the key that this node will save its state in
// Currently: kg-[keygroup]-node-[NodeID]
func (n *NameService) fmtKgNode(keygroup commons.KeygroupName) string {
	return fmt.Sprintf(fmtKgNodeString, string(keygroup), n.NodeID)
}

func getNodeNameFromKgNodeString(kgNode string) string {
	split := strings.Split(kgNode, "-")
	return split[len(split)-1]
}
