package data

import (
	"errors"
)

func checkItem(params ...Item) error {
	for _, p := range params {
		if p.Keygroup == "" {
			return errors.New("data: empty keygroup")
		}

		if p.ID == "" {
			return errors.New("data: empty ID")
		}
	}

	return nil
}

func checkKeygroup(params ...Item) error {
	for _, p := range params {
		if p.Keygroup == "" {
			return errors.New("data: empty keygroup")
		}
	}

	return nil
}
