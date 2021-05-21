package dns

import (
	"net"
	"testing"
)

func TestCanRetrieveValue(t *testing.T) {
	entries := NewDNSEntries()
	entries.AddNewEntry(AEntry, "scottrichardson.dev", net.IPv4(1, 1, 1, 1))
	entries.AddNewEntry(AEntry, "dev.scottrichardson.dev", net.IPv4(4, 0, 0, 0))

	found := entries.RetrieveEntry(AEntry, "scottrichardson.dev")
	if found.Value.String() != net.IPv4(1, 1, 1, 1).String() {
		t.Fatalf("Expected 1.1.1.1 got %s", found.Value.String())
	}

	found = entries.RetrieveEntry(AEntry, "dev.scottrichardson.dev")
	if found == nil || found.Value.String() != net.IPv4(4, 0, 0, 0).String() {
		t.Fatalf("Expected 0.0.0.0 got %v", found)
	}
}

func TestCanRetrieveDifferentType(t *testing.T) {
	entries := NewDNSEntries()
	entries.AddNewEntry(AEntry, "scott.richardson", net.IPv4(1, 1, 1, 1))
	entries.AddNewEntry(AAAAEntry, "scott.richardson", net.IPv4(0, 1, 1, 0))

	found := entries.RetrieveEntry(AEntry, "scott.richardson")
	if found == nil || found.Value.String() != net.IPv4(1, 1, 1, 1).String() {
		t.Fatalf("Expected 1.1.1.1 got %v", found)
	}

	found = entries.RetrieveEntry(AAAAEntry, "scott.richardson")
	if found == nil || found.Value.String() != net.IPv4(0, 1, 1, 0).String() {
		t.Fatalf("Expected 0.1.1.0 got %v\n", found)
	}
}
