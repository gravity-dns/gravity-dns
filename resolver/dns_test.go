package resolver

import (
	"log"
	"testing"

	"golang.org/x/net/dns/dnsmessage"
)

func TestDot(t *testing.T) {
	resp, err := ResolverOverDNS("scottrichardson.dev", "A")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range resp {
		aResource := v.Body.(*dnsmessage.AResource)
		log.Println(aResource.A)
	}
	t.Fatal(resp)
}
