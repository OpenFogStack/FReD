package fred

import (
	"context"
	"crypto/tls"

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
}

func (t *triggerService) getConnAndClient(host string) (client trigger.TriggerNodeClient, conn *grpc.ClientConn) {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(*t.tc))
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return nil, nil
	}
	log.Debug().Msgf("Interclient: Created Connection to %s", host)
	client = trigger.NewTriggerNodeClient(conn)
	return
}

// logs the response and returns the correct error message
func dealWithStatusResponse(res *trigger.TriggerResponse, err error, from string) error {
	if res != nil {
		log.Debug().Msgf("Interclient got Response from %s, Status %s with Message %s and Error %s", from, res.Status, res.ErrorMessage, err)
	} else {
		log.Debug().Msgf("Interclient got empty Response from %s", from)
	}

	if err != nil || res == nil {
		return err
	}

	if res.Status == trigger.EnumTriggerStatus_TRIGGER_ERROR {
		return errors.New(res.ErrorMessage)
	}

	return nil

}

func newTriggerService(s *storeService, certFile, keyFile string) *triggerService {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load certificates")
		return nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	tc := credentials.NewTLS(tlsConfig)

	return &triggerService{
		tc: &tc,
		s:  s,
	}
}

func (t *triggerService) triggerDelete(i Item) (e error) {
	log.Debug().Msgf("triggerDelete from triggerservice: in %#v", i)

	nodes, err := t.s.getKeygroupTrigger(i.Keygroup)

	if err != nil {
		return err
	}

	for _, node := range nodes {
		client, conn := t.getConnAndClient(node.Host)
		res, err := client.DeleteItemTrigger(context.Background(), &trigger.DeleteItemTriggerRequest{
			Keygroup: string(i.Keygroup),
			Id:       i.ID,
		})

		err = dealWithStatusResponse(res, err, "TriggerDelete")

		if err != nil {
			e = errors.Errorf("%#v %#v", e, err)
		}

		err = conn.Close()

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
		return err
	}

	for _, node := range nodes {
		client, conn := t.getConnAndClient(node.Host)
		res, err := client.PutItemTrigger(context.Background(), &trigger.PutItemTriggerRequest{
			Keygroup: string(i.Keygroup),
			Id:       i.ID,
			Val:      i.Val,
		})

		err = dealWithStatusResponse(res, err, "TriggerUpdate")

		if err != nil {
			e = errors.Errorf("%#v %#v", e, err)
		}

		err = conn.Close()

		if err != nil {
			e = errors.Errorf("%#v %#v", e, err)
		}
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
