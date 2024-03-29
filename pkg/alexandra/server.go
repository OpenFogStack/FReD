package alexandra

// Middleware listens to GRPC requests from clients (and sends them to the relevant Fred Node etc.)
// The implementation is split up into different files in this folder.
type Middleware struct {
	isProxied    bool
	proxyHost    string
	clientsMgr   *ClientsMgr
	lighthouse   string
	vcache        *vcache
	experimental bool
}

// NewMiddleware creates a new Middleware for requests from Alexandra Clients
func NewMiddleware(nodesCert string, nodesKey string, caCert string, lighthouse string, isProxied bool, proxyHost string, experimental bool, skipVerify bool) *Middleware {

	return &Middleware{
		isProxied:    isProxied,
		proxyHost:    proxyHost,
		clientsMgr:   newClientsManager(nodesCert, nodesKey, caCert, lighthouse, experimental, skipVerify),
		lighthouse:   lighthouse,
		vcache:        newVCache(),
		experimental: experimental,
	}
}
