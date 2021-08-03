package proxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"git.tu-berlin.de/mcc-fred/fred/proto/peering"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type PeeringProxy struct {
	p    *Proxy
	port int
	conn map[string]peering.NodeClient
	opts grpc.DialOption
}

func StartPeeringProxy(p *Proxy, port int, certFile string, keyFile string, caFile string) (*grpc.Server, error) {
	if certFile == "" {
		log.Fatal().Msg("peering proxy: no certificate file given")
	}

	if keyFile == "" {
		log.Fatal().Msg("peering proxy: no key file given")
	}

	if caFile == "" {
		log.Fatal().Msg("peering proxy: no root certificate file given")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		return nil, err
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := ioutil.ReadFile(caFile)

	if err != nil {
		return nil, err
	}

	rootCAs.AppendCertsFromPEM(loaded)
	// Create the credentials and return it

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootCAs,
		RootCAs:      rootCAs,
		MinVersion:   tls.VersionTLS12,
	}

	a := &PeeringProxy{
		p:    p,
		port: port,
		conn: make(map[string]peering.NodeClient),
		opts: grpc.WithTransportCredentials(credentials.NewTLS(config)),
	}

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(config)))

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

// AddReplica calls this Method on the Inthandler
func (p *PeeringProxy) AddReplica(ctx context.Context, req *peering.AddReplicaRequest) (*peering.Empty, error) {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	return c.AddReplica(ctx, req)
}

// RemoveReplica calls this Method on the Inthandler
func (p *PeeringProxy) RemoveReplica(ctx context.Context, req *peering.RemoveReplicaRequest) (*peering.Empty, error) {
	c, err := p.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	return c.RemoveReplica(ctx, req)
}
