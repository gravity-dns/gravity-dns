package dns

import (
	"net"
	"sync"

	"github.com/gravity-dns/gravity-dns/resolver"
)

type (
	GravityDNS interface {
		DNSEntries
		Sinkhole

		NumEntries() int
		RetrieveAndSet(entryType EntryType, domain string) ([]*EntryValue, error)
	}

	gravityDNS struct {
		numEntries int
		dns        DNSEntries
		lock       sync.RWMutex
	}
)

func New() GravityDNS {
	return new()
}

func new() *gravityDNS {
	return &gravityDNS{
		dns:        NewDNSEntries(),
		lock:       sync.RWMutex{},
		numEntries: 0,
	}
}

func (g *gravityDNS) RetrieveAndSet(entryType EntryType, domain string) ([]*EntryValue, error) {
	values := []*EntryValue{}
	resolved, err := g.RetrieveEntry(entryType, domain)
	if err != nil && err.Error() != ErrDomainNotFound {
		return values, err
	}

	if resolved != nil && err == nil {
		return append(values, resolved...), nil
	}

	// We don't have the value so we need to retrieve it
	resp, err := resolver.ResolveOverDoH(domain, GravtiyTypeToDNSMessage(entryType))
	if err != nil {
		return nil, err
	}

	for _, val := range resp.Answer {
		var data interface{}
		if entryType == AEntry || entryType == AAAAEntry {
			data = net.ParseIP(val.Data)
		} else {
			data = val.Data
		}
		newValue, err := newEntryValue(entryType, data)
		if err != nil {
			return values, nil
		}

		values = append(values, &newValue)
		g.AddNewEntry(entryType, domain, data)
	}

	return values, nil
}

func (g *gravityDNS) AddNewEntry(entryType EntryType, domain string, value interface{}) (*Entry, error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.numEntries++
	return g.dns.AddNewEntry(entryType, domain, value)
}

func (g *gravityDNS) RetrieveEntry(entryType EntryType, domain string) ([]*EntryValue, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.dns.RetrieveEntry(entryType, domain)
}

func (g *gravityDNS) NumEntries() int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.numEntries
}
