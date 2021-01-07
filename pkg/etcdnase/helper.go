package etcdnase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// getPrefix gets every key that starts(!) with the specified string
// the keys are sorted ascending by key for easier debugging
func (n *NameService) getPrefix(prefix string) (kv []*mvccpb.KeyValue, err error) {
	ctx, cncl := context.WithTimeout(context.Background(), timeout)

	defer cncl()

	resp, err := n.cli.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))

	if err != nil {
		return nil, errors.New(err)
	}

	kv = resp.Kvs

	return
}

// getExact gets the exact key
func (n *NameService) getExact(key string) (kv []*mvccpb.KeyValue, err error) {
	ctx, cncl := context.WithTimeout(context.Background(), timeout)

	defer cncl()

	resp, err := n.cli.Get(ctx, key)

	if err != nil {
		return nil, errors.New(err)
	}

	kv = resp.Kvs
	return
}

func (n *NameService) getKeygroupStatus(kg string) (string, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgStatusString, kg))
	if resp == nil {
		return "", err
	}
	return string(resp[0].Value), err
}

func (n *NameService) getKeygroupMutable(kg string) (string, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgMutableString, kg))
	if resp == nil {
		return "", err
	}
	return string(resp[0].Value), err
}

func (n *NameService) getKeygroupExpiry(kg string, id string) (int, error) {
	resp, err := n.getExact(fmt.Sprintf(fmtKgExpiryString, kg, id))
	if resp == nil {
		return 0, err
	}

	return strconv.Atoi(string(resp[0].Value))
}

// put puts the value into etcd.
func (n *NameService) put(key, value string) (err error) {
	ctx, cncl := context.WithTimeout(context.TODO(), timeout)

	defer cncl()

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
