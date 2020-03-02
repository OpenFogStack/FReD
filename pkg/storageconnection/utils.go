package storage

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
)

// RPCItemToDataItem transforms a RPC Item from the storage.proto spec to a data.Item
func RPCItemToDataItem(item *Item) *data.Item {
	kg := KeygroupStringToObject(item.Keygroup)
	id := item.Id
	idata := item.Data
	return &data.Item{Keygroup: kg, ID: id, Data: idata}
}

// DataItemToRPCItem transforms a data.Item to a storage.proto Item
func DataItemToRPCItem(item *data.Item) *Item {
	return &Item{Keygroup: string(item.Keygroup), Id: item.ID, Data: item.Data}
}

// DataItemToRPCKey transforms a data.Item to a storage.proto Key
func DataItemToRPCKey(item *data.Item) *Key {
	return &Key{Keygroup: string(item.Keygroup), Id: item.ID}
}

// KeygroupStringToObject transforms a string to a commons.KeygroupName
func KeygroupStringToObject(kg string) commons.KeygroupName {
	return commons.KeygroupName(kg)
}

// KeygroupObjectToString transforms a commons.Keygroupname to a string
func KeygroupObjectToString(kg commons.KeygroupName) string {
	return string(kg)
}

// RPCKeyToItem transforms a storage.proto Key to a data.Item
func RPCKeyToItem(key *Key) *data.Item {
	return &data.Item{Keygroup: KeygroupStringToObject(key.Keygroup), ID: key.Id}
}
