package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type authService struct {
	n NameService
}

// newAuthService creates a new authorization service.
func newAuthService(n NameService) *authService {
	return &authService{
		n: n,
	}
}

func (a *authService) addRoles(u string, r []Role, k KeygroupName) error {
	log.Debug().Msgf("adding roles %+v for user %s for keygroup %s", r, u, k)
	for _, role := range r {
		for m := range permissions[role] {

			err := a.n.AddUserPermissions(u, m, k)

			if err != nil {
				return errors.New(err)
			}
		}
	}

	return nil
}

func (a *authService) revokeRoles(u string, r []Role, k KeygroupName) error {
	log.Debug().Msgf("removing roles %+v from user %s for keygroup %s", r, u, k)
	for _, role := range r {
		for m := range permissions[role] {
			err := a.n.RevokeUserPermissions(u, m, k)

			if err != nil {
				return errors.New(err)
			}
		}
	}

	return nil
}

func (a *authService) isAllowed(u string, m Method, k KeygroupName) (bool, error) {
	log.Debug().Msgf("checking if user %s is allowed to perform %s on keygroup %s...", u, m, k)
	p, err := a.n.GetUserPermissions(u, k)

	if err != nil {
		return false, errors.New(err)
	}

	_, ok := p[m]

	// Only compute the string if log level is debug
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		var res string
		if ok {
			res = "true"
		} else {
			res = "false"
		}
		log.Debug().Msgf("...user is allowed: %s", res)
	}

	return ok, nil

}
