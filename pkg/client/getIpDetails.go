package client

import (
	"context"
	"net/http"
	"strings"
)

const (
	getIpDetailsEndpoint = "/api/v1/ip/details"
)

type IpDetails struct {
	Ip           string  `json:"ip"`
	CountryCode  string  `json:"country_code"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	MysteryValue int     `json:"mystery_value"`
}

//GetIpDetails Return details about ip address.
func (c *Client) GetIpDetails(ctx context.Context, ip string) (*IpDetails, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+getIpDetailsEndpoint, strings.NewReader("{\"ip\":\""+ip+"\"}"))
	if err != nil {
		return nil, err
	}

	response := new(IpDetails)

	err = c.call(req, response)

	return response, err
}
