package data

import (
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/errors"
)

func checkItem(params ...Item) error {
	for _, p := range params {
		if p.Keygroup == "" {
			log.Error().Msgf("checkItem failed for item %#v because the keygroup is empty", p)
			return errors.New(errors.StatusBadRequest, "data: empty keygroup")
		}

		if p.ID == "" {
			log.Error().Msgf("checkItem failed for item %#v because the ID is empty", p)
			return errors.New(errors.StatusBadRequest, "data: empty ID")
		}
	}

	return nil
}

func checkKeygroup(params ...Item) error {
	for _, p := range params {
		if p.Keygroup == "" {
			log.Error().Msgf("checkKeygroup failed for item %#v because the keygroup is empty", p)
			return errors.New(errors.StatusBadRequest, "data: empty keygroup")
		}
	}

	return nil
}
