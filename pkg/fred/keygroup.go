package fred

// Keygroup has a name and a list of replica nodes.
type Keygroup struct {
	Name    KeygroupName
	Replica map[NodeID]struct{}
}

// KeygroupName is a name of a keygroup.
type KeygroupName string
