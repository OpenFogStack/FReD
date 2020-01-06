package replication

// KeygroupName is a name of a keygroup.
type KeygroupName string

// Keygroup has a name and a list of replica nodes.
type Keygroup struct {
	Name KeygroupName
	Replica map[ID]struct{}
}