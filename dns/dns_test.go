package dns

import "testing"

func TestSetAndResolve(t *testing.T) {
	s := New()
	domain := "cool-domain"
	ip := "0.0.0.0"
	s.Set(domain, ip)
	if resolvedIP := s.Resolve(domain); resolvedIP != ip {
		t.Fatalf("Resolve expected %s got %s\n", ip, resolvedIP)
	}
}
