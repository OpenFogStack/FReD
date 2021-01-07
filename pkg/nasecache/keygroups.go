package nasecache

import (
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// get information about a keygroup

// IsMutable checks whether a Keygroup is mutable.
func (n *NameServiceCache) IsMutable(kg fred.KeygroupName) (bool, error) {
	key := "IsMutable" + sep + string(kg)

	// try to get from cache
	respCache, err := n.cache.Get(key)
	if err != nil {
		log.Debug().Msgf("NaSe: Key %s is not in cache", key)

		resp, err := n.regularNase.IsMutable(kg)

		if err != nil {
			return false, err
		}

		// put in cache
		err = n.cache.Set(key, boolToByteArray(resp))

		if err != nil {
			log.Err(err).Msg("Could not cache NaSe response")
		}

		return resp, nil
	}
	return byteArrayToBool(respCache), nil
}

// GetExpiry checks the expiration time for items of the keygroup on a replica.
func (n *NameServiceCache) GetExpiry(kg fred.KeygroupName) (int, error) {
	key := "GetExpiry" + sep + string(kg)

	// try to get from cache
	respCache, err := n.cache.Get(key)
	if err != nil {
		log.Debug().Msgf("NaSe: Key %s is not in cache", key)

		resp, err := n.regularNase.GetExpiry(kg)

		if err != nil {
			return 0, err
		}

		// put in cache
		err = n.cache.Set(key, intToByteArray(resp))

		if err != nil {
			log.Err(err).Msg("Could not cache NaSe response")
		}

		return resp, nil
	}
	return byteArrayToInt(respCache), nil
}

// manage keygroups

// ExistsKeygroup checks whether a Keygroup exists by checking whether there are keys with the prefix "kg|[kgname]|
func (n *NameServiceCache) ExistsKeygroup(kg fred.KeygroupName) (bool, error) {
	key := "ExistsKeygroup" + sep + string(kg)

	// try to get from cache
	respCache, err := n.cache.Get(key)
	if err != nil {
		log.Debug().Msgf("NaSe: Key %s is not in cache", key)

		resp, err := n.regularNase.ExistsKeygroup(kg)

		if err != nil {
			return false, err
		}

		// put in cache
		err = n.cache.Set(key, boolToByteArray(resp))

		if err != nil {
			log.Err(err).Msg("Could not cache NaSe response")
		}

		return resp, nil
	}
	return byteArrayToBool(respCache), nil
}

// JoinNodeIntoKeygroup joins the node into an already existing keygroup
func (n *NameServiceCache) JoinNodeIntoKeygroup(kg fred.KeygroupName, nodeID fred.NodeID, expiry int) error {
	return n.regularNase.JoinNodeIntoKeygroup(kg, nodeID, expiry)
}

// ExitOtherNodeFromKeygroup deletes the node from the NaSe
func (n *NameServiceCache) ExitOtherNodeFromKeygroup(kg fred.KeygroupName, nodeID fred.NodeID) error {
	return n.regularNase.ExitOtherNodeFromKeygroup(kg, nodeID)
}

// CreateKeygroup created the keygroup status and joins the keygroup
func (n *NameServiceCache) CreateKeygroup(kg fred.KeygroupName, mutable bool, expiry int) error {
	return n.regularNase.CreateKeygroup(kg, mutable, expiry)
}

// DeleteKeygroup marks the keygroup as "deleted" in the NaSe
func (n *NameServiceCache) DeleteKeygroup(kg fred.KeygroupName) error {
	return n.regularNase.DeleteKeygroup(kg)
}

// GetKeygroupMembers returns all IDs of the Members of a Keygroup by iterating over all saved keys that start with the keygroup name
func (n *NameServiceCache) GetKeygroupMembers(kg fred.KeygroupName, excludeSelf bool) (ids map[fred.NodeID]int, err error) {
	key := "GetKeygroupMembers" + sep + string(kg) + sep + boolToString(excludeSelf)

	// try to get from cache
	respCache, err := n.cache.Get(key)
	if err != nil {
		log.Debug().Msgf("NaSe: Key %s is not in cache", key)

		resp, err := n.regularNase.GetKeygroupMembers(kg, excludeSelf)

		if err != nil {
			return nil, err
		}

		// get byte array from response
		var b []byte
		err = genericToByteArray(resp, &b)

		if err != nil {
			log.Err(err).Msg("Could not cache NaSe response: byte array conversion failed")
		}

		// put in cache
		err = n.cache.Set(key, b)

		if err != nil {
			log.Err(err).Msg("Could not cache NaSe response")
		}

		return resp, nil
	}

	ids = make(map[fred.NodeID]int)
	err = byteArrayToGeneric(respCache, &ids)

	if err != nil {
		log.Err(err).Msg("Could not cache NaSe response: byte array conversion failed")
	}

	return ids, nil
}
