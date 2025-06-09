package valr

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// OrderBookResponse and OrderBookEntry in valr/types.go remain unchanged.

// ValrClient wraps a Resty client pre-configured for VALR.
type ValrClient struct {
	client *resty.Client
}

// NewValrClient expects a Resty client with BaseURL="https://api.valr.com" and BasicAuth set.
func NewValrClient(rc *resty.Client) *ValrClient {
	return &ValrClient{client: rc}
}

// GetOrderBook fetches the top 40 bids & asks via the public endpoint.
// GET https://api.valr.com/v1/public/{pair}/orderbook
func (c *ValrClient) GetOrderBook(pair string) (*OrderBookResponse, error) {
	resp := new(OrderBookResponse)
	res, err := c.client.R().
		SetResult(resp).
		Get(fmt.Sprintf("/v1/public/%s/orderbook", pair))
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("VALR orderbook HTTP %d", res.StatusCode())
	}
	return resp, nil
}
