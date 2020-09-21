package grpcclient

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/proto/client"
	"google.golang.org/grpc"
)

// Node represents the API to a single FReD Node
type Node struct {
	Errors int
	Client client.ClientClient
	conn   *grpc.ClientConn
	Addr   string
}

// NewNode creates a new Node that represents a connection to a single fred instance
func NewNode(addr string, port int) *Node {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", addr, port), grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return nil
	}
	c := client.NewClientClient(conn)
	return &Node{
		Errors: 0,
		Client: c,
		conn:   conn,
		Addr:   addr,
	}
}

// Close closes the connection
func (n Node) Close() {
	err := n.conn.Close()

	log.Err(err).Msg("error closing node")
}

// CreateKeygroup calls the CreateKeygroup endpoint of the GRPC interface.
func (n Node) CreateKeygroup(kgname string, mutable bool, expiry int, expectError bool) {
	status, err := n.Client.CreateKeygroup(context.Background(), &client.CreateKeygroupRequest{Keygroup: kgname, Mutable: mutable, Expiry: int64(expiry)})

	if err != nil && !expectError {
		log.Warn().Msgf("CreateKeygroup: error %s", err)
		n.Errors++
		return
	}

	if err == nil && !expectError && status.Status == client.EnumStatus_ERROR {
		log.Warn().Msgf("CreateKeygroup: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// DeleteKeygroup calls the DeleteKeygroup endpoint of the GRPC interface.
func (n Node) DeleteKeygroup(kgname string, expectError bool) {
	status, err := n.Client.DeleteKeygroup(context.Background(), &client.DeleteKeygroupRequest{Keygroup: kgname})

	if (err != nil && !expectError) || (err == nil && !expectError && status.Status == client.EnumStatus_ERROR) {
		log.Warn().Msgf("DeleteKeygroup: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// GetKeygroupReplica calls the GetKeygroupReplica endpoint of the GRPC interface.
func (n Node) GetKeygroupReplica(kgname string, expectError bool) map[string]int {
	res, err := n.Client.GetKeygroupReplica(context.Background(), &client.GetKeygroupReplicaRequest{Keygroup: kgname})

	if err != nil && !expectError {
		log.Warn().Msgf("GetKeygroupReplica: error %s", err)
		n.Errors++
	}

	if res == nil {
		return nil
	}

	nodes := make(map[string]int)

	for i, node := range res.NodeId {
		nodes[node] = int(res.Expiry[i])
	}

	return nodes
}

// AddKeygroupReplica calls the AddKeygroupReplica endpoint of the GRPC interface.
func (n Node) AddKeygroupReplica(kgname, replicaNodeID string, expiry int, expectError bool) {
	status, err := n.Client.AddReplica(context.Background(), &client.AddReplicaRequest{
		Keygroup: kgname,
		NodeId:   replicaNodeID,
		Expiry:   int64(expiry),
	})
	if (err != nil && !expectError) || (err == nil && !expectError && status.Status == client.EnumStatus_ERROR) {
		log.Warn().Msgf("AddKeygroupReplica: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// DeleteKeygroupReplica calls the DeleteKeygroupReplica endpoint of the GRPC interface.
func (n Node) DeleteKeygroupReplica(kgname, replicaNodeID string, expectError bool) {
	status, err := n.Client.RemoveReplica(context.Background(), &client.RemoveReplicaRequest{
		Keygroup: kgname,
		NodeId:   replicaNodeID,
	})
	if (err != nil && !expectError) || (err == nil && !expectError && status.Status == client.EnumStatus_ERROR) {
		log.Warn().Msgf("DeleteKeygroupReplica: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// GetKeygroupTriggers calls the GetKeygroupTriggers endpoint of the GRPC interface.
func (n Node) GetKeygroupTriggers(kgname string, expectError bool) []*client.Trigger {
	res, err := n.Client.GetKeygroupTriggers(context.Background(), &client.GetKeygroupTriggerRequest{Keygroup: kgname})
	if err != nil && !expectError {
		log.Warn().Msgf("GetKeygroupTriggers: error %s", err)
		n.Errors++
	}

	if res == nil {
		return nil
	}

	return res.Triggers
}

// AddKeygroupTrigger calls the AddKeygroupTrigger endpoint of the GRPC interface.
func (n Node) AddKeygroupTrigger(kgname, triggerNodeID, triggerNodeHost string, expectError bool) {
	status, err := n.Client.AddTrigger(context.Background(), &client.AddTriggerRequest{Keygroup: kgname, TriggerId: triggerNodeID, TriggerHost: triggerNodeHost})
	if (err != nil && !expectError) || (err == nil && !expectError && status.Status == client.EnumStatus_ERROR) {
		log.Warn().Msgf("AddKeygroupTrigger: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// DeleteKeygroupTrigger calls the DeleteKeygroupTrigger endpoint of the GRPC interface.
func (n Node) DeleteKeygroupTrigger(kgname, triggerNodeID string, expectError bool) {
	status, err := n.Client.RemoveTrigger(context.Background(), &client.RemoveTriggerRequest{Keygroup: kgname, TriggerId: triggerNodeID})
	if (err != nil && !expectError) || (err == nil && !expectError && status.Status == client.EnumStatus_ERROR) {
		log.Warn().Msgf("DeleteKeygroupTrigger: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// GetAllReplica calls the GetAllReplica endpoint of the GRPC interface.
func (n Node) GetAllReplica(expectError bool) []string {
	res, err := n.Client.GetAllReplica(context.Background(), &client.GetAllReplicaRequest{})
	if err != nil && !expectError {
		log.Warn().Msgf("GetAllReplicas: error %s", err)
		n.Errors++
	}

	if res == nil {
		return nil
	}

	ids := make([]string, len(res.Replicas))

	for i := 0; i < len(res.Replicas); i++ {
		ids[i] = res.Replicas[i].NodeId
	}

	return ids
}

// GetReplica calls the GetReplica endpoint of the GRPC interface.
func (n Node) GetReplica(nodeID string, expectError bool) string {
	res, err := n.Client.GetReplica(context.Background(), &client.GetReplicaRequest{NodeId: nodeID})

	if err != nil && !expectError {
		log.Warn().Msgf("GetReplica: error %s", err)
		n.Errors++
	}

	if res == nil {
		return ""
	}

	return res.NodeId
}

// PutItem calls the PutItem endpoint of the GRPC interface.
func (n Node) PutItem(kgname, item string, data string, expectError bool) {
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

	if err == nil && status.Status == client.EnumStatus_ERROR && !expectError {
		log.Warn().Msgf("Update: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

// GetItem calls the GetItem endpoint of the GRPC interface.
func (n Node) GetItem(kgname, item string, expectError bool) string {
	res, err := n.Client.Read(context.Background(), &client.ReadRequest{
		Keygroup: kgname,
		Id:       item,
	})

	if err != nil && !expectError {
		log.Warn().Msgf("GetReplica: error %s", err)
		n.Errors++
	}

	if res == nil {
		return ""
	}

	return res.Data
}

// DeleteItem calls the DeleteItem endpoint of the GRPC interface.
func (n Node) DeleteItem(kgname, item string, expectError bool) {

	status, err := n.Client.Delete(context.Background(), &client.DeleteRequest{
		Keygroup: kgname,
		Id:       item,
	})

	if err != nil && !expectError {
		log.Warn().Msgf("Update: error %s", err)
		n.Errors++
		return
	}

	if err == nil && !expectError && status.Status == client.EnumStatus_ERROR {
		log.Warn().Msgf("Update: error %s with status %s", err, status.Status)
		n.Errors++
	}
}

//type FredNode interface {
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
//}
