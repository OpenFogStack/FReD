package fred

import (
	"github.com/go-errors/errors"
)

func checkItem(params ...Item) error {
	for _, p := range params {
		if p.Keygroup == "" {
			return errors.Errorf("checkItem failed for item %#v because the keygroup is empty", p)
		}

		if p.ID == "" {
			return errors.Errorf("checkItem failed for item %#v because the ID is empty", p)
		}
	}

	return nil
}

func checkKeygroup(params ...KeygroupName) error {
	for _, p := range params {
		if string(p) == "" {
			return errors.Errorf("checkKeygroup failed for item %#v because the keygroup is empty", p)
		}
	}

	return nil
}

func checkID(params ...string) error {
	for _, p := range params {
		if p == "" {
			return errors.Errorf("checkID failed for item %#v because the id is an empty string", p)
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
