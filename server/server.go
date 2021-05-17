package server

import (
	"net"

	"github.com/gravity-dns/gravity-dns/dns"
	"github.com/gravity-dns/gravity-dns/resolver"
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
	log.Info("Gravity DNS server started")
	for {
		buf := make([]byte, packetLen)
		_, addr, _ := conn.ReadFromUDP(buf)
		var m dnsmessage.Message
		if err := m.Unpack(buf); err != nil {
			return err
		}

		query := m.Questions[0].Name.String()
		log.Info(m)
		log.Info("addr %v requested %s", addr, query)

		resolved := gravity.Resolve(query)
		if resolved != "" {
			log.Info("Gravity resolved ", query, " - ", resolved)

			ip := net.ParseIP(resolved).To4()
			m.Response = true
			m.Answers = append(m.Answers, dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  dnsmessage.MustNewName(query),
					Type:  m.Questions[0].Type,
					Class: m.Questions[0].Class,
					TTL:   uint32(10000),
				},
				Body: &dnsmessage.AResource{A: [4]byte{ip[0], ip[1], ip[2], ip[3]}},
			})
			packed, err := m.Pack()
			if err != nil {
				return err
			}
			if _, err := conn.WriteToUDP(packed, addr); err != nil {
				return err
			}
			continue
		}

		doh, err := resolver.ResolveOverDoH(m.Questions[0].Name.String(), "A")
		if err != nil {
			return err
		}

		var ip net.IP
		for _, ans := range doh.Answer {
			ip = net.ParseIP(ans.Data).To4()
			m.Response = doh.Status == 0
			m.Answers = append(m.Answers, dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  dnsmessage.MustNewName(ans.Name + "."),
					Type:  m.Questions[0].Type,
					Class: m.Questions[0].Class,
					TTL:   uint32(ans.TTL),
				},
				Body: &dnsmessage.AResource{A: [4]byte{ip[0], ip[1], ip[2], ip[3]}},
			})
		}
		packed, err := m.Pack()
		if err != nil {
			return err
		}
		if _, err := conn.WriteToUDP(packed, addr); err != nil {
			return err
		}
		gravity.Set(query, ip.String())
	}
}
