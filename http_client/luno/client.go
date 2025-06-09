package luno

import (
	"github.com/go-resty/resty/v2"
)

// LunoClient wraps a Resty client pre-configured for Luno.
type LunoClient struct {
	client *resty.Client
}

// NewLunoClient expects a Resty client with BaseURL set to
// "https://api.luno.com" and BasicAuth(key, secret) already applied.
func NewLunoClient(rc *resty.Client) *LunoClient {
	return &LunoClient{client: rc}
}

// GetOrderBook fetches the best bids/asks for the given pair.
// GET /api/1/orderbook?pair=XBTZAR
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
