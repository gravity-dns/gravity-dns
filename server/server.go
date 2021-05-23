package server

import (
	"net"

	"github.com/gravity-dns/gravity-dns/dns"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/dns/dnsmessage"
)

const (
	// DNS server default port
	udpPort int = 53
	// DNS packet max length
	packetLen int = 512
)

func Start() error {
	gravity := dns.New()
	if err := gravity.ParseAdFile("hostfiles/adservers.txt"); err != nil {
		return err
	}

	if err := gravity.ParseAdFile("hostfiles/facebook.txt"); err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: udpPort})
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Infof("Gravity DNS server listening on port %d", udpPort)
	for {
		buf := make([]byte, packetLen)
		_, addr, _ := conn.ReadFromUDP(buf)

		var m dnsmessage.Message
		if err := m.Unpack(buf); err != nil {
			return err
		}

		if len(m.Questions) == 0 {
			log.Error("Invalid question length of 0")
			continue
		}

		go func() {
			m.Response = true
			for _, question := range m.Questions {
				log.Infof("%v requesting %s type %d", addr, question.Name.String(), question.Type)

				query := question.Name.String()
				t := dns.DNSTypeToGravity(question.Type)
				resolved, err := gravity.RetrieveAndSet(t, query)
				if err != nil {
					log.Error(err)
					continue
				}

				for _, val := range resolved {
					dnsEntry, err := dns.GravityEntryToResourceBody(t, val)
					if err != nil {
						log.Error(err)
						continue
					}

					m.Answers = append(m.Answers, dnsmessage.Resource{
						Header: dnsmessage.ResourceHeader{
							Name:  dnsmessage.MustNewName(query),
							Type:  question.Type,
							Class: question.Class,
							TTL:   100,
						},
						Body: dnsEntry,
					})
				}
			}

			packed, err := m.Pack()
			if err != nil {
				log.Error(err)
				// return err
			}
			if _, err := conn.WriteToUDP(packed, addr); err != nil {
				log.Error(err)
				// return err
			}
		}()

	}
}
