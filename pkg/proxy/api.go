package proxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"git.tu-berlin.de/mcc-fred/fred/pkg/grpcutil"
	"git.tu-berlin.de/mcc-fred/fred/proto/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type APIProxy struct {
	p     *Proxy
	port  int
	conn  map[string]client.ClientClient
	opts  grpc.DialOption
	roots *x509.CertPool
}

func StartAPIProxy(p *Proxy, port int, cert string, key string, caCert string, skipVerify bool) (*grpc.Server, error) {
	creds, rootCA, err := grpcutil.GetCredsFromConfig(cert, key, []string{caCert}, false, skipVerify, &tls.Config{ClientAuth: tls.RequireAndVerifyClientCert})

	if err != nil {
		return nil, err
	}

	a := &APIProxy{
		p:     p,
		port:  port,
		conn:  make(map[string]client.ClientClient),
		opts:  grpc.WithTransportCredentials(creds),
		roots: rootCA,
	}

	s := grpc.NewServer(grpc.Creds(creds))

	client.RegisterClientServer(s, a)

	return s, nil
}

func (a *APIProxy) getConn(keygroup string) (client.ClientClient, error) {
	host := a.p.getHost(keygroup)

	if c, ok := a.conn[host]; ok {
		return c, nil
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, a.port), a.opts)

	if err != nil {
		return nil, err
	}

	c := client.NewClientClient(conn)

	a.conn[host] = c
	return c, nil
}

func (a *APIProxy) getAny() (client.ClientClient, error) {
	host := a.p.getAny()

	if c, ok := a.conn[host]; ok {
		return c, nil
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, a.port), a.opts)

	if err != nil {
		return nil, err
	}

	c := client.NewClientClient(conn)

	a.conn[host] = c
	return c, nil
}

func (a *APIProxy) addUserHeader(ctx context.Context) (context.Context, error) {
	// since we pass each request on to the backend with our proxy, we lose user information
	// we add that back on by adding the "user" header to the context

	// get peer information
	p, ok := peer.FromContext(ctx)

	if !ok {
		return ctx, fmt.Errorf("no peer found")
	}

	// get TLS credential information for this connection
	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)

	if !ok {
		return ctx, fmt.Errorf("unexpected peer transport credentials")
	}

	// check that the certificate exists
	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return ctx, fmt.Errorf("could not verify peer certificate: %v", tlsAuth.State)
	}

	host, _, err := net.SplitHostPort(p.Addr.String())

	if err != nil {
		return ctx, err
	}

	// verify the certificate:
	// IF we are not proxied and communicate with the client directly:
	// 1) it should be issued by a CA in our root CA pool
	// 2) any intermediates are valid for us
	// 3) the certificate should be valid for client authentication
	// 4) the certificate should have the clients address as a SAN
	_, err = tlsAuth.State.VerifiedChains[0][0].Verify(x509.VerifyOptions{
		Roots:         a.roots,
		CurrentTime:   time.Now(),
		Intermediates: x509.NewCertPool(),
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		DNSName:       host,
	})

	if err != nil {
		return ctx, err
	}

	// Check subject common name exists and return it as the user name for that client
	name := tlsAuth.State.VerifiedChains[0][0].Subject.CommonName

	if name == "" {
		return ctx, fmt.Errorf("invalid subject common name")
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return ctx, fmt.Errorf("no metadata could be found for proxied request")
	}

	newMD := md.Copy()
	newMD.Set("user", name)

	return metadata.NewOutgoingContext(ctx, newMD), nil

}

// CreateKeygroup calls this method on the exthandler
func (a *APIProxy) CreateKeygroup(ctx context.Context, req *client.CreateKeygroupRequest) (*client.Empty, error) {

	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.CreateKeygroup(ctx, req)
}

// DeleteKeygroup calls this method on the exthandler
func (a *APIProxy) DeleteKeygroup(ctx context.Context, req *client.DeleteKeygroupRequest) (*client.Empty, error) {

	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.DeleteKeygroup(ctx, req)
}

// Read calls this method on the exthandler
func (a *APIProxy) Read(ctx context.Context, req *client.ReadRequest) (*client.ReadResponse, error) {

	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.Read(ctx, req)
}

// Scan calls this method on the exthandler
func (a *APIProxy) Scan(ctx context.Context, req *client.ScanRequest) (*client.ScanResponse, error) {

	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.Scan(ctx, req)
}

// Keys calls this method on the exthandler
func (a *APIProxy) Keys(ctx context.Context, req *client.KeysRequest) (*client.KeysResponse, error) {

	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.Keys(ctx, req)
}

// Append calls this method on the exthandler
func (a *APIProxy) Append(ctx context.Context, req *client.AppendRequest) (*client.AppendResponse, error) {

	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.Append(ctx, req)
}

// Update calls this method on the exthandler
func (a *APIProxy) Update(ctx context.Context, req *client.UpdateRequest) (*client.UpdateResponse, error) {

	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.Update(ctx, req)
}

// Delete calls this method on the exthandler
func (a *APIProxy) Delete(ctx context.Context, req *client.DeleteRequest) (*client.DeleteResponse, error) {

	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.Delete(ctx, req)
}

// AddReplica calls this method on the exthandler
func (a *APIProxy) AddReplica(ctx context.Context, req *client.AddReplicaRequest) (*client.Empty, error) {

	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.AddReplica(ctx, req)
}

// GetKeygroupInfo calls this method on the exthandler
func (a *APIProxy) GetKeygroupInfo(ctx context.Context, req *client.GetKeygroupInfoRequest) (*client.GetKeygroupInfoResponse, error) {
	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.GetKeygroupInfo(ctx, req)
}

// RemoveReplica calls this method on the exthandler
func (a *APIProxy) RemoveReplica(ctx context.Context, req *client.RemoveReplicaRequest) (*client.Empty, error) {
	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.RemoveReplica(ctx, req)
}

// GetReplica calls this method on the exthandler
func (a *APIProxy) GetReplica(ctx context.Context, req *client.GetReplicaRequest) (*client.GetReplicaResponse, error) {
	c, err := a.getAny()

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.GetReplica(ctx, req)
}

// GetAllReplica calls this method on the exthandler
func (a *APIProxy) GetAllReplica(ctx context.Context, req *client.Empty) (*client.GetAllReplicaResponse, error) {
	c, err := a.getAny()

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.GetAllReplica(ctx, req)
}

// GetKeygroupTriggers calls this method on the exthandler
func (a *APIProxy) GetKeygroupTriggers(ctx context.Context, req *client.GetKeygroupTriggerRequest) (*client.GetKeygroupTriggerResponse, error) {
	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.GetKeygroupTriggers(ctx, req)
}

// AddTrigger calls this method on the exthandler
func (a *APIProxy) AddTrigger(ctx context.Context, req *client.AddTriggerRequest) (*client.Empty, error) {
	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.AddTrigger(ctx, req)
}

// RemoveTrigger calls this method on the exthandler
func (a *APIProxy) RemoveTrigger(ctx context.Context, req *client.RemoveTriggerRequest) (*client.Empty, error) {
	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.RemoveTrigger(ctx, req)
}

// AddUser calls this method on the exthandler
func (a *APIProxy) AddUser(ctx context.Context, req *client.AddUserRequest) (*client.Empty, error) {
	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.AddUser(ctx, req)
}

// RemoveUser calls this method on the exthandler
func (a *APIProxy) RemoveUser(ctx context.Context, req *client.RemoveUserRequest) (*client.Empty, error) {
	c, err := a.getConn(req.Keygroup)

	if err != nil {
		return nil, err
	}

	ctx, err = a.addUserHeader(ctx)
	if err != nil {
		return nil, err
	}

	return c.RemoveUser(ctx, req)
}
