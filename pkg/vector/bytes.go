package vector

import (
	"bytes"
	"encoding/gob"
	"sort"

	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/rs/zerolog/log"
)

type Entry struct {
	ID    string
	Ticks uint64
}

func Bytes(v vclock.VClock) []byte {
	ids := make([]string, len(v))
	i := 0
	for id := range v {
		ids[i] = id
		i++
	}

	sort.Strings(ids)

	entries := make([]Entry, len(v))

	for i := range ids {
		entries[i] = Entry{
			ID:    ids[i],
			Ticks: v[ids[i]],
		}
	}

	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	err := enc.Encode(entries)

	if err != nil {
		log.Fatal().Msgf("vector clock encode: %+v", err)
	}

	return b.Bytes()
}

func FromBytes(data []byte) (vclock.VClock, error) {
	b := new(bytes.Buffer)
	b.Write(data)

	entries := make([]Entry, 0)
	dec := gob.NewDecoder(b)

	err := dec.Decode(&entries)

	if err != nil {
		return nil, err
	}

	clock := vclock.New()
	for _, e := range entries {
		clock.Set(e.ID, e.Ticks)
	}

	return clock, nil
}
