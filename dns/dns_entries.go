package dns

import (
	"errors"
	"net"
	"strings"
)

const (
	AEntry     EntryType = iota
	AAAAEntry  EntryType = iota
	CNAMEEntry EntryType = iota
	MXEntry    EntryType = iota
	TXTEntry   EntryType = iota
	PTREntry   EntryType = iota

	InvalidEntry EntryType = iota
)

const (
	ErrInvalidDomain  = "domain must have at least two parts"
	ErrInvalidType    = "value provided was invalid. Must be either string or net.IP"
	ErrDomainNotFound = "domain specified was not found"
)

type (
	EntryType int
	Entry     struct {
		Name    string
		Entries map[string]*Entry
		Values  []EntryValue
	}

	EntryValue struct {
		EntryType EntryType

		A     net.IP
		AAAA  net.IP
		MX    string
		TXT   string
		CNAME string
		PTR   string
	}
	DNSEntries interface {
		AddNewEntry(EntryType, string, interface{}) (*Entry, error)
		RetrieveEntry(EntryType, string) ([]*EntryValue, error)
	}
)

func NewDNSEntries() DNSEntries {
	return newEntry("root")
}

func newEntry(name string) *Entry {
	return &Entry{
		Name:    name,
		Entries: make(map[string]*Entry),
		Values:  []EntryValue{},
	}
}

func (e *Entry) AddNewEntry(entryType EntryType, domain string, value interface{}) (*Entry, error) {
	splits := strings.Split(domain, ".")
	if len(splits) < 2 {
		return nil, errors.New(ErrInvalidDomain)
	}

	for i, split := range splits {
		if split == "" && i != len(splits)-1 {
			return nil, errors.New(ErrInvalidDomain)
		}
	}

	prev := e
	for i := len(splits) - 1; i >= 0; i-- {
		if splits[i] == "" && i == len(splits)-1 {
			continue
		}
		curr := prev.Entries[splits[i]]
		if curr == nil {
			curr = newEntry(splits[i])
			prev.Entries[splits[i]] = curr
		}
		prev = curr
	}
	entryValue, err := newEntryValue(entryType, value)
	if err != nil {
		return nil, err
	}

	prev.Values = append(prev.Values, entryValue)
	return prev, nil
}

func (e *Entry) RetrieveEntry(entryType EntryType, domain string) ([]*EntryValue, error) {
	splits := strings.Split(domain, ".")
	if len(splits) < 2 {
		return nil, errors.New(ErrInvalidDomain)
	}

	prev := e
	for i := len(splits) - 1; i >= 0; i-- {
		if splits[i] == "" {
			if i == len(splits)-1 {
				continue
			}
			return nil, errors.New(ErrInvalidDomain)
		}
		curr := prev.Entries[splits[i]]
		if curr == nil {
			return nil, errors.New(ErrDomainNotFound)
		}
		prev = curr
	}

	values := []*EntryValue{}
	for _, val := range prev.Values {
		if val.EntryType == entryType {
			values = append(values, &val)
		}
	}
	if len(values) == 0 {
		return values, errors.New(ErrDomainNotFound)
	}
	return values, nil
}

func newEntryValue(entryType EntryType, value interface{}) (EntryValue, error) {
	newValue := EntryValue{EntryType: entryType}
	strVal, strType := value.(string)
	ipVal, ipType := value.(net.IP)

	switch true {
	case entryType == AEntry && ipType:
		newValue.A = ipVal
	case entryType == AAAAEntry && ipType:
		newValue.AAAA = value.(net.IP)
	case entryType == CNAMEEntry && strType:
		newValue.CNAME = strVal
	case entryType == MXEntry && strType:
		newValue.MX = strVal
	case entryType == TXTEntry && strType:
		newValue.TXT = strVal
	case entryType == PTREntry && strType:
		newValue.PTR = strVal
	default:
		return newValue, errors.New(ErrInvalidType)
	}
	return newValue, nil
}
