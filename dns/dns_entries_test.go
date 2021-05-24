package dns

import (
	"errors"
	"net"
	"strings"
	"testing"
)

func TestCanRetrieveValue(t *testing.T) {
	entries := NewDNSEntries()
	entries.AddNewEntry(AEntry, "scottrichardson.dev", net.IPv4(1, 1, 1, 1))
	entries.AddNewEntry(AEntry, "dev.scottrichardson.dev", net.IPv4(4, 0, 0, 0))

	found, err := entries.RetrieveEntry(AEntry, "scottrichardson.dev")
	if err != nil {
		t.Fatal(err)
	}

	if found[0].A.String() != net.IPv4(1, 1, 1, 1).String() {
		t.Fatalf("Expected 1.1.1.1 got %v", found)
	}

	found, err = entries.RetrieveEntry(AEntry, "dev.scottrichardson.dev")
	if err != nil {
		t.Fatal(err)
	}
	if found[0].A.String() != net.IPv4(4, 0, 0, 0).String() {
		t.Fatalf("Expected 0.0.0.0 got %v", found)
	}
}

func TestCanRetrieveDifferentType(t *testing.T) {
	entries := NewDNSEntries()
	entries.AddNewEntry(AEntry, "scott.richardson", net.IPv4(1, 1, 1, 1))
	entries.AddNewEntry(AAAAEntry, "scott.richardson", net.ParseIP("::1"))

	found, err := entries.RetrieveEntry(AEntry, "scott.richardson")
	if err != nil {
		t.Fatal(err)
	}
	if found[0].A.String() != net.IPv4(1, 1, 1, 1).String() {
		t.Fatalf("Expected 1.1.1.1 got %v", found[0].A.String())
	}

	found, err = entries.RetrieveEntry(AAAAEntry, "scott.richardson")
	if err != nil {
		t.Fatal(err)
	}

	if found[0].AAAA.String() != net.ParseIP("::1").String() {
		t.Fatalf("Expected ::1 got %v\n", found)
	}
}

func TestCanAddAndRetrieveTextValues(t *testing.T) {
	entries := NewDNSEntries()
	domain := "scottrichardson.dev."
	entry := "iamtext"
	if _, err := entries.AddNewEntry(TXTEntry, domain, "iamtext"); err != nil {
		t.Fatal(err)
	}

	if found, err := entries.RetrieveEntry(TXTEntry, domain); err != nil {
		t.Fatal(err)
	} else if found == nil {
		t.Fatal("Found is nil")
	} else {
		if found[0].TXT != entry {
			t.Fatalf("TXT entry invalid. Expected %s got %s", domain, found[0].TXT)
		}
	}
}

func TestThrowsErrorOnInvalidDomain(t *testing.T) {
	entries := NewDNSEntries()
	invalidDomains := []string{
		".dev",
		"dev",
		"",
	}
	for _, domain := range invalidDomains {
		_, err := entries.AddNewEntry(AEntry, domain, net.IPv4(0, 0, 0, 0))
		if err == nil || err.Error() != ErrInvalidDomain {
			t.Fatalf("AddNewEntry should have failed from domain %s", domain)
		}
	}
}

func TestThrowsErrorOnInvalidValue(t *testing.T) {
	entries := NewDNSEntries()
	domain := "gravity-dns.com"
	invalidValues := []interface{}{
		struct {
			Name string
		}{"test"},
		errors.New("I am invalid"),
	}
	for _, value := range invalidValues {
		if _, err := entries.AddNewEntry(AEntry, domain, value); err == nil || err.Error() != ErrInvalidType {
			t.Fatalf("AddNewEntry should have failed for type %v", value)
		}
	}
}

func TestThrowsErrorOnInvalidIPType(t *testing.T) {
	entries := NewDNSEntries()
	domain := "gravity-dns.com"
	invalidValues := []interface{}{
		struct {
			Name string
		}{"test"},
		errors.New("I am invalid"),
		"Iaminvalid",
	}
	for _, value := range invalidValues {
		if _, err := entries.AddNewEntry(AEntry, domain, value); err == nil || err.Error() != ErrInvalidType {
			t.Fatalf("AddNewEntry should have failed for type %v", value)
		}
	}
}

func TestThrowsErrorOnInvalidStringType(t *testing.T) {
	entries := NewDNSEntries()
	domain := "gravity-dns.com"
	invalidValues := []interface{}{
		struct {
			Name string
		}{"test"},
		errors.New("I am invalid"),
		net.IPv4(0, 0, 0, 0),
	}
	for _, value := range invalidValues {
		if _, err := entries.AddNewEntry(TXTEntry, domain, value); err == nil || err.Error() != ErrInvalidType {
			t.Fatalf("AddNewEntry should have failed for type %v", value)
		}
	}
}

func BenchmarkStringtoken(b *testing.B) {
	stringToToken := "i.token.string.cool.very"
	for n := 0; n < b.N; n++ {
		strings.Split(stringToToken, ".")
	}
}
