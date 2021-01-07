package nasecache

import (
	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// manage permissions

// RevokeUserPermissions removes user's permission to perform method on kg by deleting the key in etcd.
func (n *NameServiceCache) RevokeUserPermissions(user string, method fred.Method, kg fred.KeygroupName) error {
	return n.regularNase.RevokeUserPermissions(user, method, kg)
}

// AddUserPermissions adds user's permission to perform method on kg by adding the key to etcd.
func (n *NameServiceCache) AddUserPermissions(user string, method fred.Method, kg fred.KeygroupName) error {
	return n.regularNase.AddUserPermissions(user, method, kg)
}

// GetUserPermissions returns a set of all of the user's permissions on kg from etcd.
func (n *NameServiceCache) GetUserPermissions(user string, kg fred.KeygroupName) (map[fred.Method]struct{}, error) {
	return n.regularNase.GetUserPermissions(user, kg)
}
