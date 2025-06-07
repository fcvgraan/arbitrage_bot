package valr

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const marketPrefix = "/api/v1/marketdata"

// ValrClient wraps a Resty client for VALR’s REST API.
type ValrClient struct {
	client *resty.Client
}

// NewValrClient returns a ValrClient.
// restyClient must have HostURL=https://api.valr.com and BasicAuth already set.
func NewValrClient(restyClient *resty.Client) *ValrClient {
	return &ValrClient{client: restyClient}
}

// GetOrderBook fetches the BTCZAR orderbook.
// If full==false you get the top 40 bids/asks only.
// GET /api/v1/marketdata/{pair}/orderbook?full={true|false}
func (c *ValrClient) GetOrderBook(pair string, full bool) (*OrderBookResponse, error) {
	resp := new(OrderBookResponse)
	req := c.client.R().SetResult(resp)
	req = req.SetQueryParam("full", fmt.Sprintf("%t", full))
	_, err := req.Get(fmt.Sprintf("%s/%s/orderbook", marketPrefix, pair))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
