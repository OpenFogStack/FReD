package fred

import (
	"context"
	"crypto/tls"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/proto/trigger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Trigger is one trigger node with an ID and a host address.
type Trigger struct {
	ID   string
	Host string
}

type triggerService struct {
	tc           *credentials.TransportCredentials
	s            *storeService
	MissedDelete map[string][]Item
	MissedUpdate map[string][]Item
}

func newTriggerService(s *storeService, certFile, keyFile string) *triggerService {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load certificates")
		return nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	tc := credentials.NewTLS(tlsConfig)

	return &triggerService{
		tc:           &tc,
		s:            s,
		MissedDelete: make(map[string][]Item),
		MissedUpdate: make(map[string][]Item),
	}
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
func dealWithStatusResponse(res *trigger.TriggerResponse, err error, from string) {
	if res != nil {
		log.Debug().Msgf("Interclient got Response from %s, Status %s with Message %s and Error %s", from, res.Status, res.ErrorMessage, err)
	} else {
		log.Debug().Msgf("Interclient got empty Response from %s", from)
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

		dealWithStatusResponse(res, err, "TriggerDelete")

		if err != nil {
			// Message did not reach trigger, so store in messages-to-send
			// But only if there are less than 1000 messages, if someone forgot to delete a trigger...
			if len(t.MissedDelete[node.ID]) < 1000 {
				t.MissedDelete[node.ID] = append(t.MissedDelete[node.ID], i)
				log.Info().Msgf("Trigger: triggerDelete: was not able to reach trigger, saving message to send in the future")
			} else {
				log.Warn().Msgf("Trigger: triggerDelete: node %s has missed more than 1000 triggers, so this trigger will be lost", node.ID)
			}
		} else {
			t.checkMissedTriggers(node)
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

		dealWithStatusResponse(res, err, "TriggerUpdate")

		if err != nil {
			// Message did not reach trigger, so store in messages-to-send
			// But only if there are less than 1000 messages, if someone forgot to delete a trigger...
			if len(t.MissedUpdate[node.ID]) < 1000 {
				t.MissedUpdate[node.ID] = append(t.MissedUpdate[node.ID], i)
				log.Info().Msgf("Trigger: triggerUpdate: was not able to reach trigger, saving message to send in the future")
			} else {
				log.Warn().Msgf("Trigger: triggerUpdate: node %s has missed more than 1000 triggers, so this trigger will be lost", node.ID)
			}
		} else {
			t.checkMissedTriggers(node)
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
	t.MissedDelete[tn.ID] = nil
	t.MissedUpdate[tn.ID] = nil

	return t.s.deleteKeygroupTrigger(k.Name, tn)
}

// checkMissedTriggers checks whether this trigger node has missed any updates
func (t *triggerService) checkMissedTriggers(node Trigger) {
	missedU := t.MissedUpdate[node.ID]
	missedD := t.MissedDelete[node.ID]
	// Missed Updates:
	if len(missedU) != 0 {
		log.Info().Msgf("triggerservice: checkMissed: node %s is reachable again and will receive %d updates", node.ID, len(missedU))
		client, conn := t.getConnAndClient(node.Host)
		for count, i := range missedU {
			_, err := client.PutItemTrigger(context.Background(), &trigger.PutItemTriggerRequest{
				Keygroup: string(i.Keygroup),
				Id:       i.ID,
				Val:      i.Val,
			})
			if err == nil {
				// No errors == remove from Array
				t.MissedUpdate[node.ID] = append(t.MissedUpdate[node.ID][:count], t.MissedUpdate[node.ID][count+1:]...)
			}
		}
		conn.Close()
	}
	// Missed Deletes
	if len(missedD) != 0 {
		log.Info().Msgf("triggerservice: checkMissed: node %s is reachable again and will receive %d deletes", node.ID, len(missedU))
		client, conn := t.getConnAndClient(node.Host)
		for count, i := range missedD {
			_, err := client.DeleteItemTrigger(context.Background(), &trigger.DeleteItemTriggerRequest{
				Keygroup: string(i.Keygroup),
				Id:       i.ID,
			})
			if err == nil {
				// No errors == remove from Array
				t.MissedDelete[node.ID] = append(t.MissedDelete[node.ID][:count], t.MissedDelete[node.ID][count+1:]...)
			}
		}
		conn.Close()
	}
}
