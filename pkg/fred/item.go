package fred

// Item is an item in the key-value store.
type Item struct {
	Keygroup KeygroupName
	ID       string
	Val      string
}
