package fred

import (
	"context"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/trigger"
	"google.golang.org/grpc"
	"sync"
)

// Trigger is one trigger node with an ID and a host address.
type Trigger struct {
	ID   string
	Host string
}

type triggerMap struct {
	nodes map[string]*Trigger
	*sync.RWMutex
}

type triggerService struct {
	triggers map[KeygroupName]triggerMap
}

func getConnAndClient(host string) (client trigger.TriggerNodeClient, conn *grpc.ClientConn) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return nil, nil
	}
	log.Info().Msgf("Interclient: Created Connection to %s", host)
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

	if err != nil {
		return err
	} else if res.Status == trigger.EnumTriggerStatus_TRIGGER_ERROR {
		return errors.New(res.ErrorMessage)
	} else {
		return nil
	}
}

func newTriggerService() *triggerService {
	return &triggerService{triggers: make(map[KeygroupName]triggerMap)}
}

func (t *triggerService) triggerDelete(i Item) (e error) {
	log.Debug().Msgf("triggerDelete from triggerservice: in %#v", i)

	t.triggers[i.Keygroup].RLock()
	defer t.triggers[i.Keygroup].RUnlock()
	for _, node := range t.triggers[i.Keygroup].nodes {
		client, conn := getConnAndClient(node.Host)
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
	log.Debug().Msgf("triggerUpdate from triggerservice: in %#v, have %#v", i, t.triggers[i.Keygroup])

	t.triggers[i.Keygroup].RLock()
	defer t.triggers[i.Keygroup].RUnlock()
	for _, node := range t.triggers[i.Keygroup].nodes {
		client, conn := getConnAndClient(node.Host)
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

	if _, ok := t.triggers[k.Name]; !ok {
		return errors.Errorf("keygroup %s not in trigger store", k.Name)
	}

	t.triggers[k.Name].Lock()
	defer t.triggers[k.Name].Unlock()

	t.triggers[k.Name].nodes[tn.ID] = &tn

	return nil
}

func (t *triggerService) getTrigger(k Keygroup) ([]Trigger, error) {
	log.Debug().Msgf("getTrigger from triggerservice: in %#v", k)

	if _, ok := t.triggers[k.Name]; !ok {
		return nil, errors.Errorf("keygroup %s not in trigger store", k.Name)
	}

	t.triggers[k.Name].RLock()
	defer t.triggers[k.Name].RUnlock()

	var n []Trigger

	for _, node := range t.triggers[k.Name].nodes {
		n = append(n, *node)
	}

	return n, nil
}

func (t *triggerService) removeTrigger(k Keygroup, tn Trigger) error {
	log.Debug().Msgf("removeTrigger from triggerservice: in %#v %#v", k, tn)

	if _, ok := t.triggers[k.Name]; !ok {
		return errors.Errorf("keygroup %s not in trigger store", k.Name)
	}

	t.triggers[k.Name].Lock()
	defer t.triggers[k.Name].Unlock()

	delete(t.triggers[k.Name].nodes, tn.ID)

	return nil
}

func (t *triggerService) createKeygroup(k Keygroup) error {
	log.Debug().Msgf("createKeygroup from triggerservice: in %#v", k)

	t.triggers[k.Name] = triggerMap{
		nodes:   make(map[string]*Trigger),
		RWMutex: &sync.RWMutex{},
	}

	log.Debug().Msgf("createKeygroup from triggerservice: out %#v, have %#v", k, t.triggers[k.Name])

	return nil
}

func (t *triggerService) deleteKeygroup(name KeygroupName) error {
	log.Debug().Msgf("deleteKeygroup from triggerservice: in %#v", name)

	delete(t.triggers, name)
	return nil
}
