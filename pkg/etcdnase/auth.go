package etcdnase

import (
	"fmt"
	"strings"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// RevokeUserPermissions removes user's permission to perform method on kg by deleting the key in etcd.
// It will not return an error if the user does not have the permission.
func (n *NameService) RevokeUserPermissions(user string, method fred.Method, kg fred.KeygroupName) error {
	log.Trace().Msgf("Revoking user %s's permission to %s on keygroup %s", user, method, kg)

	exists, err := n.ExistsKeygroup(kg)
	if err != nil {
		return err
	}

	if !exists {
		return errors.Errorf("keygroup does not exist")
	}

	prefix := fmt.Sprintf(fmtUserPermissionStringPrefix, user, string(kg))
	return n.delete(prefix+string(method), prefix)
}

// AddUserPermissions adds user's permission to perform method on kg by adding the key to etcd.
// It will not return an error if the user already has the permission.
func (n *NameService) AddUserPermissions(user string, method fred.Method, kg fred.KeygroupName) error {
	log.Trace().Msgf("Adding user %s's permission to %s on keygroup %s", user, method, kg)

	exists, err := n.ExistsKeygroup(kg)
	if err != nil {
		return err
	}

	if !exists {
		return errors.Errorf("keygroup does not exist")
	}

	prefix := fmt.Sprintf(fmtUserPermissionStringPrefix, user, string(kg))
	return n.put(prefix+string(method), "ok", prefix)
}

// GetUserPermissions returns a set of all of the user's permissions on kg from etcd.
func (n *NameService) GetUserPermissions(user string, kg fred.KeygroupName) (map[fred.Method]struct{}, error) {
	log.Trace().Msgf("Getting user %s's permissions on keygroup %s", user, kg)

	exists, err := n.ExistsKeygroup(kg)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.Errorf("keygroup does not exist")
	}

	res, err := n.getPrefix(fmt.Sprintf(fmtUserPermissionStringPrefix, user, string(kg)))

	if err != nil {
		return nil, err
	}

	permissions := make(map[fred.Method]struct{})

	for k := range res {
		permissions[fred.Method(strings.Split(k, sep)[5])] = struct{}{}
	}

	return permissions, nil
}
