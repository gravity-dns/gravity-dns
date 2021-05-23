package dns

import (
	"net"
	"testing"
)

func TestSetAndResolve(t *testing.T) {
	s := New()
	domain := "cool-domain.dev"
	ip := net.IPv4(0, 0, 0, 0)
	s.AddNewEntry(AEntry, domain, ip)
	if resolvedIP, err := s.RetrieveEntry(AEntry, domain); err != nil && resolvedIP.A.String() != ip.String() {
		t.Fatalf("Resolve expected %s got %s\n", ip.String(), resolvedIP.A.String())
	}
}
