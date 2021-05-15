package dns

import (
	"fmt"
	"os"
	"testing"
)

var sink *gravityDNS

const (
	adserverHosts = "../hostfiles/adservers.txt"
	facebookHosts = "../hostfiles/facebook.txt"
)

func init() {
	sink = new()
	sink.ParseAdFile(adserverHosts)
	sink.ParseAdFile(facebookHosts)
}

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

	if len(s.domains) == 0 {
		t.Fail()
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

func TestIsAdIP(t *testing.T) {
	badIP := "0.0.0.0"
	goodIP := "1.1.1.1"
	if !IsAdIP(badIP) {
		t.Fatalf("IP: %s should be an ad\n", badIP)
	}
	if IsAdIP(goodIP) {
		t.Fatalf("IP: %s should be a good ip not an ad\n", goodIP)
	}
}

func BenchmarkStandardLib(b *testing.B) {
	fmt.Printf("Size of map - %d\n", len(sink.domains))
	for n := 0; n < b.N; n++ {
		if sink.Resolve("00v07c3k7o.kameleoon.eu") != "0.0.0.0" {
			b.Fail()
		}
	}
}
