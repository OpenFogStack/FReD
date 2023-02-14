package alexandra

import (
	"fmt"
	"sync"

	"git.tu-berlin.de/mcc-fred/vclock"
	"github.com/rs/zerolog/log"
)

type cache struct {
	keygroups map[string]struct {
		items map[string]*struct {
			clocks []vclock.VClock
			*sync.Mutex
		}
		*sync.RWMutex
	}
	*sync.RWMutex
}

func newCache() *cache {
	return &cache{
		keygroups: make(map[string]struct {
			items map[string]*struct {
				clocks []vclock.VClock
				*sync.Mutex
			}
			*sync.RWMutex
		}),
		RWMutex: &sync.RWMutex{},
	}
}

func (c *cache) cLock(kg string, id string) {
	c.RLock()
	if _, ok := c.keygroups[kg]; !ok {
		c.RUnlock()
		c.Lock()
		if _, ok := c.keygroups[kg]; !ok {
			c.keygroups[kg] = struct {
				items map[string]*struct {
					clocks []vclock.VClock
					*sync.Mutex
				}
				*sync.RWMutex
			}{
				items: make(map[string]*struct {
					clocks []vclock.VClock
					*sync.Mutex
				}),
				RWMutex: &sync.RWMutex{},
			}
		}
		c.Unlock()
		c.RLock()
	}
	c.keygroups[kg].RLock()
	if _, ok := c.keygroups[kg].items[id]; !ok {
		c.keygroups[kg].RUnlock()
		c.keygroups[kg].Lock()
		if _, ok := c.keygroups[kg].items[id]; !ok {
			c.keygroups[kg].items[id] = &struct {
				clocks []vclock.VClock
				*sync.Mutex
			}{
				clocks: make([]vclock.VClock, 0),
				Mutex:  &sync.Mutex{},
			}
		}
		c.keygroups[kg].Unlock()
		c.keygroups[kg].RLock()
	}
	c.keygroups[kg].items[id].Lock()
}

func (c *cache) cUnlock(kg string, id string) {
	c.keygroups[kg].items[id].Unlock()
	c.keygroups[kg].RUnlock()
	c.RUnlock()
}

func (c *cache) add(kg string, id string, version vclock.VClock) error {
	c.cLock(kg, id)
	defer c.cUnlock(kg, id)

	newClocks := make([]vclock.VClock, 0, len(c.keygroups[kg].items[id].clocks))

	// let's first go to the existing versions
	for _, v := range c.keygroups[kg].items[id].clocks {
		switch v.Order(version) {
		// seems like we have that version in store already!
		// just return, nothing to do here
		case vclock.Equal:
			{
				return nil
			}
		// if there is a concurrent version to our new version, we need to keep that
		case vclock.Concurrent:
			{
				newClocks = append(newClocks, v)
				continue
			}
		// if by chance we come upon a newer version, this has all been pointless
		// actually, this is a bad sign: we seem to have read outdated data!
		case vclock.Ancestor:
			{
				return fmt.Errorf("add failed because your version %v is older than cached version %v", version, v)
			}
		}
	}

	// now actually add the new version
	newClocks = append(newClocks, version.Copy())

	log.Debug().Msgf("Cache Add: new: %+v old: %+v", newClocks, c.keygroups[kg].items[id].clocks)

	// and store our new clocks
	c.keygroups[kg].items[id].clocks = newClocks
	return nil
}

func (c *cache) supersede(kg string, id string, known []vclock.VClock, version vclock.VClock) error {
	c.cLock(kg, id)
	defer c.cUnlock(kg, id)

	newClocks := []vclock.VClock{
		version.Copy(),
	}

	// add all clocks from cache as well, unless they are in the "known" array
	for _, v := range c.keygroups[kg].items[id].clocks {
		discard := false
		for _, k := range known {
			if v.Compare(k, vclock.Equal) {
				discard = true
				break
			}
		}
		if discard {
			continue
		}

		newClocks = append(newClocks, v)
	}

	c.keygroups[kg].items[id].clocks = newClocks

	return nil
}

func (c *cache) get(kg string, id string) ([]vclock.VClock, error) {
	c.cLock(kg, id)
	defer c.cUnlock(kg, id)

	return c.keygroups[kg].items[id].clocks, nil
}
