package fred

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"git.tu-berlin.de/mcc-fred/fred/proto/trigger"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Trigger is one trigger node with an ID and a host address.
type Trigger struct {
	ID   string
	Host string
}

type triggerService struct {
	tc *credentials.TransportCredentials
	s  *storeService
	c  map[string]trigger.TriggerNodeClient
}

func (t *triggerService) getClient(host string) (trigger.TriggerNodeClient, error) {
	if client, ok := t.c[host]; ok {
		return client, nil
	}

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(*t.tc))

	if err != nil {
		log.Error().Err(err).Msg("Cannot create Grpc connection")
		return nil, err
	}

	log.Debug().Msgf("Interclient: Created Connection to %s", host)

	t.c[host] = trigger.NewTriggerNodeClient(conn)
	return t.c[host], nil
}

func newTriggerService(s *storeService, certFile string, keyFile string, caFiles []string) *triggerService {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load certificates")
		return nil
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs, err := x509.SystemCertPool()

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load root certificates")
		return nil
	}

	for _, f := range caFiles {
		loaded, err := ioutil.ReadFile(f)

		if err != nil {
			log.Fatal().Msgf("unexpected missing certfile: %v", err)
		}

		rootCAs.AppendCertsFromPEM(loaded)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      rootCAs,
	}

	tc := credentials.NewTLS(tlsConfig)

	return &triggerService{
		tc: &tc,
		s:  s,
		c:  make(map[string]trigger.TriggerNodeClient),
	}
}

func (t *triggerService) triggerDelete(i Item) (e error) {
	log.Debug().Msgf("triggerDelete from triggerservice: in %#v", i)

	nodes, err := t.s.getKeygroupTrigger(i.Keygroup)

	if err != nil {
		return err
	}

	for _, node := range nodes {
		client, err := t.getClient(node.Host)
		if err != nil {
			e = errors.Errorf("%#v %#v", e, err)
			continue
		}
		_, err = client.DeleteItemTrigger(context.Background(), &trigger.DeleteItemTriggerRequest{
			Keygroup: string(i.Keygroup),
			Id:       i.ID,
		})

		if err != nil {
			e = errors.Errorf("%#v %#v", e, err)
		}
	}

	return e
}

func (t *triggerService) triggerUpdate(i Item) (e error) {
	log.Debug().Msgf("triggerUpdate from triggerservice: in %#v", i)

	nodes, err := t.s.getKeygroupTrigger(i.Keygroup)

	if err != nil {
		log.Error().Msgf("error in triggerUpdate: %s", err.Error())
		return err
	}

	for _, node := range nodes {
		log.Debug().Msgf("triggerUpdate: updating %s %s", node.ID, node.Host)
		client, err := t.getClient(node.Host)
		if err != nil {
			log.Error().Msgf("error in triggerUpdate: %s", err.Error())
			e = errors.Errorf("%#v\n%s", e, err.Error())
			continue
		}
		_, err = client.PutItemTrigger(context.Background(), &trigger.PutItemTriggerRequest{
			Keygroup: string(i.Keygroup),
			Id:       i.ID,
			Val:      i.Val,
		})

		if err != nil {
			log.Error().Msgf("error in triggerUpdate: %s", err.Error())
			e = errors.Errorf("%#v\n%s", e, err.Error())
		}
	}

	if e != nil {
		log.Error().Msgf("error in triggerUpdate: %s", e.Error())
	}

	return e
}

func (t *triggerService) addTrigger(k Keygroup, tn Trigger) error {
	log.Debug().Msgf("addTrigger from triggerservice: in %#v %#v", k, tn)

	return t.s.addKeygroupTrigger(k.Name, tn)
}

func (t *triggerService) getTrigger(k Keygroup) ([]Trigger, error) {
	log.Debug().Msgf("getTrigger from triggerservice: in %#v", k)

	return t.s.getKeygroupTrigger(k.Name)
}

func (t *triggerService) removeTrigger(k Keygroup, tn Trigger) error {
	log.Debug().Msgf("removeTrigger from triggerservice: in %#v %#v", k, tn)

	return t.s.deleteKeygroupTrigger(k.Name, tn)
}
