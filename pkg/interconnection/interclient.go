package interconnection

import (
	"context"
	"errors"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// InterClient communicates another nodes inthandler
type InterClient struct {
}

// NewClient creates a new Interclient
func NewClient() *InterClient {
	return &InterClient{}
}

// createConnAndClient creates a new connection to a server.
// Maybe it could be useful to reuse these?
// IDK whether this would be faster to store them in a map
func (i InterClient) getConnAndClient(host fred.Address, port int) (client NodeClient, conn *grpc.ClientConn) {
	conn, err := grpc.Dial(net.JoinHostPort(host.Addr, string(port)), grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return nil, nil
	}
	log.Info().Msgf("Interclient: Created Connection to %s:%d", host, port)
	client = NewNodeClient(conn)
	return
}

// logs the response and returns the correct error message
func (i InterClient) dealWithStatusResponse(res *StatusResponse, err error, from string) error {
	log.Debug().Msgf("Interclient got Response from %s, Status %s with Message %s and Error %s", from, res.Status, res.ErrorMessage, err)
	if err != nil {
		return err
	} else if res.Status == EnumStatus_ERROR {
		return errors.New(res.ErrorMessage)
	} else {
		return nil
	}
}

// Destroy currently does nothing, but might delete open connections if they are implemented
func (i InterClient) Destroy(){
}

// SendCreateKeygroup sends this command to the server at this address
func (i InterClient) SendCreateKeygroup(addr fred.Address, port int, kgname fred.KeygroupName) error {
	client, conn := i.getConnAndClient(addr, port)
	res, err := client.CreateKeygroup(context.Background(), &CreateKeygroupRequest{Keygroup: string(kgname)})
	conn.Close()
	return i.dealWithStatusResponse(res, err, "CreateKeygroup")
}

// SendDeleteKeygroup sends this command to the server at this address
func (i InterClient) SendDeleteKeygroup(addr fred.Address, port int, kgname fred.KeygroupName) error {
	client, conn := i.getConnAndClient(addr, port)
	res, err := client.DeleteKeygroup(context.Background(), &DeleteKeygroupRequest{Keygroup: string(kgname)})
	conn.Close()
	return i.dealWithStatusResponse(res, err, "DeleteKeygroup")
}

// SendUpdate sends this command to the server at this address
func (i InterClient) SendUpdate(addr fred.Address, port int, kgname fred.KeygroupName, id string, value string) error {
	client, conn := i.getConnAndClient(addr, port)
	res, err := client.PutItem(context.Background(), &PutItemRequest{
		Keygroup: string(kgname),
		Id:       id,
		Data:     value,
	})
	conn.Close()
	return i.dealWithStatusResponse(res, err, "SendUpdate")
}

// SendDelete sends this command to the server at this address
func (i InterClient) SendDelete(addr fred.Address, port int, kgname fred.KeygroupName, id string) error {
	client, _ := i.getConnAndClient(addr, port)
	res, err := client.DeleteItem(context.Background(), &DeleteItemRequest{
		Keygroup: string(kgname),
		Id: id,
	})
	return i.dealWithStatusResponse(res, err, "DeleteItem")
}

// SendAddReplica sends this command to the server at this address
func (i InterClient) SendAddReplica(addr fred.Address, port int, kgname fred.KeygroupName, node fred.Node) error {
	client, _ := i.getConnAndClient(addr, port)
	res, err := client.AddReplica(context.Background(), &AddReplicaRequest{
		NodeId:   string(node.ID),
		Keygroup: string(kgname),
	})
	return i.dealWithStatusResponse(res, err, "AddReplica")
}

// SendRemoveReplica sends this command to the server at this address
func (i InterClient) SendRemoveReplica(addr fred.Address, port int, kgname fred.KeygroupName, node fred.Node) error {
	client, _ := i.getConnAndClient(addr, port)
	res, err := client.RemoveReplica(context.Background(), &RemoveReplicaRequest{
		NodeId:   string(node.ID),
		Keygroup: string(kgname),
	})
	return i.dealWithStatusResponse(res, err, "RemoveReplica")
}
