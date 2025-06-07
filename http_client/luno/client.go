package luno

import "github.com/go-resty/resty/v2"

// LunoClient wraps Resty for Luno's API.
type LunoClient struct{ client *resty.Client }

// NewLunoClient creates a LunoClient (resty.BaseURL/auth already set).
func NewLunoClient(c *resty.Client) *LunoClient { return &LunoClient{client: c} }

// GetOrderBook fetches top-of-book for a pair.
func (c *LunoClient) GetOrderBook(pair string) (*OrderBookResponse, error) {
	r := new(OrderBookResponse)
	_, err := c.client.R().
		SetQueryParam("pair", pair).
		SetResult(r).
		Get("/api/1/orderbook")
	return r, err
}
