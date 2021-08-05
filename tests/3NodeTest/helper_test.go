package threenodetest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"git.tu-berlin.de/mcc-fred/fred/proto/client"
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// node represents the API to a single FReD node
type node struct {
	id     string
	client client.ClientClient
	conn   *grpc.ClientConn
	addr   string
}

// newNode creates a new node that represents a connection to a single fred instance
func newNode(connAddr string, port int, actualAddr string, id string, certFile string, keyFile string, caFile string) *node {

	if certFile == "" {
		log.Fatal().Msg("node: no certificate file given")
	}

	if keyFile == "" {
		log.Fatal().Msg("node: no key file given")
	}

	if caFile == "" {
		log.Fatal().Msg("node: no root certificate file given")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("node: Cannot load certificates")
		return nil
	}

	// from https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("node: Failed to append %q to RootCAs: %v", caFile, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Fatal().Err(err).Msgf("node: No certs appended, using system certs only")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      rootCAs,
	}

	tc := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", connAddr, port), grpc.WithTransportCredentials(tc))

	if err != nil {
		log.Fatal().Err(err).Msg("node: Cannot create Grpc connection")
		return nil
	}

	c := client.NewClientClient(conn)

	return &node{
		client: c,
		id:     id,
		conn:   conn,
		addr:   actualAddr,
	}
}

// Close closes the connection
func (n *node) close() {
	err := n.conn.Close()

	log.Err(err).Msg("error closing node")
}

// createKeygroup calls the createKeygroup endpoint of the GRPC interface.
func (n *node) createKeygroup(kgname string, mutable bool, expiry int) error {
	_, err := n.client.CreateKeygroup(context.Background(), &client.CreateKeygroupRequest{Keygroup: kgname, Mutable: mutable, Expiry: int64(expiry)})

	if err != nil {
		return err
	}

	return nil
}

// deleteKeygroup calls the deleteKeygroup endpoint of the GRPC interface.
func (n *node) deleteKeygroup(kgname string) error {
	_, err := n.client.DeleteKeygroup(context.Background(), &client.DeleteKeygroupRequest{Keygroup: kgname})

	if err != nil {
		return err
	}

	return nil

}

// getKeygroupReplica calls the getKeygroupReplica endpoint of the GRPC interface.
func (n *node) getKeygroupReplica(kgname string) (map[string]int, error) {
	res, err := n.client.GetKeygroupReplica(context.Background(), &client.GetKeygroupReplicaRequest{Keygroup: kgname})

	if err != nil {
		return nil, err
	}

	nodes := make(map[string]int)

	for _, node := range res.Replica {
		nodes[node.NodeId] = int(node.Expiry)
	}

	return nodes, nil
}

// addKeygroupReplica calls the addKeygroupReplica endpoint of the GRPC interface.
func (n *node) addKeygroupReplica(kgname, replicaNodeID string, expiry int) error {
	_, err := n.client.AddReplica(context.Background(), &client.AddReplicaRequest{
		Keygroup: kgname,
		NodeId:   replicaNodeID,
		Expiry:   int64(expiry),
	})

	if err != nil {
		return err
	}

	return nil

}

// deleteKeygroupReplica calls the deleteKeygroupReplica endpoint of the GRPC interface.
func (n *node) deleteKeygroupReplica(kgname, replicaNodeID string) error {
	_, err := n.client.RemoveReplica(context.Background(), &client.RemoveReplicaRequest{
		Keygroup: kgname,
		NodeId:   replicaNodeID,
	})

	if err != nil {
		return err
	}

	return nil
}

// getKeygroupTriggers calls the getKeygroupTriggers endpoint of the GRPC interface.
func (n *node) getKeygroupTriggers(kgname string) ([]*client.Trigger, error) {
	res, err := n.client.GetKeygroupTriggers(context.Background(), &client.GetKeygroupTriggerRequest{Keygroup: kgname})

	if err != nil {
		return nil, err
	}

	return res.Triggers, nil
}

// addKeygroupTrigger calls the addKeygroupTrigger endpoint of the GRPC interface.
func (n *node) addKeygroupTrigger(kgname, triggerNodeID, triggerNodeHost string) error {
	_, err := n.client.AddTrigger(context.Background(), &client.AddTriggerRequest{Keygroup: kgname, TriggerId: triggerNodeID, TriggerHost: triggerNodeHost})

	if err != nil {
		return err
	}

	return nil
}

// deleteKeygroupTrigger calls the deleteKeygroupTrigger endpoint of the GRPC interface.
func (n *node) deleteKeygroupTrigger(kgname, triggerNodeID string) error {
	_, err := n.client.RemoveTrigger(context.Background(), &client.RemoveTriggerRequest{Keygroup: kgname, TriggerId: triggerNodeID})

	if err != nil {
		return err
	}

	return nil
}

// getAllReplica calls the getAllReplica endpoint of the GRPC interface.
func (n *node) getAllReplica() (map[string]string, error) {
	res, err := n.client.GetAllReplica(context.Background(), &client.Empty{})

	if err != nil {
		return nil, err
	}

	ids := make(map[string]string)

	for _, r := range res.Replicas {
		ids[r.NodeId] = r.Host
	}

	return ids, nil
}

// getReplica calls the getReplica endpoint of the GRPC interface.
func (n *node) getReplica(nodeID string) (string, string, error) {
	res, err := n.client.GetReplica(context.Background(), &client.GetReplicaRequest{NodeId: nodeID})

	if err != nil {
		return "", "", err
	}

	return res.NodeId, res.Host, nil
}

// putItem calls the putItem endpoint of the GRPC interface.
func (n *node) putItem(kgname, item string, data string) (vclock.VClock, error) {
	res, err := n.client.Update(context.Background(), &client.UpdateRequest{
		Keygroup: kgname,
		Id:       item,
		Data:     data,
	})

	if err != nil {
		return nil, err
	}

	return res.Version.Version, nil
}

// putItemVersion calls the putItem endpoint of the GRPC interface with a version.
func (n *node) putItemVersion(kgname, item string, data string, version vclock.VClock) (vclock.VClock, error) {
	res, err := n.client.Update(context.Background(), &client.UpdateRequest{
		Keygroup: kgname,
		Id:       item,
		Data:     data,
		Versions: []*client.Version{{
			Version: version.GetMap(),
		}}})

	if err != nil {
		return nil, err
	}

	return res.Version.Version, nil
}

// appendItem calls the appendItem endpoint of the GRPC interface.
func (n *node) appendItem(kgname string, id uint64, data string) (string, error) {
	res, err := n.client.Append(context.Background(), &client.AppendRequest{
		Keygroup: kgname,
		Id:       id,
		Data:     data,
	})

	if err != nil {
		return "", err
	}

	return res.Id, nil
}

// getItem calls the getItem endpoint of the GRPC interface.
func (n *node) getItem(kgname, item string) ([]string, []vclock.VClock, error) {
	res, err := n.client.Read(context.Background(), &client.ReadRequest{
		Keygroup: kgname,
		Id:       item,
	})

	if err != nil {
		return nil, nil, err
	}

	vals := make([]string, len(res.Data))
	vvectors := make([]vclock.VClock, len(res.Data))

	for i, it := range res.Data {
		vals[i] = it.Val
		vvectors[i] = it.Version
	}

	return vals, vvectors, nil
}

// ScanItems calls the Scan endpoint of the GRPC interface.
func (n *node) scanItems(kgname, item string, count uint64) (map[string]string, error) {
	res, err := n.client.Scan(context.Background(), &client.ScanRequest{
		Keygroup: kgname,
		Id:       item,
		Count:    count,
	})

	if res == nil {
		return nil, err
	}

	items := make(map[string]string)

	for _, d := range res.Data {
		items[d.Id] = d.Val
	}

	return items, nil
}

// deleteItem calls the deleteItem endpoint of the GRPC interface.
func (n *node) deleteItem(kgname, item string) error {

	_, err := n.client.Delete(context.Background(), &client.DeleteRequest{
		Keygroup: kgname,
		Id:       item,
	})

	if err != nil {
		return err
	}

	return nil

}

// deleteItemVersion calls the deleteItem endpoint of the GRPC interface with a version.
func (n *node) deleteItemVersion(kgname, item string, version vclock.VClock) (vclock.VClock, error) {

	res, err := n.client.Delete(context.Background(), &client.DeleteRequest{
		Keygroup: kgname,
		Id:       item,
		Versions: []*client.Version{{
			Version: version.GetMap(),
		}},
	})

	if err != nil {
		return nil, err
	}

	return res.Version.Version, nil

}

func strToRole(role string) client.UserRole {
	switch role {
	case "ReadKeygroup":
		return client.UserRole_ReadKeygroup
	case "WriteKeygroup":
		return client.UserRole_WriteKeygroup
	case "ConfigureReplica":
		return client.UserRole_ConfigureReplica
	case "ConfigureTrigger":
		return client.UserRole_ConfigureTrigger
	case "ConfigureKeygroups":
		return client.UserRole_ConfigureKeygroups
	}
	return -1
}

// addUser calls the addUser endpoint of the GRPC interface.
func (n *node) addUser(user string, kgname string, role string) error {
	_, err := n.client.AddUser(context.Background(), &client.AddUserRequest{
		User:     user,
		Keygroup: kgname,
		Role:     strToRole(role),
	})
	if err != nil {
		return err
	}

	return nil
}

// removeUser calls the removeUser endpoint of the GRPC interface.
func (n *node) removeUser(user string, kgname string, role string) error {
	_, err := n.client.RemoveUser(context.Background(), &client.RemoveUserRequest{
		User:     user,
		Keygroup: kgname,
		Role:     strToRole(role),
	})

	if err != nil {
		return err
	}

	return nil
}
