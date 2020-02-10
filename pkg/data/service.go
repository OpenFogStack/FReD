package data

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
)

// Service provides methods to manipulate items in the key-value store.
type Service interface {
	Read(kg commons.KeygroupName, id string) (string, error)
	ReadAll(kg commons.KeygroupName) ([]Item, error)
	Update(i Item) error
	Delete(kg commons.KeygroupName, id string) error
	CreateKeygroup(kg commons.KeygroupName) error
	DeleteKeygroup(kg commons.KeygroupName) error
}
