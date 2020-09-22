package fred

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

const (
	fmtKgNodeString     = "kg-%s-node-%s"
	fmtKgStatusString   = "kg-%s-status"
	fmtKgMutableString  = "kg-%s-mutable"
	fmtKgExpiryString   = "kg-%s-expiry-node-%s"
	fmtNodeAdressString = "node-%s-address"
	nodePrefixString    = "node-"
	timeout             = 5 * time.Second
)

// nameService is the interface to the etcd server that serves as NaSe
// It is used by the replservice to sync updates to keygroups with other nodes and thereby makes sure that ReplicationStorage always has up to date information
type nameService struct {
	cli    *clientv3.Client
	NodeID string
}

// newNameService creates a new NameService
func newNameService(nodeID string, endpoints []string) (*nameService, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		// Deadline Exceeded
		return nil, errors.Errorf("Error starting the etcd client")
	}

	return &nameService{
		cli: cli, NodeID: nodeID,
	}, nil
}

// registerSelf stores information about this node
func (n *nameService) registerSelf(host Address) error {
	key := fmt.Sprintf(fmtNodeAdressString, n.NodeID)
	log.Debug().Msgf("NaSe: registering self as %s // %s", key, host.Addr)
	err := n.put(key, host.Addr)
	if err != nil {
		return err
	}
	return n.put(key, host.Addr)
}

/*
// registerOtherNode stores information about another node in NaSe
func (n *nameService) registerOtherNode(id NodeID, adress Address, port int) error {
	key := fmt.Sprintf(fmtNodeAdressString, id)
	value := fmt.Sprintf("%s:%d", string(adress.Addr), port)
	log.Debug().Msgf("NaSe: registering other node as %s // %s", key, value)
	return n.put(key, value)
}
*/

// existsKeygroup checks whether a Keygroup exists by checking whether there are keys with the prefix "kg-[kgname]-
func (n *nameService) existsKeygroup(key KeygroupName) (bool, error) {
	status, err := n.getKeygroupStatus(key)
	if err != nil {
		return false, err
	}
	return status == "created", nil
}

// isMutable checks whether a Keygroup is mutable.
func (n *nameService) isMutable(key KeygroupName) (bool, error) {
	status, err := n.getKeygroupMutable(key)
	if err != nil {
		return false, err
	}
	return status == "true", nil
}

// getExpiry checks the expiration time for items of the keygroup on a replica.
func (n *nameService) getExpiry(key KeygroupName) (int, error) {
	expiry, err := n.getKeygroupExpiry(key, n.NodeID)
	if err != nil {
		return 0, err
	}
	return expiry, nil
}

// createKeygroup created the keygroup status and joins the keygroup
func (n *nameService) createKeygroup(key KeygroupName, mutable bool, expiry int) error {
	exists, err := n.existsKeygroup(key)
	if err != nil {
		return err
	}
	if exists {
		return errors.Errorf("keygroup already exists in name service")
	}

	// Save the mutable attribute of the keygroup
	err = n.addKgMutableEntry(key, mutable)

	if err != nil {
		return err
	}

	// Save the expiry attribute of the keygroup for this replica
	err = n.addKgExpiryEntry(key, n.NodeID, expiry)

	if err != nil {
		return err
	}

	// Save the status of the keygroup
	err = n.addKgStatusEntry(key, "created")

	if err != nil {
		return err
	}

	// If the keygroup has existed before and was deleted it still has the old members
	// Why? Because all the nodes in this keygroup should be able to know that they are in a deleted keygroup
	// This is only a problem if a node doesn't see the delete state and only the new state in which it is not a member
	// But in this case it should just delete itself from the keygroup
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	_, err = n.cli.Delete(ctx, fmt.Sprintf(fmtKgNodeString, key, ""), clientv3.WithPrefix())

	if err != nil {
		return errors.New(err)
	}
	return n.addOwnKgNodeEntry(key, "ok")
}

// getKeygroupMembers returns all IDs of the Members of a Keygroup by iterating over all saved keys that start with the keygroup name
func (n *nameService) getKeygroupMembers(key KeygroupName, excludeSelf bool) (ids map[NodeID]int, err error) {
	nodes, err := n.getPrefix(fmt.Sprintf(fmtKgNodeString, string(key), ""))

	if err != nil {
		return nil, err
	}

	ids = make(map[NodeID]int)

	for i, value := range nodes {
		log.Debug().Msgf("NaSe: GetKeygroupMembers: Got result %d, key: %s value: %s", i, value.Key, value.Value)
		// If status is OK then add to available replicas
		if bytes.Equal(value.Value, []byte("ok")) {
			// If we are to exclude ourselfes
			if excludeSelf && n.NodeID == getNodeNameFromKgNodeString(string(value.Key)) {
				log.Debug().Msg("Excluding this node from results since this is the own node")
			} else {
				id := getNodeNameFromKgNodeString(string(value.Key))
				ids[NodeID(id)], err = n.getKeygroupExpiry(key, id)
				if err != nil {
					return
				}
			}

		} else {
			log.Debug().Msg("NaSe: GetKeygroupMembers: Above node has a status != OK, not returning it.")
		}
	}
	return
}

/*
// isKeygroupMember returns whether nodeID is in the keygroup kg
func (n *nameService) isKeygroupMember(nodeID string, kg KeygroupName) (bool, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgNodeString, kg, nodeID))
	if err != nil {
		return false, errors.Errorf("Node %s is not member of keygroup %s", nodeID, kg)
	}
	status := string(resp[0].Value)

	log.Debug().Msgf("NaSe: Is node %s member of keygroup %s ? Status is %s, so %#v", nodeID, kg, status, status == "ok")

	return status == "ok", nil

}*/

/*
// joinKeygroup joins an already existing keygroup
func (n *nameService) joinKeygroup(key KeygroupName) error {
	exists, err := n.existsKeygroup(key)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("keygroup does not exists so it cannot be joined")
	}
	return n.addOwnKgNodeEntry(key, "ok")
}
*/

// joinNodeIntoKeygroup joins the node into an already existing keygroup
func (n *nameService) joinNodeIntoKeygroup(key KeygroupName, nodeID NodeID, expiry int) error {
	exists, err := n.existsKeygroup(key)
	if err != nil {
		return err
	}

	if !exists {
		return errors.Errorf("keygroup does not exists so it cannot be joined by another node")
	}

	// Check whether the other node exists
	_, _, err = n.getNodeAddress(nodeID)
	if err != nil {
		log.Err(err).Msgf("Cannot join other node into a keygroup because the other nodes does not exist according to NaSe. Key: %s, otherNode: %s", key, nodeID)
		return err
	}

	// set expiry attribute for that particular node
	err = n.addKgExpiryEntry(key, string(nodeID), expiry)

	if err != nil {
		return err
	}

	return n.addOtherKgNodeEntry(nodeID, key, "ok")
}

/*
// exitKeygroup exits the local node from the keygroup
func (n *nameService) exitKeygroup(key KeygroupName) error {
	exists, err := n.existsKeygroup(key)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("keygroup does not exists so it cannot be exited")
	}
	return n.addOwnKgNodeEntry(key, "removed")
}
*/

// exitOtherNodeFromKeygroup deletes the node from the NaSe
func (n *nameService) exitOtherNodeFromKeygroup(key KeygroupName, nodeID NodeID) error {
	exists, err := n.existsKeygroup(key)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("keygroup does not exists so another node cannot exit it")
	}

	// Check whether the other node exists
	_, _, err = n.getNodeAddress(nodeID)
	if err != nil {
		log.Err(err).Msgf("Cannot exit other node from a keygroup because the other nodes does not exist according to NaSe. Key: %s, otherNode: %s", key, nodeID)
		return err
	}

	return n.addOtherKgNodeEntry(nodeID, key, "removed")
}

// deleteKeygroup marks the keygroup as "deleted" in the NaSe
func (n *nameService) deleteKeygroup(key KeygroupName) error {
	// Save the status of the keygroup
	err := n.addKgStatusEntry(key, "deleted")

	return err
}

// getNodeAddress returns the ip and port of a node
func (n *nameService) getNodeAddress(nodeID NodeID) (addr Address, port int, err error) {
	resp, err := n.getExact(fmt.Sprintf(fmtNodeAdressString, nodeID))

	if err != nil {
		return Address{}, 0, errors.New(err)
	}

	if len(resp) == 0 {
		return Address{}, 0, errors.Errorf("no such node %s", nodeID)
	}

	split := strings.Split(string(resp[0].Value), ":")
	addr, _ = ParseAddress(split[0])
	port, _ = strconv.Atoi(split[1])
	log.Debug().Msgf("Getting adress of node %s:: %s:%d", nodeID, addr.Addr, port)
	return addr, port, nil
}

// getAllNodes returns all nodes that are stored in the NaSe
func (n *nameService) getAllNodes() (nodes []Node, err error) {
	resp, err := n.getPrefix(nodePrefixString)
	for _, value := range resp {
		key := string(value.Key)
		res := strings.Split(string(value.Value), ":")
		// TODO status checks
		if strings.HasSuffix(key, "-status") {
			continue
		}
		// Now add node to return []
		nodeID := NodeID(strings.Split(key, "-")[1])
		addr, _ := ParseAddress(res[0])
		port, _ := strconv.Atoi(res[1])
		log.Debug().Msgf("NaSe: GetAllNodes: Got Response %s // %s", nodeID, res)
		node := &Node{
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
func (n *nameService) getPrefix(prefix string) (kv []*mvccpb.KeyValue, err error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	resp, err := n.cli.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))

	if err != nil {
		return nil, errors.New(err)
	}

	kv = resp.Kvs
	return
}

// getExact gets the exact key
func (n *nameService) getExact(key string) (kv []*mvccpb.KeyValue, err error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	resp, err := n.cli.Get(ctx, key)

	if err != nil {
		return nil, errors.New(err)
	}

	kv = resp.Kvs
	return
}

func (n *nameService) getKeygroupStatus(key KeygroupName) (string, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgStatusString, key))
	if resp == nil {
		return "", err
	}
	return string(resp[0].Value), err
}

func (n *nameService) getKeygroupMutable(key KeygroupName) (string, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgMutableString, key))
	if resp == nil {
		return "", err
	}
	return string(resp[0].Value), err
}

func (n *nameService) getKeygroupExpiry(key KeygroupName, id string) (int, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgExpiryString, key, id))
	if resp == nil {
		return 0, err
	}

	return strconv.Atoi(string(resp[0].Value))
}

// getCount returns the number of results getPrefix would return
// func (n *NameService) getCount(prefix string) (count int64, err error) {
// 	resp, err := n.cli.Get(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
// if err != nil {
// return 0, errors.New(err)
// }
// 	return resp.Count, nil
// }

// put puts the value into etcd.
func (n *nameService) put(key, value string) (err error) {
	ctx, _ := context.WithTimeout(context.TODO(), timeout)
	_, err = n.cli.Put(ctx, key, value)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// addOwnKgNodeEntry adds the entry for this node with a status.
func (n *nameService) addOwnKgNodeEntry(keygroup KeygroupName, status string) error {
	return n.put(n.fmtKgNode(keygroup), status)
}

// addOtherKgNodeEntry adds the entry for a remote node with a status.
func (n *nameService) addOtherKgNodeEntry(node NodeID, keygroup KeygroupName, status string) error {
	key := fmt.Sprintf(fmtKgNodeString, string(keygroup), node)
	return n.put(key, status)
}

// addKgStatusEntry adds the entry for a (new!) keygroup with a status.
func (n *nameService) addKgStatusEntry(keygroup KeygroupName, status string) error {
	return n.put(fmt.Sprintf(fmtKgStatusString, string(keygroup)), status)
}

// addKgMutableEntry adds the ismutable entry for a keygroup with a status.
func (n *nameService) addKgMutableEntry(keygroup KeygroupName, mutable bool) error {
	var data string

	if mutable {
		data = "true"
	} else {
		data = "false"
	}

	return n.put(fmt.Sprintf(fmtKgMutableString, string(keygroup)), data)
}

// addKgExpiryEntry adds the expiry entry for a keygroup with a status.
func (n *nameService) addKgExpiryEntry(keygroup KeygroupName, id string, expiry int) error {
	return n.put(fmt.Sprintf(fmtKgExpiryString, string(keygroup), id), strconv.Itoa(expiry))
}

// fmtKgNode turns a keygroup name into the key that this node will save its state in
// Currently: kg-[keygroup]-node-[NodeID]
func (n *nameService) fmtKgNode(keygroup KeygroupName) string {
	return fmt.Sprintf(fmtKgNodeString, string(keygroup), n.NodeID)
}

func getNodeNameFromKgNodeString(kgNode string) string {
	split := strings.Split(kgNode, "-")
	return split[len(split)-1]
}
