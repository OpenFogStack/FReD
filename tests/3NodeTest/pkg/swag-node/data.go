package node

import (
	"context"

	"github.com/rs/zerolog/log"

	client "gitlab.tu-berlin.de/mcc-fred/fred/ext/go-client"
)

// PutItem puts a key-value pair into a (already created) keygroup.
func (n *Node) PutItem(kgname, item string, data string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Put for Item %s in KG %s; expecting %d", item, kgname, expectedStatusCode)

	resp, recvErr := n.Client.DataApi.KeygroupGroupIdDataItemIdPut(context.Background(), kgname, item, client.Item{
		Id:       item,
		Value:    data,
		Keygroup: kgname,
	})

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("PutItem: %s", err.Error())
		n.Errors++
	}
}

// GetItem returns the stored item.
func (n *Node) GetItem(kgname, item string, expectedStatusCode int, expectEmptyResponse bool) string {
	log.Debug().Str("node", n.URL).Msgf("Sending a Get for Item %s in KG %s; expecting %d", item, kgname, expectedStatusCode)

	data, resp, recvErr := n.Client.DataApi.KeygroupGroupIdDataItemIdGet(context.Background(), kgname, item)

	if err := checkResponse(resp, recvErr, data, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("GetItem: %s", err.Error())
		n.Errors++
	}

	return data.Value
}

// DeleteItem deletes the item from the keygroup.
func (n *Node) DeleteItem(kgname, item string, expectedStatusCode int, expectEmptyResponse bool) {
	log.Debug().Str("node", n.URL).Msgf("Sending a Delete for Item %s in KG %s; expecting %d", item, kgname, expectedStatusCode)

	resp, recvErr := n.Client.DataApi.KeygroupGroupIdDataItemIdDelete(context.Background(), kgname, item)

	if err := checkResponse(resp, recvErr, nil, expectedStatusCode, expectEmptyResponse); err != nil {
		log.Warn().Str("node", n.URL).Msgf("DeleteItem: %s", err.Error())
		n.Errors++
	}
}
