package alexandra

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	clientsProto "git.tu-berlin.de/mcc-fred/fred/proto/client"
	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
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
	clients     map[string]*clientExpiry
}

// ClientsMgr manages all Clients to Fred that Alexandra has. Is used to get fastest clients to keygroups etc. and to read from anywhere
// there are 3 variables important for this configuration: keygruopTimeout, otherNodesToAsk, getSlowerNodeProb. Please see their documentation.
type ClientsMgr struct {
	// Mutex for the keygroups map, because it might be changed while iterated over
	sync.Mutex
	clients                             map[string]*Client
	clientsCert, clientsKey, lighthouse string
	keygroups                           map[string]*keygroupSet
}

func NewClientsManager(clientsCert, clientsKey, lighthouse string) *ClientsMgr {
	mgr := &ClientsMgr{
		clients:     make(map[string]*Client),
		clientsCert: clientsCert,
		clientsKey:  clientsKey,
		lighthouse:  lighthouse,
		keygroups:   make(map[string]*keygroupSet),
	}
	// add the lighthouse client to the clients list
	mgr.GetClientTo(lighthouse)
	rand.Seed(time.Now().UnixNano())
	return mgr
}

func (m *ClientsMgr) ReadFromAnywhere(ctx context.Context, request *alexandraProto.ReadRequest) (*alexandraProto.ReadResponse, error) {
	log.Debug().Msgf("ClientsManager is reading from anywhere. Req= %#v", request)
	go m.maybeUpdateKeygroupClients(request.Keygroup)
	type readResponse struct {
		error bool
		data  string
	}
	responses := make(chan readResponse)
	responsesClosed := false
	sentAsks := 0

	// Start a coroutine to every fastestClient we ask
	m.Lock()
	set, exists := m.keygroups[request.Keygroup]
	m.Unlock()
	clts := filterClientsToExpiry(set.clients, request.MinExpiry)
	if !exists || len(clts) == 0 {
		log.Error().Msgf("...found no clients whith minimum expiry. Clients with longer expiry: %#v", set.clients)
		return nil, errors.New("there are no members of this keygroup")
	}
	askNode := func(client *Client) {
		// TODO select channels?
		// TODO buffered Channel
		log.Debug().Msgf("...asking Client %#v for Keygroup %s", client, request.Keygroup)
		res, err := client.Client.Read(context.Background(), &clientsProto.ReadRequest{Id: request.Id, Keygroup: request.Keygroup})
		if err != nil {
			log.Err(err).Msg("Reading from client returned error")
			if !responsesClosed {
				responses <- readResponse{error: true, data: ""}
			}
		} else {
			log.Debug().Msgf("Reading from client returned data: %#v", res)
			if !responsesClosed {
				responses <- readResponse{error: false, data: res.Data}
			}
		}
	}

	// Ask the fastest fastestClient
	fastestClient, err := m.GetFastestClientWithKeygroup(request.Keygroup, request.MinExpiry)
	if err == nil {
		go askNode(fastestClient)
		sentAsks++
	}

	// Ask $otherNodesToAsk other Clients
	if len(clts) > 2 { // If its only one element long the one node is also the fastest node
		if otherNodesToAsk > len(clts) {
			for _, client := range clts {
				go askNode(client.client)
				sentAsks++
			}
		} else {
			i := 0
			otherClientsNames := make([]string, len(clts))
			for k := range clts {
				otherClientsNames[i] = k
				i++
			}
			for i := 0; i < otherNodesToAsk; i++ {
				id := rand.Intn(len(otherClientsNames))
				go askNode(clts[otherClientsNames[id]].client)
				sentAsks++
			}
		}
	}

	// Wait for results and return the first good one
	var res readResponse
	rectRes := 0
	for res = range responses {
		log.Debug().Msgf("...waiting for the first answer to return it. res=%#v", res)
		rectRes++
		if !res.error {
			log.Debug().Msgf("...got Response without error (closing channel): %#v", res)
			responsesClosed = true
			return &alexandraProto.ReadResponse{Data: res.data}, nil
		}
		if rectRes >= sentAsks {
			log.Warn().Msgf("ReadFromAnywhere: no fastestClient was able to answer the read (closing channel). Kg=%s", request.Keygroup)
			responsesClosed = true
			break
		}
	}

	// There was no successful response -- Update the keygroup information and try one last time
	log.Info().Msg("ReadFromAnywhere: Was not able to reach any queried node, updating cache and retrying...")
	m.updateKeygroupClients(request.Keygroup)

	client, err := m.GetFastestClientWithKeygroup(request.Keygroup, request.MinExpiry)
	if err != nil {
		return nil, fmt.Errorf("ReadFromAnywhere: there is no client with keygroup %s and expiry %d", request.Keygroup, request.MinExpiry)
	}
	result, err := client.Read(ctx, request.Keygroup, request.Id)
	if err != nil {
		return nil, fmt.Errorf("ReadFromAnywhere: cannot read from fastest client. err=%v", err)
	}

	return &alexandraProto.ReadResponse{Data: result.Data}, nil
}
