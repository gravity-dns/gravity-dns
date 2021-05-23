package dns

import (
	"net"
	"os"
	"testing"
)

const (
	adserverHosts = "../hostfiles/adservers.txt"
	facebookHosts = "../hostfiles/facebook.txt"
)

func createFakeFile(filename, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func deleteFakeFile(filename string) error {
	return os.Remove(filename)
}

func TestParser(t *testing.T) {
	s := new()
	if err := s.ParseAdFile(adserverHosts); err != nil {
		t.Fatal(err)
	}

	if s.NumEntries() == 0 {
		t.Fatalf("Num entries should be grater than zero")
	}
}

func TestThrowsErrorWhenNoFile(t *testing.T) {
	s := New()
	if err := s.ParseAdFile("i-dont-exists"); err == nil {
		t.Fatal("Parse file should have thrown error for file not found")
	}
}

func TestThrowsAnErrorWhenInvalidFile_TooManyParts(t *testing.T) {
	s := New()
	createFakeFile("test.txt", "i am invalid because I have too many parts")
	defer func() {
		if err := deleteFakeFile("test.txt"); err != nil {
			t.Fatal(err)
		}
	}()
	if err := s.ParseAdFile("test.txt"); err == nil {
		t.Fatal("Expected error for invalid file")
	}
}

func TestThrowsAnErrorWhenInvalidFile_TooFewParts(t *testing.T) {
	s := New()
	createFakeFile("test.txt", "invalid")
	defer func() {
		if err := deleteFakeFile("test.txt"); err != nil {
			t.Fatal(err)
		}
	}()
	if err := s.ParseAdFile("test.txt"); err == nil {
		t.Fatal("Expected error for invalid file")
	}
}

func BenchmarkStandardLib(b *testing.B) {
	sink := new()
	sink.ParseAdFile(adserverHosts)
	sink.ParseAdFile(facebookHosts)

	for n := 0; n < b.N; n++ {
		found, err := sink.RetrieveEntry(AEntry, "00v07c3k7o.kameleoon.eu")
		if err != nil {
			b.Fatal(err)
		} else if found == nil {
			b.Fatal("Retrived value is nil")
		} else if found.A.String() != net.IPv4(0, 0, 0, 0).String() {
			b.Fatal("Retrieved value is not correcte")
		}
	}
}

func BenchmarkRetrievalEmpty(b *testing.B) {
	sink := new()
	sink.AddNewEntry(AEntry, "scott.dev", net.IPv4(0, 0, 0, 0))

	for n := 0; n < b.N; n++ {
		if _, err := sink.RetrieveEntry(AEntry, "scott.dev"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRetrievalFull(b *testing.B) {
	sink := new()
	sink.ParseAdFile(facebookHosts)
	sink.ParseAdFile(adserverHosts)

	sink.AddNewEntry(AEntry, "scott.dev", net.IPv4(0, 0, 0, 0))

	for n := 0; n < b.N; n++ {
		if _, err := sink.RetrieveEntry(AEntry, "scott.dev"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWithoutLock(b *testing.B) {
	sink := NewDNSEntries()
	sink.AddNewEntry(AEntry, "scott.dev", net.IPv4(0, 0, 0, 0))
	for n := 0; n < b.N; n++ {
		if _, err := sink.RetrieveEntry(AEntry, "scott.dev"); err != nil {
			b.Fatal(err)
		}
	}
}
