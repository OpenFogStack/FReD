package etcdnase

import (
	"bytes"
	"context"
	"fmt"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"go.etcd.io/etcd/clientv3"
)

// ExitOtherNodeFromKeygroup deletes the node from the NaSe
func (n *NameService) ExitOtherNodeFromKeygroup(kg fred.KeygroupName, nodeID fred.NodeID) error {
	exists, err := n.ExistsKeygroup(kg)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("keygroup does not exists so another node cannot exit it")
	}

	// Check whether the other node exists
	_, err = n.GetNodeAddress(nodeID)
	if err != nil {
		log.Err(err).Msgf("Cannot exit other node from a keygroup because the other nodes does not exist according to NaSe. Key: %s, otherNode: %s", kg, nodeID)
		return err
	}

	return n.addOtherKgNodeEntry(string(nodeID), string(kg), "removed")
}

// DeleteKeygroup marks the keygroup as "deleted" in the NaSe
func (n *NameService) DeleteKeygroup(kg fred.KeygroupName) error {
	// Save the status of the keygroup
	err := n.addKgStatusEntry(string(kg), "deleted")

	return err
}

// ExistsKeygroup checks whether a Keygroup exists by checking whether there are keys with the prefix "kg|[kgname]|
func (n *NameService) ExistsKeygroup(kg fred.KeygroupName) (bool, error) {
	status, err := n.getKeygroupStatus(string(kg))
	if err != nil {
		return false, err
	}
	return status == "created", nil
}

// IsMutable checks whether a Keygroup is mutable.
func (n *NameService) IsMutable(kg fred.KeygroupName) (bool, error) {
	status, err := n.getKeygroupMutable(string(kg))
	if err != nil {
		return false, err
	}
	return status == "true", nil
}

// GetExpiry checks the expiration time for items of the keygroup on a replica.
func (n *NameService) GetExpiry(kg fred.KeygroupName) (int, error) {
	expiry, err := n.getKeygroupExpiry(string(kg), n.NodeID)
	if err != nil {
		return 0, err
	}
	return expiry, nil
}

// CreateKeygroup created the keygroup status and joins the keygroup
func (n *NameService) CreateKeygroup(kg fred.KeygroupName, mutable bool, expiry int) error {
	exists, err := n.ExistsKeygroup(kg)
	if err != nil {
		return err
	}
	if exists {
		return errors.Errorf("keygroup already exists in name service")
	}

	// Save the mutable attribute of the keygroup
	err = n.addKgMutableEntry(string(kg), mutable)

	if err != nil {
		return err
	}

	// Save the expiry attribute of the keygroup for this replica
	err = n.addKgExpiryEntry(string(kg), n.NodeID, expiry)

	if err != nil {
		return err
	}

	// Save the status of the keygroup
	err = n.addKgStatusEntry(string(kg), "created")

	if err != nil {
		return err
	}

	// If the keygroup has existed before and was deleted it still has the old members
	// Why? Because all the nodes in this keygroup should be able to know that they are in a deleted keygroup
	// This is only a problem if a node doesn't see the delete state and only the new state in which it is not a member
	// But in this case it should just delete itself from the keygroup
	ctx, cncl := context.WithTimeout(context.Background(), timeout)

	defer cncl()

	_, err = n.cli.Delete(ctx, fmt.Sprintf(fmtKgNodeString, string(kg), ""), clientv3.WithPrefix())

	if err != nil {
		return errors.New(err)
	}
	return n.addOwnKgNodeEntry(string(kg), "ok")
}

// GetKeygroupMembers returns all IDs of the Members of a Keygroup by iterating over all saved keys that start with the keygroup name.
// The value of the map is the expiry in seconds.
func (n *NameService) GetKeygroupMembers(kg fred.KeygroupName, excludeSelf bool) (ids map[fred.NodeID]int, err error) {
	nodes, err := n.getPrefix(fmt.Sprintf(fmtKgNodeString, string(kg), ""))

	if err != nil {
		return nil, err
	}

	ids = make(map[fred.NodeID]int)

	for i, value := range nodes {
		// If status is OK then add to available replicas
		if bytes.Equal(value.Value, []byte("ok")) {
			// If we are to exclude ourselves
			if excludeSelf && n.NodeID == getNodeNameFromKgNodeString(string(value.Key)) {
				log.Debug().Msgf("NaSe: GetKeygroupMembers: Got result %d, key: %s value: %s", i, value.Key, value.Value)
				log.Debug().Msg("...Excluding this node from results since this is the own node")
			} else {
				id := getNodeNameFromKgNodeString(string(value.Key))
				ids[fred.NodeID(id)], err = n.getKeygroupExpiry(string(kg), id)
				if err != nil {
					return
				}
			}

		} else {
			log.Debug().Msgf("NaSe: GetKeygroupMembers: Got result %d, key: %s value: %s", i, value.Key, value.Value)
			log.Debug().Msg("... node has a status != OK, not returning it.")
		}
	}
	return
}

// JoinNodeIntoKeygroup joins the node into an already existing keygroup
func (n *NameService) JoinNodeIntoKeygroup(kg fred.KeygroupName, nodeID fred.NodeID, expiry int) error {
	exists, err := n.ExistsKeygroup(kg)
	if err != nil {
		return err
	}

	if !exists {
		return errors.Errorf("keygroup does not exists so it cannot be joined by another node")
	}

	// Check whether the other node exists
	_, err = n.GetNodeAddress(nodeID)
	if err != nil {
		log.Err(err).Msgf("Cannot join other node into a keygroup because the other nodes does not exist according to NaSe. Key: %s, otherNode: %s", kg, nodeID)
		return err
	}

	// set expiry attribute for that particular node
	err = n.addKgExpiryEntry(string(kg), string(nodeID), expiry)

	if err != nil {
		return err
	}

	return n.addOtherKgNodeEntry(string(nodeID), string(kg), "ok")
}
