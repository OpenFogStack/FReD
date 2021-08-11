package alexandra

import (
	"sync"

	"github.com/DistributedClocks/GoVector/govec/vclock"
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

	for _, v := range c.keygroups[kg].items[id].clocks {
		if version.Compare(v, vclock.Concurrent) {
			newClocks = append(newClocks, v)
			continue
		}
		if version.Compare(v, vclock.Descendant) {
			return nil
		}
	}

	c.keygroups[kg].items[id].clocks = newClocks
	return nil
}

func (c *cache) supersede(kg string, id string, version vclock.VClock) error {
	c.cLock(kg, id)
	defer c.cUnlock(kg, id)

	c.keygroups[kg].items[id].clocks = []vclock.VClock{
		version.Copy(),
	}

	return nil
}

func (c *cache) get(kg string, id string) ([]vclock.VClock, error) {
	c.cLock(kg, id)
	defer c.cUnlock(kg, id)

	return c.keygroups[kg].items[id].clocks, nil
}
