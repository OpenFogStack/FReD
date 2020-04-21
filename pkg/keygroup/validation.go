package keygroup

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/errors"
)

// CheckKeygroup checks every keygroup for errors in the name
func CheckKeygroup(params ...Keygroup) error {
	for _, p := range params {
		if p.Name == "" {
			return errors.New(errors.StatusBadRequest, "empty name")
		}
	}

	return nil
}
