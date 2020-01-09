package keygroup

import (
	"github.com/rs/zerolog/log"
)

type service struct {
	store  Store
	nodeID string
}

// New creates a new keygroup management service.
func New(k Store, n string) Service {
	return &service{
		store:  k,
		nodeID: n,
	}
}

// Create creates a new keygroup in the fred system.
func (s *service) Create(k Keygroup) error {
	err := checkKeygroup(k)

	if err != nil {
		log.Err(err).Msg("Keygroup service can not create a new keygroup")
		return err
	}

	err = s.store.Create(k)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes a keygroup from the fred system.
func (s *service) Delete(k Keygroup) error {
	err := checkKeygroup(k)

	if err != nil {
		log.Err(err).Msg("Keygroup service can not delete a keygroup")
		return err
	}

	err = s.store.Delete(k)

	if err != nil {
		log.Err(err).Msg("Keygroup service can not delete a keygroup")
		return err
	}

	return nil
}

// Exists checks if the given keygroup exists on this fred node.
func (s *service) Exists(k Keygroup) bool {
	return s.store.Exists(k)
}