package fred

// NameService interface abstracts from the features of the nameservice, whether that is etcd or a distributed implementation.
type NameService interface {
	// manage information about this node
	GetNodeID() NodeID
	RegisterSelf(host string, externalHost string) error

	// manage permissions
	AddUserPermissions(user string, method Method, keygroup KeygroupName) error
	RevokeUserPermissions(user string, method Method, keygroup KeygroupName) error
	GetUserPermissions(user string, keygroup KeygroupName) (map[Method]struct{}, error)

	// get information about a keygroup
	IsMutable(kg KeygroupName) (bool, error)
	GetExpiry(kg KeygroupName) (int, error)

	// manage information about another node
	GetNodeAddress(nodeID NodeID) (addr string, err error)
	GetAllNodes() (nodes []Node, err error)
	GetAllNodesExternal() (nodes []Node, err error)

	// manage keygroups
	ExistsKeygroup(kg KeygroupName) (bool, error)
	JoinNodeIntoKeygroup(key KeygroupName, nodeID NodeID, expiry int) error
	ExitOtherNodeFromKeygroup(kg KeygroupName, nodeID NodeID) error
	CreateKeygroup(kg KeygroupName, mutable bool, expiry int) error
	DeleteKeygroup(kg KeygroupName) error
	GetKeygroupMembers(kg KeygroupName, excludeSelf bool) (ids map[NodeID]int, err error)

	// handle node failures
	ReportFailedNode(nodeID NodeID, kg KeygroupName, id string) error
	RequestNodeStatus(nodeID NodeID) []Item
	GetNodeWithBiggerExpiry(kg KeygroupName) (nodeID NodeID, addr string)
}
