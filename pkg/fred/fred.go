package fred

// Config holds configuration parameters for an instance of FReD.
type Config struct {
	Store   Store
	Client  Client
	ZmqPort int
	NodeID  string
}

// Fred is an instance of FReD.
type Fred struct {
	E ExtHandler
	I IntHandler
}

// IntHandler is an interface that abstracts the methods of the handler that handles internal requests.
type IntHandler interface {
	HandleCreateKeygroup(k Keygroup, nodes []Node) error
	HandleDeleteKeygroup(k Keygroup) error
	HandleUpdate(i Item) error
	HandleDelete(i Item) error
	HandleAddReplica(k Keygroup, n Node) error
	HandleRemoveReplica(k Keygroup, n Node) error
	HandleAddNode(n Node) error
	HandleRemoveNode(n Node) error
	HandleIntroduction(introducer Node, self Node, node []Node) error
	HandleDetroduction() error
}

// ExtHandler is an interface that abstracts the methods of the handler that handles external requests.
type ExtHandler interface {
	HandleCreateKeygroup(k Keygroup) error
	HandleDeleteKeygroup(k Keygroup) error
	HandleRead(i Item) (Item, error)
	HandleUpdate(i Item) error
	HandleDelete(i Item) error
	HandleAddReplica(k Keygroup, n Node) error
	HandleGetKeygroupReplica(k Keygroup) ([]Node, error)
	HandleRemoveReplica(k Keygroup, n Node) error
	HandleAddNode(n []Node) error
	HandleGetReplica(n Node) (Node, error)
	HandleGetAllReplica() ([]Node, error)
	HandleRemoveNode(n Node) error
}

// New creates a new FReD instance.
func New(config *Config) (f Fred) {
	s := newStoreService(config.Store)
	k := newKeygroupStore(config.NodeID)
	r := newReplicationService(config.ZmqPort, config.Client)

	return Fred{
		E: newExthandler(s, k, r),
		I: newInthandler(s, k, r),
	}
}
