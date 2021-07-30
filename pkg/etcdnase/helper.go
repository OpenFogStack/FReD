package etcdnase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"go.etcd.io/etcd/client/v3"
)

// getPrefix gets every key that starts(!) with the specified string
// the keys are sorted ascending by key for easier debugging
func (n *NameService) getPrefix(prefix string) (kv map[string]string, err error) {
	// the hard part of caching isn't storing a key-value pair locally
	// it's actually knowing when to remove an entry from the cache because it's outdated
	// sure, you can set a timeout or other eviction policy but that's more or less arbitrary
	// we remove an item from the cache if we delete it from the nase ourselves or nase informs us about deletion via watchers
	// we update an item if we update it ourselves or nase informs us about an update via watchers
	// prefixes are the hardest part about this
	if n.cached {
		// let's check the local cache first
		// store prefix directly in cache
		val, ok := n.local.Get(prefix)

		// found something!
		if ok {
			log.Debug().Msgf("prefix: %s cache hit", prefix)
			return val.(map[string]string), nil
		}

		log.Debug().Msgf("prefix: %s cache miss", prefix)
	}

	// didn't find anything? ask nameservice, cache, and be sure to invalidate on change
	ctx, cncl := context.WithTimeout(context.Background(), timeout)

	defer cncl()

	resp, err := n.cli.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))

	if err != nil {
		return nil, errors.New(err)
	}

	kv = make(map[string]string)

	for _, val := range resp.Kvs {
		kv[string(val.Key)] = string(val.Value)
	}

	if n.cached {
		n.local.Set(prefix, kv, 1)

		// we start the watcher directly
		go func() {
			watchCtx, watchCncl := context.WithCancel(context.Background())
			c := n.watcher.Watch(watchCtx, prefix, clientv3.WithPrefix())
			log.Debug().Msgf("nase cache: watching for changes to prefix %s", prefix)

			defer watchCncl()
			for r := range c {
				if err := r.Err(); err != nil {
					log.Err(err).Msgf("nase cache: error getting changes to prefix %s", prefix)
				}
				log.Debug().Msgf("nase cache: got %d changes to prefix %s", len(r.Events), prefix)
				val, ok := n.local.Get(prefix)

				if !ok {
					n.local.Del(prefix)
					return
				}

				prefixMap, ok := val.(map[string]string)

				if !ok {
					n.local.Del(prefix)
					return
				}

				for _, ev := range r.Events {

					if ev.Type == clientv3.EventTypeDelete {
						delete(prefixMap, string(ev.Kv.Key))
						log.Debug().Msgf("prefix: %s remote cache invalidation for key %s", prefix, string(ev.Kv.Key))
					}

					if ev.Type == clientv3.EventTypePut {
						prefixMap[string(ev.Kv.Key)] = string(ev.Kv.Value)
						log.Debug().Msgf("prefix: %s remote cache update for key %s", prefix, string(ev.Kv.Key))
					}
				}
				n.local.Set(prefix, prefixMap, 1)
			}
		}()
	}
	return kv, nil
}

// getExact gets the exact key
func (n *NameService) getExact(key string) (v string, err error) {

	if n.cached {
		// let's check the local cache first
		val, ok := n.local.Get(key)

		// found something!
		if ok {
			log.Debug().Msgf("key: %s cache hit", key)
			return val.(string), nil
		}

		log.Debug().Msgf("key: %s cache miss", key)
	}

	// didn't find anything? ask nameservice, cache, and be sure to invalidate on change
	ctx, cncl := context.WithTimeout(context.Background(), timeout)

	defer cncl()

	resp, err := n.cli.Get(ctx, key)

	if err != nil {
		return "", errors.New(err)
	}

	if len(resp.Kvs) != 0 {
		v = string(resp.Kvs[0].Value)
	}

	if n.cached {
		n.local.Set(key, v, 1)

		// we can always assume that no watcher exists for this key because otherwise we would have had a cache hit
		// so we can safely start a new watcher
		// that watcher should exit on key deletion, though!
		go func() {
			watchCtx, watchCncl := context.WithCancel(context.Background())
			c := n.watcher.Watch(watchCtx, key)

			defer watchCncl()
			for r := range c {
				if err := r.Err(); err != nil {
					log.Err(err).Msgf("nase cache: error getting changes to key %s", key)
				}
				log.Debug().Msgf("nase cache: got %d changes to key %s", len(r.Events), key)
				for _, ev := range r.Events {
					if ev.Type == clientv3.EventTypeDelete {
						n.local.Del(key)
						log.Debug().Msgf("key: %s remote cache invalidation", key)
						return
					}

					if ev.Type == clientv3.EventTypePut {
						n.local.Set(key, string(ev.Kv.Value), 1)
						log.Debug().Msgf("key: %s remote cache update", key)
					}
				}
			}
		}()

	}

	return v, nil
}

func (n *NameService) getKeygroupStatus(kg string) (string, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgStatusString, kg))

	return resp, err
}

func (n *NameService) getKeygroupMutable(kg string) (string, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgMutableString, kg))

	return resp, err
}

func (n *NameService) getKeygroupExpiry(kg string, id string) (int, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgExpiryStringPrefix, kg) + id)
	if resp == "" {
		return 0, err
	}

	return strconv.Atoi(resp)
}

// put puts the value into etcd.
func (n *NameService) put(key, value string, prefix ...string) (err error) {
	ctx, cncl := context.WithTimeout(context.TODO(), timeout)

	defer cncl()

	if n.cached {
		for _, p := range prefix {
			val, ok := n.local.Get(p)

			if !ok {
				n.local.Del(p)
				continue
			}

			prefixMap, ok := val.(map[string]string)

			if !ok {
				n.local.Del(p)
				continue
			}

			prefixMap[key] = value

			n.local.Set(p, prefixMap, 1)
			log.Debug().Msgf("prefix: %s local cache update for key %s", p, key)
		}
		n.local.Set(key, value, 1)
		log.Debug().Msgf("key: %s local cache update", key)
	}

	_, err = n.cli.Put(ctx, key, value)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// delete removes the value from etcd.
func (n *NameService) delete(key string, prefix ...string) (err error) {
	ctx, cncl := context.WithTimeout(context.TODO(), timeout)

	defer cncl()

	if n.cached {
		for _, p := range prefix {
			val, ok := n.local.Get(p)

			if !ok {
				n.local.Del(p)
				continue
			}

			prefixMap, ok := val.(map[string]string)

			if !ok {
				n.local.Del(p)
				continue
			}

			delete(prefixMap, key)

			n.local.Set(p, prefixMap, 1)
			log.Debug().Msgf("prefix: %s local cache invalidation for key %s", p, key)

		}
		n.local.Del(key)
		log.Debug().Msgf("key: %s local cache invalidation", key)
	}

	_, err = n.cli.Delete(ctx, key)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// addOwnKgNodeEntry adds the entry for this node with a status.
func (n *NameService) addOwnKgNodeEntry(kg string, status string) error {
	prefix, id := n.fmtKgNode(kg)
	return n.put(prefix+id, status, prefix)
}

// addOtherKgNodeEntry adds the entry for a remote node with a status.
func (n *NameService) addOtherKgNodeEntry(node string, kg string, status string) error {
	prefix := fmt.Sprintf(fmtKgNodeStringPrefix, kg)
	key := prefix + node
	return n.put(key, status, prefix)
}

// addKgStatusEntry adds the entry for a (new!) keygroup with a status.
func (n *NameService) addKgStatusEntry(kg string, status string) error {
	return n.put(fmt.Sprintf(fmtKgStatusString, kg), status)
}

// addKgMutableEntry adds the ismutable entry for a keygroup with a status.
func (n *NameService) addKgMutableEntry(kg string, mutable bool) error {
	var data string

	if mutable {
		data = "true"
	} else {
		data = "false"
	}

	return n.put(fmt.Sprintf(fmtKgMutableString, kg), data)
}

// addKgExpiryEntry adds the expiry entry for a keygroup with a status.
func (n *NameService) addKgExpiryEntry(kg string, id string, expiry int) error {
	prefix := fmt.Sprintf(fmtKgExpiryStringPrefix, kg)
	return n.put(prefix+id, strconv.Itoa(expiry), prefix)
}

// fmtKgNode turns a keygroup name into the key that this node will save its state in
// Currently: kg|[keygroup]|node|[NodeID]
func (n *NameService) fmtKgNode(kg string) (string, string) {
	prefix := fmt.Sprintf(fmtKgNodeStringPrefix, kg)
	return prefix, n.NodeID
}

func getNodeNameFromKgNodeString(kgNode string) string {
	split := strings.Split(kgNode, sep)
	return split[len(split)-1]
}
