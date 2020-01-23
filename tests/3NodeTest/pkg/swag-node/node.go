package node

import (
	client "gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/vendor/go-client"
)

// Node represents the API to a single FReD Node
type Node struct {
	URL    string
	Errors int
	Client *client.APIClient
}

// NewNode creates a new Node with the specified url (should have format: http://%s:%d/v%d/)
func NewNode(url string) (node *Node) {
	cfg := client.NewConfiguration()
	cfg.BasePath = url

	c := client.NewAPIClient(cfg)

	node = &Node{URL: url, Errors: 0, Client: c}
	return
}
