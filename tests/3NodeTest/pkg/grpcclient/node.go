package grpcclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"git.tu-berlin.de/mcc-fred/fred/proto/client"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Node represents the API to a single FReD Node
type Node struct {
	Errors int
	ID     string
	Client client.ClientClient
	conn   *grpc.ClientConn
	Addr   string
}

// NewNode creates a new Node that represents a connection to a single fred instance
func NewNode(addr string, port int, id string, certFile string, keyFile string, caFile string) *Node {

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load certificates")
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
		log.Fatal().Err(err).Msgf("Failed to append %q to RootCAs: %v", caFile, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Fatal().Err(err).Msgf("No certs appended, using system certs only")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      rootCAs,
	}

	tc := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", addr, port), grpc.WithTransportCredentials(tc))

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return nil
	}

	c := client.NewClientClient(conn)

	return &Node{
		Errors: 0,
		Client: c,
		ID:     id,
		conn:   conn,
		Addr:   addr,
	}
}

// Close closes the connection
func (n *Node) Close() {
	err := n.conn.Close()

	log.Err(err).Msg("error closing node")
}

// CreateKeygroup calls the CreateKeygroup endpoint of the GRPC interface.
func (n *Node) CreateKeygroup(kgname string, mutable bool, expiry int, expectError bool) {
	status, err := n.Client.CreateKeygroup(context.Background(), &client.CreateKeygroupRequest{Keygroup: kgname, Mutable: mutable, Expiry: int64(expiry)})

	if err != nil && !expectError {
		log.Warn().Msgf("CreateKeygroup: error %s", err)
		n.Errors++
		return
	}

	if err == nil && expectError {
		log.Warn().Msg("CreateKeygroup: Expected Error bot got no error :(")
		n.Errors++
		return
	}

	if err == nil && !expectError && status.Status == client.EnumStatus_ERROR {
		log.Warn().Msgf("CreateKeygroup: error %s with status %s", err, status.Status)
		n.Errors++
		return
	}
}

// DeleteKeygroup calls the DeleteKeygroup endpoint of the GRPC interface.
func (n *Node) DeleteKeygroup(kgname string, expectError bool) {
	status, err := n.Client.DeleteKeygroup(context.Background(), &client.DeleteKeygroupRequest{Keygroup: kgname})

	if err != nil && !expectError {
		log.Warn().Msgf("DeleteKeygroup: error %s", err)
		n.Errors++
		return
	}

	if err == nil && expectError {
		log.Warn().Msg("DeleteKeygroup: Expected Error bot got no error")
		n.Errors++
		return
	}

	if err == nil && status.Status == client.EnumStatus_ERROR && !expectError {
		log.Warn().Msgf("DeleteKeygroup: error %s with status %s", err, status.Status)
		n.Errors++
		return
	}
}

// GetKeygroupReplica calls the GetKeygroupReplica endpoint of the GRPC interface.
func (n *Node) GetKeygroupReplica(kgname string, expectError bool) map[string]int {
	res, err := n.Client.GetKeygroupReplica(context.Background(), &client.GetKeygroupReplicaRequest{Keygroup: kgname})

	if err != nil && !expectError {
		log.Warn().Msgf("GetKeygroupReplica: error %s", err)
		n.Errors++
	}

	if err == nil && expectError {
		log.Warn().Msg("GetKeygroupReplica: Expected Error bot got no error")
		n.Errors++
	}

	if res == nil {
		return nil
	}

	nodes := make(map[string]int)

	for _, node := range res.Replica {
		nodes[node.NodeId] = int(node.Expiry)
	}

	return nodes
}

// AddKeygroupReplica calls the AddKeygroupReplica endpoint of the GRPC interface.
func (n *Node) AddKeygroupReplica(kgname, replicaNodeID string, expiry int, expectError bool) {
	status, err := n.Client.AddReplica(context.Background(), &client.AddReplicaRequest{
		Keygroup: kgname,
		NodeId:   replicaNodeID,
		Expiry:   int64(expiry),
	})

	if err != nil && !expectError {
		log.Warn().Msgf("AddKeygroupReplica: error %s", err)
		n.Errors++
		return
	}

	if err == nil && expectError {
		log.Warn().Msg("AddKeygroupReplica: Expected Error bot got no error")
		n.Errors++
		return
	}

	if err == nil && status.Status == client.EnumStatus_ERROR && !expectError {
		log.Warn().Msgf("AddKeygroupReplica: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// DeleteKeygroupReplica calls the DeleteKeygroupReplica endpoint of the GRPC interface.
func (n *Node) DeleteKeygroupReplica(kgname, replicaNodeID string, expectError bool) {
	status, err := n.Client.RemoveReplica(context.Background(), &client.RemoveReplicaRequest{
		Keygroup: kgname,
		NodeId:   replicaNodeID,
	})

	if err != nil && !expectError {
		log.Warn().Msgf("DeleteKeygroupReplica: error %s", err)
		n.Errors++
		return
	}
	if err == nil && expectError {
		log.Warn().Msg("DeleteKeygroupReplica: Expected Error bot got no error")
		n.Errors++
		return
	}

	if err == nil && status.Status == client.EnumStatus_ERROR && !expectError {
		log.Warn().Msgf("DeleteKeygroupReplica: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// GetKeygroupTriggers calls the GetKeygroupTriggers endpoint of the GRPC interface.
func (n *Node) GetKeygroupTriggers(kgname string, expectError bool) []*client.Trigger {
	res, err := n.Client.GetKeygroupTriggers(context.Background(), &client.GetKeygroupTriggerRequest{Keygroup: kgname})

	if err != nil && !expectError {
		log.Warn().Msgf("GetKeygroupTriggers: error %s", err)
		n.Errors++
	}

	if err == nil && expectError {
		log.Warn().Msg("GetKeygroupTriggers: Expected Error bot got no error")
		n.Errors++
	}

	if res == nil {
		return nil
	}

	return res.Triggers
}

// AddKeygroupTrigger calls the AddKeygroupTrigger endpoint of the GRPC interface.
func (n *Node) AddKeygroupTrigger(kgname, triggerNodeID, triggerNodeHost string, expectError bool) {
	status, err := n.Client.AddTrigger(context.Background(), &client.AddTriggerRequest{Keygroup: kgname, TriggerId: triggerNodeID, TriggerHost: triggerNodeHost})

	if err != nil && !expectError {
		log.Warn().Msgf("AddKeygroupTrigger: error %s", err)
		n.Errors++
		return
	}

	if err == nil && expectError {
		log.Warn().Msg("AddKeygroupTrigger: Expected Error bot got no error")
		n.Errors++
		return
	}

	if err == nil && status.Status == client.EnumStatus_ERROR && !expectError {
		log.Warn().Msgf("AddKeygroupTrigger: error %s with status %s", err, status.Status)
		n.Errors++
	}

}

// DeleteKeygroupTrigger calls the DeleteKeygroupTrigger endpoint of the GRPC interface.
func (n *Node) DeleteKeygroupTrigger(kgname, triggerNodeID string, expectError bool) {
	status, err := n.Client.RemoveTrigger(context.Background(), &client.RemoveTriggerRequest{Keygroup: kgname, TriggerId: triggerNodeID})

	if err != nil && !expectError {
		log.Warn().Msgf("DeleteKeygroupTrigger: error %s", err)
		n.Errors++
		return
	}

	if err == nil && expectError {
		log.Warn().Msg("DeleteKeygroupTrigger: Expected Error bot got no error")
		n.Errors++
		return
	}

	if err == nil && status.Status == client.EnumStatus_ERROR && !expectError {
		log.Warn().Msgf("DeleteKeygroupTrigger: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// GetAllReplica calls the GetAllReplica endpoint of the GRPC interface.
func (n *Node) GetAllReplica(expectError bool) map[string]string {
	res, err := n.Client.GetAllReplica(context.Background(), &client.GetAllReplicaRequest{})
	if err != nil && !expectError {
		log.Warn().Msgf("GetAllReplicas: error %s", err)
		n.Errors++
	}

	if err == nil && expectError {
		log.Warn().Msg("GetAllReplicas: Expected Error bot got no error")
		n.Errors++
	}

	if res == nil {
		return nil
	}

	ids := make(map[string]string)

	for _, r := range res.Replicas {
		ids[r.NodeId] = r.Host
	}

	return ids
}

// GetReplica calls the GetReplica endpoint of the GRPC interface.
func (n *Node) GetReplica(nodeID string, expectError bool) string {
	res, err := n.Client.GetReplica(context.Background(), &client.GetReplicaRequest{NodeId: nodeID})

	if err != nil && !expectError {
		log.Warn().Msgf("GetReplica: error %s", err)
		n.Errors++
	}

	if err == nil && expectError {
		log.Warn().Msg("GetReplica: Expected Error bot got no error")
		n.Errors++
	}

	if res == nil {
		return ""
	}

	return res.NodeId
}

// PutItem calls the PutItem endpoint of the GRPC interface.
func (n *Node) PutItem(kgname, item string, data string, expectError bool) {
	status, err := n.Client.Update(context.Background(), &client.UpdateRequest{
		Keygroup: kgname,
		Id:       item,
		Data:     data,
	})

	if err != nil && !expectError {
		log.Warn().Msgf("Update: error %s", err)
		n.Errors++
		return
	}

	if err == nil && expectError {
		log.Warn().Msg("Update: Expected Error bot got no error")
		n.Errors++
		return
	}

	if err == nil && status.Status == client.EnumStatus_ERROR && !expectError {
		log.Warn().Msgf("Update: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// AppendItem calls the AppendItem endpoint of the GRPC interface.
func (n *Node) AppendItem(kgname string, data string, expectError bool) string {
	status, err := n.Client.Append(context.Background(), &client.AppendRequest{
		Keygroup: kgname,
		Data:     data,
	})

	if err != nil != expectError {
		log.Warn().Msgf("Append: error %s", err)
		n.Errors++
	}

	if err == nil && expectError {
		log.Warn().Msg("Append: Expected Error bot got no error")
		n.Errors++
	}

	if err != nil {
		return ""
	}

	return status.Id
}

// GetItem calls the GetItem endpoint of the GRPC interface.
func (n *Node) GetItem(kgname, item string, expectError bool) string {
	res, err := n.Client.Read(context.Background(), &client.ReadRequest{
		Keygroup: kgname,
		Id:       item,
	})

	if err != nil && !expectError {
		log.Warn().Msgf("GetItem: error %s", err)
		n.Errors++
	}

	if err == nil && expectError {
		log.Warn().Msg("GetItem: Expected Error bot got no error")
		n.Errors++
	}

	if res == nil {
		return ""
	}

	return res.Data
}

// ScanItems calls the Scan endpoint of the GRPC interface.
func (n *Node) ScanItems(kgname, item string, count uint64, expectError bool) map[string]string {
	res, err := n.Client.Scan(context.Background(), &client.ScanRequest{
		Keygroup: kgname,
		Id:       item,
		Count:    count,
	})

	if err != nil && !expectError {
		log.Warn().Msgf("ScanItems: error %s", err)
		n.Errors++
	}

	if err == nil && expectError {
		log.Warn().Msg("ScanItems: Expected Error bot got no error")
		n.Errors++
	}

	if res == nil {
		return nil
	}

	items := make(map[string]string)

	for _, d := range res.Data {
		items[d.Id] = d.Val
	}

	return items
}

// DeleteItem calls the DeleteItem endpoint of the GRPC interface.
func (n *Node) DeleteItem(kgname, item string, expectError bool) {

	status, err := n.Client.Delete(context.Background(), &client.DeleteRequest{
		Keygroup: kgname,
		Id:       item,
	})

	if err != nil && !expectError {
		log.Warn().Msgf("Delete: error %s", err)
		n.Errors++
		return
	}

	if err == nil && expectError {
		log.Warn().Msg("Delete: Expected Error bot got no error")
		n.Errors++
		return
	}

	if err == nil && !expectError && status.Status == client.EnumStatus_ERROR {
		log.Warn().Msgf("Delete: error %s with status %s", err, status.Status)
		n.Errors++
	}
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

// AddUser calls the AddUser endpoint of the GRPC interface.
func (n *Node) AddUser(user string, kgname string, role string, expectError bool) {
	status, err := n.Client.AddUser(context.Background(), &client.UserRequest{
		User:     user,
		Keygroup: kgname,
		Role:     strToRole(role),
	})

	if err != nil && !expectError {
		log.Warn().Msgf("AddUser: error %s", err)
		n.Errors++
		return
	}

	if err == nil && expectError {
		log.Warn().Msg("AddUser: Expected Error bot got no error")
		n.Errors++
		return
	}

	if err == nil && !expectError && status.Status == client.EnumStatus_ERROR {
		log.Warn().Msgf("AddUser: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// RemoveUser calls the RemoveUser endpoint of the GRPC interface.
func (n *Node) RemoveUser(user string, kgname string, role string, expectError bool) {
	status, err := n.Client.RemoveUser(context.Background(), &client.UserRequest{
		User:     user,
		Keygroup: kgname,
		Role:     strToRole(role),
	})

	if err != nil && !expectError {
		log.Warn().Msgf("RemoveUser: error %s", err)
		n.Errors++
		return
	}

	if err == nil && expectError {
		log.Warn().Msg("RemoveUser: Expected Error bot got no error")
		n.Errors++
		return
	}

	if err == nil && !expectError && status.Status == client.EnumStatus_ERROR {
		log.Warn().Msgf("RemoveUser: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// type FredNode interface {
//	CreateKeygroup(kgname string, expectError bool)
//	DeleteKeygroup(kgname string, expectError bool)
//	GetKeygroupReplica(kgname string, expectError bool)
//	AddKeygroupReplica(kgname, replicaNodeID string, expectError bool)
//	DeleteKeygroupReplica(kgname, replicaNodeID string, expectError bool)
//	GetKeygroupTriggers(kgname string, expectError bool)
//	AddKeygroupTrigger(kgname, triggerNodeID, triggerNodeHost string, expectError bool)
//	DeleteKeygroupTrigger(kgname, triggerNodeID string, expectError bool)
//	GetAllReplica(expectError bool) []string
//	GetReplica(nodeID string, expectError bool) string
//	DeleteReplica(nodeID string, expectError bool)
//	PutItem(kgname, item string, data string, expectError bool)
//	GetItem(kgname, item string, expectError bool) string
//	DeleteItem(kgname, item string, expectError bool)
//  AddUser(user string, kgname string, role string, expectError bool)
//  RemoveUser(user string, kgname string, role string, expectError bool)
//}
