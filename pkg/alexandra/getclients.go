package alexandra

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

// GetClientTo returns a client with this address
func (m *ClientsMgr) GetClientTo(host string) (client *Client) {
	log.Info().Msgf("GetClientTo: Trying to get Fred Client to host %s", host)
	client = m.clients[host]
	if client != nil {
		return
	}
	client = NewClient(host, m.clientsCert, m.clientsKey)
	m.clients[host] = client
	return
}

func getFastestClient(clts map[string]*Client) (client *Client) {
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

func getFastestClientByClientExpiry(clts map[string]*clientExpiry) (client *Client) {
	if clts == nil {
		return nil
	}
	var clientsMap = make(map[string]*Client)
	for key, value := range clts {
		clientsMap[key] = value.client
	}
	return getFastestClient(clientsMap)
}

// filterClientsToExpiry if param=-1 then only exp==-1; if param=0 then anything; if param>=0 then anything exp >= than param
func filterClientsToExpiry(clientEx map[string]*clientExpiry, expiry int64) (out map[string]*clientExpiry) {
	if clientEx == nil {
		return nil
	}
	out = make(map[string]*clientExpiry)
	for k, v := range clientEx {
		if expiry == -1 && v.expiry == -1 {
			out[k] = v
		} else if expiry == 0 {
			out[k] = v
		} else if v.expiry >= expiry {
			out[k] = v
		}
	}
	return
}

// GetFastestClient searches for the fastest of the already existing clients
func (m *ClientsMgr) GetFastestClient() (client *Client) {
	if len(m.clients) == 0 {
		log.Info().Msg("Fredclients: GetFastestClient was called but there are not clients. Using lighthouse client")
		return m.GetClientTo(m.lighthouse)
	}
	return getFastestClient(m.clients)
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
		go m.updateKeygroupClients(keygroup)
	}
}

// updateKeygroupClients updates the clients a keygroup has in a blocking way
func (m *ClientsMgr) updateKeygroupClients(keygroup string) {
	log.Debug().Msgf("Updating Clients for Keygroup %s", keygroup)
	replica, err := m.GetFastestClient().GetKeygroupReplica(context.Background(), keygroup)
	if err != nil {
		replica, err = m.GetClientTo(m.lighthouse).GetKeygroupReplica(context.Background(), keygroup)
		if err != nil {
			log.Error().Msgf("updateKeygroupClients cannot reach fastest client OR lighthouse...")
			return
		}
	}
	log.Debug().Msgf("updateKeygroupClients: Got replicas: %#v", replica)

	m.Lock()
	defer m.Unlock()
	set, exists := m.keygroups[keygroup]
	if !exists {
		m.keygroups[keygroup] = &keygroupSet{
			lastUpdated: time.Now(),
			clients:     make(map[string]*clientExpiry),
		}
		set = m.keygroups[keygroup]
	}
	set.clients = make(map[string]*clientExpiry)
	for _, client := range replica.Replica {
		set.clients[client.NodeId] = &clientExpiry{
			client: m.GetClientTo(client.Host),
			expiry: client.Expiry,
		}
	}
	set.lastUpdated = time.Now()
	m.keygroups[keygroup] = set
	log.Debug().Msgf("updateKeygroupClients: new Clients are: %#v", set)
}

// GetFastestClientWithKeygroup returns the fastest client that has the keygroup with an expiry bigger than the parameter
// set expiry to 1 to get any client, 0=no expiry
func (m *ClientsMgr) GetFastestClientWithKeygroup(keygroup string, expiry int64) (client *Client, err error) {
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
	log.Debug().Msgf("Clients before filtering: %#v", clients)
	filteredClients := filterClientsToExpiry(clients.clients, expiry)
	fastestClient := getFastestClientByClientExpiry(filteredClients)
	if fastestClient == nil {
		return nil, fmt.Errorf("was not able to find any client to keygroup %s with expiry > %d", keygroup, expiry)
	}
	return fastestClient, nil
}

func (m *ClientsMgr) GetRandomClientWithKeygroup(keygroup string, expiry int64) (client *Client, err error) {
	m.maybeUpdateKeygroupClients(keygroup)
	m.Lock()
	clients := m.keygroups[keygroup]
	m.Unlock()
	filtered := filterClientsToExpiry(clients.clients, expiry)
	// Get random element from this list
	log.Debug().Msgf("Len filtered is %#v", len(filtered))
	if len(filtered) == 0 {
		return nil, fmt.Errorf("was not able to find ANY client to keygroup %s with expiry > %d. Clients: %#v", keygroup, expiry, clients)
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
