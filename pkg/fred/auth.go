package fred

import (
	"github.com/go-errors/errors"
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
	for _, role := range r {
		for m := range permissions[role] {
			log.Debug().Msgf("adding permission %s from user %s for keygroup %s", m, u, k)

			err := a.n.AddUserPermissions(u, m, k)

			if err != nil {
				return errors.New(err)
			}
		}
	}

	return nil
}

func (a *authService) revokeRoles(u string, r []Role, k KeygroupName) error {
	for _, role := range r {
		for m := range permissions[role] {
			log.Debug().Msgf("removing permission %s from user %s for keygroup %s", m, u, k)

			err := a.n.RevokeUserPermissions(u, m, k)

			if err != nil {
				return errors.New(err)
			}
		}
	}

	return nil
}

func (a *authService) isAllowed(u string, m Method, k KeygroupName) (bool, error) {
	log.Debug().Msgf("checking if user %s is allowed to perform %s on keygroup %s", u, m, k)
	p, err := a.n.GetUserPermissions(u, k)

	log.Debug().Msgf("result of check if user %s is allowed to perform %s on keygroup %s: %v", u, m, k, p)

	if err != nil {
		return false, errors.New(err)
	}

	_, ok := p[m]

	return ok, nil

}
