package resolver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	dohPath = "/dns-query?"
)

type (
	DoHQuestion struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	}

	DoHAnswer struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int
		Data string `json:"data"`
	}

	DoHResponse struct {
		Status int

		TC bool
		RD bool
		RA bool
		AD bool
		CD bool

		Question []DoHQuestion
		Answer   []DoHAnswer
	}
)

func doHTTPRequest(domain, queryType string) (*http.Response, error) {
	client := &http.Client{}

	params := url.Values{}
	params.Add("name", domain)
	params.Add("type", queryType)

	req, err := http.NewRequest("GET", "https://1.1.1.1"+dohPath+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/dns-json")
	return client.Do(req)
}

func ResolveOverDoH(domain, queryType string) (*DoHResponse, error) {
	resp, err := doHTTPRequest(domain, queryType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	dohResponse := &DoHResponse{}
	if err = json.Unmarshal(body, dohResponse); err != nil {
		return nil, err
	}
	return dohResponse, nil
}
