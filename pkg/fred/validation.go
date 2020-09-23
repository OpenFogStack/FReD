package fred

import (
	"regexp"

	"github.com/go-errors/errors"
)

var expr = "^[a-zA-Z0-9]+$"
var reg = regexp.MustCompile(expr)

func checkItem(params ...Item) error {
	for _, p := range params {
		if !reg.MatchString(string(p.Keygroup)) {
			return errors.Errorf("checkItem failed for item %#v because the keygroup name does not match %s", p, expr)
		}

		if !reg.MatchString(p.ID) {
			return errors.Errorf("checkItem failed for item %#v because the ID does not match %s", p, expr)
		}
	}

	return nil
}

func checkKeygroup(params ...KeygroupName) error {
	for _, p := range params {
		if !reg.MatchString(string(p)) {
			return errors.Errorf("checkKeygroup failed for keygroup %s because the keygroup name does not match %s", p, expr)
		}
	}

	return nil
}

func checkID(params ...string) error {
	for _, p := range params {
		if !reg.MatchString(p) {
			return errors.Errorf("checkID failed for item %s because the id does not match %s", p, expr)
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
