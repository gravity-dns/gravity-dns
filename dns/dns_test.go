package dns

import (
	"net"
	"testing"
)

func TestSetAndResolve(t *testing.T) {
	s := New()
	domain := "cool-domain.dev"
	ip := net.IPv4(0, 0, 0, 0)
	if _, err := s.AddNewEntry(AEntry, domain, ip); err != nil {
		t.Fatal(err)
	}
	if resolvedIP, err := s.RetrieveEntry(AEntry, domain); err != nil && resolvedIP[0].A.String() != ip.String() {
		t.Fatalf("Resolve expected %s got %s\n", ip.String(), resolvedIP[0].A.String())
	}
}
