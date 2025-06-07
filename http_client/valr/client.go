package valr

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const apiPrefix = "/api/v1/marketdata"

// ValrClient wraps Resty for VALR's API.
type ValrClient struct{ client *resty.Client }

// NewValrClient initializes ValrClient (resty.BaseURL/auth already set).
func NewValrClient(c *resty.Client) *ValrClient { return &ValrClient{client: c} }

// GetOrderBook fetches top-of-book; full=false for top 40 entries.
func (c *ValrClient) GetOrderBook(pair string, full bool) (*OrderBookResponse, error) {
	r := new(OrderBookResponse)
	req := c.client.R().SetResult(r).
		SetQueryParam("full", fmt.Sprintf("%t", full))
	_, err := req.Get(fmt.Sprintf("%s/%s/orderbook", apiPrefix, pair))
	return r, err
}
