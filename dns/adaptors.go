package dns

import (
	"errors"

	"golang.org/x/net/dns/dnsmessage"
)

const (
	ErrInvalidIPSize = "invalid size of IP address"
	ErrInvalidValue  = "the value provided was invalid"
)

func DNSTypeToGravity(t dnsmessage.Type) EntryType {
	switch t {
	case dnsmessage.TypeA:
		return AEntry
	case dnsmessage.TypeAAAA:
		return AAAAEntry
	case dnsmessage.TypeCNAME:
		return CNAMEEntry
	case dnsmessage.TypeMX:
		return MXEntry
	case dnsmessage.TypeTXT:
		return TXTEntry
	case dnsmessage.TypePTR:
		return PTREntry
	default:
		return InvalidEntry
	}
}

func GravtiyTypeToDNSMessage(t EntryType) dnsmessage.Type {
	switch t {
	case AEntry:
		return dnsmessage.TypeA
	case AAAAEntry:
		return dnsmessage.TypeAAAA
	case CNAMEEntry:
		return dnsmessage.TypeCNAME
	case MXEntry:
		return dnsmessage.TypeMX
	case TXTEntry:
		return dnsmessage.TypeTXT
	case PTREntry:
		return dnsmessage.TypePTR
	default:
		return dnsmessage.TypeALL
	}
}

func GravityEntryToResourceBody(entryType EntryType, data *EntryValue) (dnsmessage.ResourceBody, error) {
	switch entryType {
	case AEntry:
		val := data.A.To4()
		if len(val) < 4 {
			return nil, errors.New(ErrInvalidIPSize)
		}
		byteArr := [4]byte{}
		copy(byteArr[:], val[0:4])
		return &dnsmessage.AResource{A: byteArr}, nil
	case AAAAEntry:
		val := data.AAAA.To16()
		if len(val) < 16 {
			return nil, errors.New(ErrInvalidIPSize)
		}
		byteArr := [16]byte{}
		copy(byteArr[:], val[0:16])
		return &dnsmessage.AAAAResource{AAAA: byteArr}, nil
	case CNAMEEntry:
		name, err := dnsmessage.NewName(data.CNAME)
		if err != nil {
			return nil, errors.New(ErrInvalidValue)
		}
		return &dnsmessage.CNAMEResource{CNAME: name}, nil
	case MXEntry:
		name, err := dnsmessage.NewName(data.MX)
		if err != nil {
			return nil, errors.New(ErrInvalidValue)
		}
		return &dnsmessage.MXResource{MX: name}, nil
	case TXTEntry:
		return &dnsmessage.TXTResource{TXT: []string{data.TXT}}, nil
	case PTREntry:
		name, err := dnsmessage.NewName(data.PTR)
		if err != nil {
			return nil, errors.New(ErrInvalidValue)
		}
		return &dnsmessage.PTRResource{PTR: name}, nil
	}

	return nil, errors.New(ErrInvalidType)
}
