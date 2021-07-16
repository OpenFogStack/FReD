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
	if n.cached {
		// let's check the local cache first
		// store prefix directly in cache
		val, ok := n.local.Get(prefix)

		// found something!
		if ok {
			log.Debug().Msgf("prefix: %s cache hit", prefix)
			return val.(map[string]string), nil
		}
	}

	log.Debug().Msgf("prefix: %s cache miss", prefix)

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

		go func() {
			watchCtx, watchCncl := context.WithCancel(context.Background())
			c := n.watcher.Watch(watchCtx, prefix, clientv3.WithPrefix())

			defer watchCncl()
			for r := range c {
				for _, ev := range r.Events {
					if ev.IsModify() {
						n.local.Del(prefix)
						log.Debug().Msgf("prefix: %s remote cache invalidation", prefix)
						return
					}
				}
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
	}

	log.Debug().Msgf("key: %s cache miss", key)

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

		go func() {
			watchCtx, watchCncl := context.WithCancel(context.Background())
			c := n.watcher.Watch(watchCtx, key)

			defer watchCncl()
			for r := range c {
				for _, ev := range r.Events {
					if ev.IsModify() {
						n.local.Del(key)
						log.Debug().Msgf("key: %s remote cache invalidation", key)
						return

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
	resp, err := n.getExact(fmt.Sprintf(fmtKgExpiryString, kg, id))
	if resp == "" {
		return 0, err
	}

	return strconv.Atoi(resp)
}

// put puts the value into etcd.
func (n *NameService) put(key, value string) (err error) {
	ctx, cncl := context.WithTimeout(context.TODO(), timeout)

	defer cncl()

	if n.cached {
		n.local.Del(key)
		log.Debug().Msgf("key: %s local cache invalidation", key)
	}

	_, err = n.cli.Put(ctx, key, value)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// delete removes the value from etcd.
func (n *NameService) delete(key string) (err error) {
	ctx, cncl := context.WithTimeout(context.TODO(), timeout)

	defer cncl()

	if n.cached {
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
	return n.put(n.fmtKgNode(kg), status)
}

// addOtherKgNodeEntry adds the entry for a remote node with a status.
func (n *NameService) addOtherKgNodeEntry(node string, kg string, status string) error {
	key := fmt.Sprintf(fmtKgNodeString, kg, node)
	return n.put(key, status)
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
	return n.put(fmt.Sprintf(fmtKgExpiryString, kg, id), strconv.Itoa(expiry))
}

// fmtKgNode turns a keygroup name into the key that this node will save its state in
// Currently: kg|[keygroup]|node|[NodeID]
func (n *NameService) fmtKgNode(kg string) string {
	return fmt.Sprintf(fmtKgNodeString, kg, n.NodeID)
}

func getNodeNameFromKgNodeString(kgNode string) string {
	split := strings.Split(kgNode, sep)
	return split[len(split)-1]
}
