package alexandra

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	clientsProto "git.tu-berlin.de/mcc-fred/fred/proto/client"
	"git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"git.tu-berlin.de/mcc-fred/vclock"
	"github.com/rs/zerolog/log"
)

const (
	// When to update stored information about clients
	keygroupTimeout = 2 * time.Minute
	// How many other nodes to ask when ReadFromAnywhere is called
	otherNodesToAsk = 3
	// UseSlowerNodeProb In how many percent of the cases: instead of using the fastest client, use a random client to update its readSpeed
	// only used for Read,Write,Delete,Append (since these are the only operations that update the readSpeed)
	UseSlowerNodeProb = 0.089
)

type clientExpiry struct {
	client *Client
	expiry int64
}

// keygroupSet represents a keygroups clients and the last time this information was updated from FreD
type keygroupSet struct {
	lastUpdated time.Time
	preferred   *Client
	clients     []*clientExpiry
}

// ClientsMgr manages all Clients to Fred that Alexandra has. Is used to get fastest clients to keygroups etc. and to read from anywhere
// there are 3 variables important for this configuration: keygroupTimeout, otherNodesToAsk, getSlowerNodeProb. Please see their documentation.
type ClientsMgr struct {
	// Mutex for the keygroups map, because it might be changed while iterated over
	sync.Mutex
	clients                                     map[string]*Client
	clientsCert, clientsKey, caCert, lighthouse string
	keygroups                                   map[string]*keygroupSet
	experimental                                bool
}

func newClientsManager(clientsCert string, clientsKey string, caCert string, lighthouse string, experimental bool) *ClientsMgr {
	mgr := &ClientsMgr{
		clients:      make(map[string]*Client),
		clientsCert:  clientsCert,
		clientsKey:   clientsKey,
		caCert:       caCert,
		lighthouse:   lighthouse,
		keygroups:    make(map[string]*keygroupSet),
		experimental: experimental,
	}
	// add the lighthouse client to the clients list
	mgr.getLightHouse()
	// rand.Seed(time.Now().UnixNano())
	return mgr
}

func (m *ClientsMgr) readFromAnywhere(request *middleware.ReadRequest) ([]string, []vclock.VClock, error) {
	log.Debug().Msgf("ClientsManager is reading from anywhere. Req= %+v", request)
	m.maybeUpdateKeygroupClients(request.Keygroup)

	// Start a coroutine to every fastestClient we ask
	m.Lock()
	set, exists := m.keygroups[request.Keygroup]
	m.Unlock()

	type readResponse struct {
		vals     []string
		versions []vclock.VClock
	}

	// if there is a preferred client, only try to ask that
	if set.preferred != nil {
		log.Debug().Msgf("there is a preferred node, reading from that: %s", set.preferred.nodeID)
		res, err := set.preferred.Client.Read(context.Background(), &clientsProto.ReadRequest{Id: request.Id, Keygroup: request.Keygroup})

		if err != nil {
			log.Err(err).Msgf("Reading from preferred client %s returned error", set.preferred.nodeID)
		} else {
			vals := make([]string, len(res.Data))
			versions := make([]vclock.VClock, len(res.Data))

			for i := range res.Data {
				vals[i] = res.Data[i].Val
				versions[i] = res.Data[i].Version.Version
				log.Debug().Msgf("Reading from client %s returned data: %+v %+v", set.preferred.nodeID, res.Data[i].Val, res.Data[i].Version.Version)
			}

			return vals, versions, nil
		}
	}

	clientsToAsk := make(map[*Client]struct{})

	clts := filterClientsToExpiry(set.clients, request.MinExpiry)

	if !exists || len(clts) == 0 {
		log.Error().Msgf("...found no clients with minimum expiry. Clients with longer expiry: %+v", set.clients)
		return nil, nil, errors.New("there are no members of this keygroup")
	}

	// let's figure out who we want to ask
	// Ask the fastest fastestClient
	fastestClient, err := m.getFastestClientWithKeygroup(request.Keygroup, request.MinExpiry)
	if err == nil {
		clientsToAsk[fastestClient] = struct{}{}
	}

	// and add a maximum of otherNodeToAsk other clients
	for i := 0; i < otherNodesToAsk; i++ {
		clientsToAsk[clts[rand.Intn(len(clts))].client] = struct{}{}
	}

	log.Debug().Msgf("asking %d nodes", len(clientsToAsk))

	var wg sync.WaitGroup
	responses := make(chan readResponse, len(clientsToAsk))
	done := make(chan struct{})

	for c := range clientsToAsk {
		wg.Add(1)

		go func(c *Client) {
			defer wg.Done()

			log.Debug().Msgf("...asking Client %s for Keygroup %s", c.nodeID, request.Keygroup)
			res, err := c.Client.Read(context.Background(), &clientsProto.ReadRequest{Id: request.Id, Keygroup: request.Keygroup})

			if err != nil {
				log.Err(err).Msgf("Reading from client %s returned error", c.nodeID)
				return
			}

			r := readResponse{
				vals:     make([]string, len(res.Data)),
				versions: make([]vclock.VClock, len(res.Data)),
			}
			for i := range res.Data {
				r.vals[i] = res.Data[i].Val
				r.versions[i] = res.Data[i].Version.Version
				log.Debug().Msgf("Reading from client %s returned data: %+v %+v", c.nodeID, res.Data[i].Val, res.Data[i].Version.Version)
			}

			responses <- r

		}(c)
	}

	// wait for all responses to come in and close the channel
	go func() {
		wg.Wait()
		done <- struct{}{}
		close(responses)
	}()

	// if we get a response, return that
	// otherwise, if done is called
	select {
	case r := <-responses:
		return r.vals, r.versions, nil
	case <-done:
		break
	}

	// There was no successful response -- Update the keygroup information and try one last time
	log.Info().Msg("ReadFromAnywhere: Was not able to reach any queried node, updating cache and retrying...")
	m.updateKeygroupClients(request.Keygroup)

	client, err := m.getFastestClientWithKeygroup(request.Keygroup, request.MinExpiry)

	if err != nil {
		return nil, nil, fmt.Errorf("ReadFromAnywhere: there is no client with keygroup %s and expiry %d", request.Keygroup, request.MinExpiry)
	}

	result, err := client.read(context.Background(), request.Keygroup, request.Id)
	if err != nil {
		return nil, nil, fmt.Errorf("ReadFromAnywhere: cannot read from fastest client. err=%v", err)
	}

	vals := make([]string, len(result.Data))
	versions := make([]vclock.VClock, len(result.Data))

	for i := range result.Data {
		vals[i] = result.Data[i].Val
		versions[i] = result.Data[i].Version.Version
	}

	return vals, versions, nil
}

func (m *ClientsMgr) getLightHouse() (client *Client) {
	return m.getClientTo(m.lighthouse, "__lighthouse")
}

// GetClientTo returns a client with this address
func (m *ClientsMgr) getClientTo(host string, nodeID string) (client *Client) {
	log.Info().Msgf("GetClientTo: Trying to get Fred Client to node %s host %s", nodeID, host)
	client = m.clients[nodeID]

	if client != nil {
		return
	}

	client = newClient(nodeID, host, m.clientsCert, m.clientsKey, m.caCert)
	m.clients[nodeID] = client
	return
}

func getFastestClient(clts []*Client) (client *Client) {
	var minTime float32 = math.MaxFloat32
	var minClient *Client
	// Set the first client up as fastest client, so that it gets returned if no other client is found.
	for _, value := range clts {
		minClient = value
		break
	}
	for _, value := range clts {
		if value.ReadSpeed < minTime && value.ReadSpeed != -1 {
			minTime = value.ReadSpeed
			minClient = value
		}
	}
	return minClient
}

func getFastestClientByClientExpiry(clts []*clientExpiry) (client *Client) {
	if clts == nil {
		return nil
	}
	var clients = make([]*Client, len(clts))
	for i, c := range clts {
		clients[i] = c.client
	}
	return getFastestClient(clients)
}

// filterClientsToExpiry if param=-1 then only exp==-1; if param=0 then anything; if param>=0 then anything exp >= than param
func filterClientsToExpiry(clientEx []*clientExpiry, expiry int64) (out []*clientExpiry) {
	if clientEx == nil {
		return nil
	}
	out = make([]*clientExpiry, 0)
	for _, v := range clientEx {
		if expiry == -1 && v.expiry == -1 {
			out = append(out, v)
		} else if expiry == 0 {
			out = append(out, v)
		} else if v.expiry >= expiry {
			out = append(out, v)
		}
	}
	return
}

func (m *ClientsMgr) setPreferred(keygroup string, nodeID string) error {
	c, ok := m.clients[nodeID]

	if !ok {
		return fmt.Errorf("unknown node %s", nodeID)
	}

	m.Lock()
	defer m.Unlock()

	if _, ok := m.keygroups[keygroup]; !ok {
		// turns out we don't know about this keygroup yet
		m.Unlock()
		m.updateKeygroupClients(keygroup)
		m.Lock()
	}

	m.keygroups[keygroup].preferred = c

	return nil
}

func (m *ClientsMgr) getClient(keygroup string) (*Client, error) {
	m.Lock()
	// if there is a preferred client for that keygroup, use that
	if k, ok := m.keygroups[keygroup]; ok && k.preferred != nil {
		m.Unlock()
		return k.preferred, nil
	}
	m.Unlock()

	if m.experimental && rand.Float64() < UseSlowerNodeProb {
		return m.getRandomClientWithKeygroup(keygroup, 0)
	}

	return m.getFastestClientWithKeygroup(keygroup, 0)
}

// getFastestClient searches for the fastest of the already existing clients
func (m *ClientsMgr) getFastestClient() (client *Client) {
	if len(m.clients) == 0 {
		log.Info().Msg("ClientsMgr: GetFastestClient was called but there are not clients. Using lighthouse client")
		return m.getLightHouse()
	}

	clts := make([]*Client, 0, len(m.clients))

	for _, c := range m.clients {
		clts = append(clts, c)
	}

	return getFastestClient(clts)
}

// GetFastestClientWithKeygroup returns the fastest client that has the keygroup with an expiry bigger than the parameter
// set expiry to 1 to get any client, 0=no expiry
func (m *ClientsMgr) getFastestClientWithKeygroup(keygroup string, expiry int64) (client *Client, err error) {
	m.maybeUpdateKeygroupClients(keygroup)
	m.Lock()
	clients := m.keygroups[keygroup]
	m.Unlock()
	if clients == nil {
		log.Debug().Msgf("GetFastestClientWithKeygroup: No clients to keygroup %s", keygroup)
		m.updateKeygroupClients(keygroup)
		m.Lock()
		clients = m.keygroups[keygroup]
		m.Unlock()
	}
	log.Debug().Msgf("Clients before filtering: %+v", clients)
	filteredClients := filterClientsToExpiry(clients.clients, expiry)
	fastestClient := getFastestClientByClientExpiry(filteredClients)
	if fastestClient == nil {
		return nil, fmt.Errorf("was not able to find any client to keygroup %s with expiry > %d", keygroup, expiry)
	}
	return fastestClient, nil
}

func (m *ClientsMgr) getRandomClientWithKeygroup(keygroup string, expiry int64) (client *Client, err error) {
	m.maybeUpdateKeygroupClients(keygroup)
	m.Lock()
	clients := m.keygroups[keygroup]
	m.Unlock()
	filtered := filterClientsToExpiry(clients.clients, expiry)
	// Get random element from this list
	log.Debug().Msgf("Len filtered is %+v", len(filtered))
	if len(filtered) == 0 {
		return nil, fmt.Errorf("was not able to find ANY client to keygroup %s with expiry > %d. Clients: %+v", keygroup, expiry, clients)
	}
	nodeI := rand.Intn(len(filtered))
	curI := 0
	for _, v := range filtered {
		if nodeI == curI {
			return v.client, nil
		}
		curI++
	}
	return nil, fmt.Errorf("was not able to find RANDOM client to keygroup %s with expiry > %d", keygroup, expiry)
}

// maybeUpdateKeygroupClients updates the cached keygroups of a client if it hasn't happened $keygroupCacheTimeout
func (m *ClientsMgr) maybeUpdateKeygroupClients(keygroup string) {
	m.Lock()
	set, exists := m.keygroups[keygroup]
	m.Unlock()
	if !exists {
		log.Debug().Msgf("Keygroup %s has no entries, updating them now", keygroup)
		m.updateKeygroupClients(keygroup)
	} else if time.Since(set.lastUpdated) > keygroupTimeout {
		log.Debug().Msgf("Keygroup %s has not been updated in %.0f minutes, doing it now", keygroup, keygroupTimeout.Minutes())
		m.updateKeygroupClients(keygroup)
	}
}

// updateKeygroupClients updates the clients a keygroup has in a blocking way
func (m *ClientsMgr) updateKeygroupClients(keygroup string) {
	log.Debug().Msgf("Updating Clients for Keygroup %s", keygroup)
	replica, err := m.getFastestClient().getKeygroupReplica(context.Background(), keygroup)
	if err != nil {
		log.Debug().Msgf("couldn't get replicas for keygroup %s from fastest client: %s", keygroup, err.Error())
		replica, err = m.getLightHouse().getKeygroupReplica(context.Background(), keygroup)
		if err != nil {
			log.Debug().Msgf("couldn't get replicas for keygroup %s from lighthouse client: %s", keygroup, err.Error())
			log.Error().Msgf("updateKeygroupClients cannot reach fastest client OR lighthouse...")
			return
		}
	}

	log.Debug().Msgf("updateKeygroupClients: Got replicas: %+v", replica)

	m.Lock()
	defer m.Unlock()
	set, exists := m.keygroups[keygroup]
	if !exists {
		m.keygroups[keygroup] = &keygroupSet{
			lastUpdated: time.Now(),
		}
		set = m.keygroups[keygroup]
	}

	// we also need to find out if the preferred node for that keygroup (if any) still exists
	preferred := m.keygroups[keygroup].preferred
	removePreferred := preferred == nil

	set.clients = make([]*clientExpiry, len(replica.Replica))
	for i, client := range replica.Replica {
		set.clients[i] = &clientExpiry{
			client: m.getClientTo(client.Host, client.NodeId),
			expiry: client.Expiry,
		}

		if !removePreferred && client.NodeId == preferred.nodeID {
			removePreferred = false
		}
	}

	if removePreferred {
		m.keygroups[keygroup].preferred = nil
	}

	set.lastUpdated = time.Now()
	m.keygroups[keygroup] = set
	log.Debug().Msgf("updateKeygroupClients: new Clients are: %+v", set)
}
