package keygroup

import "errors"

func checkKeygroup(params ...Keygroup) error {
	for _, p := range params {
		if p.Name == "" {
			return errors.New("empty name")
		}
	}

	return nil
}
