package dns

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	rootEntry       = iota
	subdoamainEntry = iota

	AEntry     = iota
	AAAAEntry  = iota
	CNAMEEntry = iota
	MXEntry    = iota
	TXTEntry   = iota
	NSEntry    = iota
	SOAEntry   = iota
	SRVEntry   = iota
	PTREntry   = iota
)

type (
	EntryType int

	Entry struct {
		EntryType EntryType
		Value     net.IP
		Entries   map[string]*Entry
	}

	DNSEntries interface {
		AddNewEntry(EntryType, string, net.IP) *Entry
		RetrieveEntry(EntryType, string) *Entry
	}
)

func NewDNSEntries() DNSEntries {
	return newEntry(rootEntry)
}

func newEntry(entryType EntryType) *Entry {
	return &Entry{
		EntryType: entryType,
		Entries:   make(map[string]*Entry),
	}
}

func (e *Entry) AddNewEntry(entryType EntryType, domain string, value net.IP) *Entry {
	splits := strings.Split(domain, ".")
	currEntry := e.Entries[fmt.Sprint(entryType)]
	if currEntry == nil {
		createdEntry := newEntry(entryType)
		currEntry = createdEntry
		e.Entries[fmt.Sprint(entryType)] = createdEntry
	}

	prev := currEntry
	for i := len(splits) - 1; i >= 0; i-- {
		log.Println(prev)
		currEntry = prev.Entries[splits[i]]
		if currEntry == nil {
			createdEntry := newEntry(subdoamainEntry)
			currEntry = createdEntry
			prev.Entries[splits[i]] = createdEntry
		}
		prev = currEntry
	}
	prev.Value = value
	return currEntry
}

func (e *Entry) RetrieveEntry(entryType EntryType, domain string) *Entry {
	splits := strings.Split(domain, ".")
	currEntry := e.Entries[fmt.Sprint(entryType)]
	if currEntry == nil {
		return nil
	}

	prev := currEntry
	for i := len(splits) - 1; i >= 0; i-- {
		currEntry = prev.Entries[splits[i]]
		if currEntry == nil {
			return nil
		}
		prev = currEntry
	}

	return currEntry
}
