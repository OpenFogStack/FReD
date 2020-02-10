package data

import (
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
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

func checkKeygroup(params ...commons.KeygroupName) error {
	for _, p := range params {
		if string(p) == "" {
			log.Error().Msgf("checkKeygroup failed for item %#v because the keygroup is empty", p)
			return errors.New(errors.StatusBadRequest, "data: empty keygroup")
		}
	}

	return nil
}

func checkID(params ...string) error {
	for _, p := range params {
		if p == "" {
			log.Error().Msgf("checkID failed for item %#v because the id is an empty string", p)
			return errors.New(errors.StatusBadRequest, "data: empty keygroup")
		}
	}

	return nil
}

// Check a single keygroup and id at once. If both error it returns the error of the KG
func checkKGandID(kg commons.KeygroupName, id string) error {
	err := checkKeygroup(kg)
	if err == nil {
		err = checkID(id)
	}
	return err
}
