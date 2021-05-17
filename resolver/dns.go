package resolver

import (
	"log"
	"math"
	"math/rand"
	"net"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

const (
	udpPort      = 53
	queryUDPPort = 8008
	packetLen    = 512
)

func getNewID() uint32 {
	// Seed the random number generator using the current time (nanoseconds since epoch)
	rand.Seed(time.Now().UnixNano())

	// Much harder to predict...but it is still possible if you know the day, and hour, minute...
	return rand.Uint32() % uint32(math.Pow(2, 15))
}

func createDNSQuery(id uint16, domain, queryType string) ([]byte, error) {
	msg := dnsmessage.NewBuilder(nil, dnsmessage.Header{
		ID:                 id,
		Response:           false,
		OpCode:             0,
		Authoritative:      false,
		Truncated:          false,
		RecursionDesired:   true,
		RecursionAvailable: false,
		RCode:              dnsmessage.RCodeSuccess,
	})

	if err := msg.StartQuestions(); err != nil {
		return nil, err
	}

	name, err := dnsmessage.NewName(domain + ".")
	log.Println("domain " + name.String())
	if err != nil {
		return nil, err
	}

	msg.Question(dnsmessage.Question{
		Name:  name,
		Type:  dnsmessage.TypeA,
		Class: dnsmessage.ClassINET,
	})

	msg.EnableCompression()
	return msg.Finish()
}

func ResolverOverDNS(domain, queryType string) ([]dnsmessage.Resource, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: queryUDPPort})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	queryID := uint16(getNewID())
	log.Printf("queryID %v", queryID)
	query, err := createDNSQuery(queryID, domain, queryType)
	if err != nil {
		return nil, err
	}

	if _, err = conn.WriteToUDP(query, &net.UDPAddr{
		IP:   net.IPv4(1, 1, 1, 1),
		Port: udpPort,
	}); err != nil {
		return nil, err
	}

	for {
		buf := make([]byte, packetLen)
		_, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			return nil, err
		}

		m := dnsmessage.Message{}
		if err = m.Unpack(buf); err != nil {
			return nil, err
		}
		if !m.Header.Response || m.Header.ID != queryID {
			continue
		}
		log.Println(m)
		return m.Answers, nil
	}
}
