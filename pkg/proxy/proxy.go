package proxy

import (
	"math/rand"

	"github.com/cespare/xxhash/v2"
)

type Proxy struct {
	hosts []string
}

func NewProxy(hosts []string) *Proxy {
	return &Proxy{
		hosts: hosts,
	}
}

func (p *Proxy) getHost(keygroup string) string {
	return p.hosts[xxhash.Sum64String(keygroup)%uint64(len(p.hosts))]
}

func (p *Proxy) getAny() string {
	return p.hosts[rand.Intn(len(p.hosts))]
}
