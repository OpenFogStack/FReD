package proxy

import (
	"context"
	"fmt"

	"git.tu-berlin.de/mcc-fred/fred/pkg/grpcutil"
	"git.tu-berlin.de/mcc-fred/fred/proto/peering"
	"google.golang.org/grpc"
)

type PeeringProxy struct {
	p    *Proxy
	port int
	conn map[string]peering.NodeClient
	opts grpc.DialOption
}

func StartPeeringProxy(p *Proxy, port int, certFile string, keyFile string, caFile string) (*grpc.Server, error) {

	creds, _, err := grpcutil.GetCreds(certFile, keyFile, []string{caFile}, false)

	if err != nil {
		return nil, err
	}

	a := &PeeringProxy{
		p:    p,
		port: port,
		conn: make(map[string]peering.NodeClient),
		opts: grpc.WithTransportCredentials(creds),
	}

	s := grpc.NewServer(grpc.Creds(creds))

	peering.RegisterNodeServer(s, a)

	return s, nil
}

func (p *PeeringProxy) getConn(keygroup string) (peering.NodeClient, error) {
	host := p.p.getHost(keygroup)

	if c, ok := p.conn[host]; ok {
		return c, nil
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, p.port), p.opts)

	if err != nil {
		return nil, err
	}

	c := peering.NewNodeClient(conn)

	p.conn[host] = c
	return c, nil
}

// CreateKeygroup calls this Method on the Inthandler
func (p *PeeringProxy) CreateKeygroup(ctx context.Context, req *peering.CreateKeygroupRequest) (*peering.Empty, error) {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	return c.CreateKeygroup(ctx, req)
}

// DeleteKeygroup calls this Method on the Inthandler
func (p *PeeringProxy) DeleteKeygroup(ctx context.Context, req *peering.DeleteKeygroupRequest) (*peering.Empty, error) {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	return c.DeleteKeygroup(ctx, req)
}

// PutItem calls HandleUpdate on the Inthandler
func (p *PeeringProxy) PutItem(ctx context.Context, req *peering.PutItemRequest) (*peering.Empty, error) {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	return c.PutItem(ctx, req)
}

// AppendItem calls HandleAppend on the Inthandler
func (p *PeeringProxy) AppendItem(ctx context.Context, req *peering.AppendItemRequest) (*peering.Empty, error) {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	return c.AppendItem(ctx, req)
}

// GetItem has no implementation
func (p *PeeringProxy) GetItem(ctx context.Context, req *peering.GetItemRequest) (*peering.GetItemResponse, error) {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	return c.GetItem(ctx, req)
}

// GetAllItems has no implementation
func (p *PeeringProxy) GetAllItems(ctx context.Context, req *peering.GetAllItemsRequest) (*peering.GetAllItemsResponse, error) {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	return c.GetAllItems(ctx, req)
}
