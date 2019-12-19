package app

import "errors"

func checkParameters(params ...string) error{
	for _, p := range params {
		if p == "" {
			return errors.New("empty parameter")
		}
	}

	return nil
}
