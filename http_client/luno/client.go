package luno

import (
	"github.com/go-resty/resty/v2"
)

// LunoClient wraps a Resty client for Luno’s REST API.
type LunoClient struct {
	client *resty.Client
}

// NewLunoClient returns a LunoClient.
// restyClient must have BaseURL=https://api.luno.com and BasicAuth already set.
func NewLunoClient(restyClient *resty.Client) *LunoClient {
	return &LunoClient{client: restyClient}
}

// GetOrderBook fetches the top‐of‐book for the given pair (e.g. "XBTZAR").
// GET https://api.luno.com/api/1/orderbook?pair=XBTZAR
func (c *LunoClient) GetOrderBook(pair string) (*OrderBookResponse, error) {
	resp := new(OrderBookResponse)
	_, err := c.client.R().
		SetQueryParam("pair", pair).
		SetResult(resp).
		Get("/api/1/orderbook")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
