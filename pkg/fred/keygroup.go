package fred

// Keygroup has a name and a list of replica nodes and trigger nodes.
type Keygroup struct {
	Name KeygroupName
}

// KeygroupName is a name of a keygroup.
type KeygroupName string
