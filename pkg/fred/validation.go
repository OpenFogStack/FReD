package fred

import (
	"github.com/rs/zerolog/log"
)

func checkItem(params ...Item) error {
	for _, p := range params {
		if p.Keygroup == "" {
			log.Error().Msgf("checkItem failed for item %#v because the keygroup is empty", p)
			return newError(StatusBadRequest, "val: empty keygroup")
		}

		if p.ID == "" {
			log.Error().Msgf("checkItem failed for item %#v because the ID is empty", p)
			return newError(StatusBadRequest, "val: empty ID")
		}
	}

	return nil
}

func checkKeygroup(params ...KeygroupName) error {
	for _, p := range params {
		if string(p) == "" {
			log.Error().Msgf("checkKeygroup failed for item %#v because the keygroup is empty", p)
			return newError(StatusBadRequest, "val: empty keygroup")
		}
	}

	return nil
}

func checkID(params ...string) error {
	for _, p := range params {
		if p == "" {
			log.Error().Msgf("checkID failed for item %#v because the id is an empty string", p)
			return newError(StatusBadRequest, "val: empty keygroup")
		}
	}

	return nil
}

// Check a single keygroup and id at once. If both error it returns the error of the KG
func checkKGandID(kg KeygroupName, id string) error {
	err := checkKeygroup(kg)
	if err == nil {
		err = checkID(id)
	}
	return err
}
