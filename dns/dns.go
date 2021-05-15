package dns

import "sync"

type (
	GravityDNS interface {
		Resolve(domain string) string
		Set(domain, up string)
		ParseAdFile(filename string) error
	}

	gravityDNS struct {
		domains map[string]string
		lock    sync.RWMutex
	}
)

func New() GravityDNS {
	return new()
}

func new() *gravityDNS {
	return &gravityDNS{
		domains: make(map[string]string),
		lock:    sync.RWMutex{},
	}
}

func (s *gravityDNS) Resolve(domain string) string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.domains[domain]
}

func (s *gravityDNS) Set(domain, ip string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.domains[domain] = ip
}
