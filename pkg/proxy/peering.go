package proxy

import (
	"context"
	"fmt"
	"io"
	"time"

	"git.tu-berlin.de/mcc-fred/fred/pkg/grpcutil"
	"git.tu-berlin.de/mcc-fred/fred/proto/peering"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type PeeringProxy struct {
	p    *Proxy
	port int
	conn map[string]peering.NodeClient
	opts grpc.DialOption
}

func StartPeeringProxy(p *Proxy, port int, certFile string, keyFile string, caFile string, skipVerify bool) (*grpc.Server, error) {

	creds, _, err := grpcutil.GetCreds(certFile, keyFile, []string{caFile}, false, skipVerify)

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
	start := time.Now()

	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("PeeringProxy: PutItem: %s %s starting after %s", req.Keygroup, req.Id, time.Since(start))

	resp, err := c.PutItem(ctx, req)

	log.Debug().Msgf("PeeringProxy: PutItem: %s %s: %s", req.Keygroup, req.Id, time.Since(start))

	return resp, err
}

func (p *PeeringProxy) StreamPut(server peering.Node_StreamPutServer) error {
	req, err := server.Recv()

	if err != nil {
		return err
	}

	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return err
	}

	s, err := c.StreamPut(server.Context())

	if err != nil {
		return err
	}

	err = s.Send(req)

	if err != nil {
		return err
	}

	for item, err := server.Recv(); err != io.EOF; item, err = server.Recv() {
		if err != nil {
			return err
		}

		err = s.Send(item)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetItem has no implementation
func (p *PeeringProxy) GetItem(ctx context.Context, req *peering.GetItemRequest) (*peering.ItemResponse, error) {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	return c.GetItem(ctx, req)
}

// GetAllItems has no implementation
func (p *PeeringProxy) GetAllItems(req *peering.GetAllItemsRequest, server peering.Node_GetAllItemsServer) error {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return err
	}

	s, err := c.GetAllItems(server.Context(), req)

	if err != nil {
		return err
	}

	for item, err := s.Recv(); err != io.EOF; item, err = s.Recv() {
		if err != nil {
			return err
		}
		err = server.Send(item)

		if err != nil {
			return err
		}
	}

	return nil
}
