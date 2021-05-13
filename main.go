package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/dns/dnsmessage"
)

var (
	initialDNS = net.IPv4(1, 1, 1, 1)
)

const (
	// DNS server default port
	udpPort int = 53
	// DNS packet max length
	packetLen  int    = 512
	cloudflare string = "/dns-query?"
)

type question struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}
type answer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int
	Data string `json:"data"`
}
type response struct {
	Status   int
	TC       bool
	RD       bool
	RA       bool
	AD       bool
	CD       bool
	Question []question
	Answer   []answer
}

func queryIP() net.IP {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: udpPort})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	return nil
}
func getDoH(addr string, t string) *response {
	params := url.Values{}
	params.Add("name", addr)
	params.Add("type", t)
	client := &http.Client{}
	fmt.Println("query", cloudflare+params.Encode())
	req, _ := http.NewRequest("GET", cloudflare+params.Encode(), nil)
	req.Header.Set("accept", "application/dns-json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	dohResponse := &response{}
	if err = json.Unmarshal(body, dohResponse); err != nil {
		panic(err)
	}
	return dohResponse
}
func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: udpPort})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("running dns server...")
	for {
		buf := make([]byte, packetLen)
		_, addr, _ := conn.ReadFromUDP(buf)
		var m dnsmessage.Message
		if err := m.Unpack(buf); err != nil {
			panic(err)
		}
		fmt.Println(addr)
		fmt.Println(m)
		doh := getDoH(m.Questions[0].Name.String(), "A")
		fmt.Println("verified ", doh.CD)
		for _, ans := range doh.Answer {
			ip := net.ParseIP(ans.Data).To4()
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
			panic(err)
		}
		if _, err := conn.WriteToUDP(packed, addr); err != nil {
			panic(err)
		}
	}
}
