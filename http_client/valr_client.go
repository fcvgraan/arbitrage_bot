package http_client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	publicPrefix  = "/api/v1/public"
	marketPrefix  = "/api/v1/marketdata"
	accountPrefix = "/api/v1/account"
	orderPrefix   = "/api/v1/orders"
)

// ValrClient wraps a Resty client and exposes all VALR endpoints.
type ValrClient struct {
	client *resty.Client
}

// NewValrClient returns a ValrClient configured to talk to VALR’s REST API.
func NewValrClient(restyClient *resty.Client) *ValrClient {
	return &ValrClient{client: restyClient}
}

// —————————————————————————————————————
// 1. PUBLIC / MARKET DATA
// —————————————————————————————————————

// GetCurrencyPairs lists all active (and inactive) pairs.
// GET /api/v1/public/currencyPairs
func (c *ValrClient) GetCurrencyPairs() (*CurrencyPairsResponse, error) {
	resp := new(CurrencyPairsResponse)
	_, err := c.client.R().
		SetResult(resp).
		Get(publicPrefix + "/currencyPairs")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetOrderBook returns the top bids/asks for a pair.
// GET /api/v1/marketdata/{pair}/orderbook?full={full}
func (c *ValrClient) GetOrderBook(pair string, full bool) (*OrderBookResponse, error) {
	resp := new(OrderBookResponse)
	req := c.client.R().SetResult(resp)
	if full {
		req = req.SetQueryParam("full", "true")
	} else {
		req = req.SetQueryParam("full", "false")
	}
	_, err := req.Get(fmt.Sprintf("%s/%s/orderbook", marketPrefix, pair))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMarketSummary returns 24h summary for a given pair.
// GET /api/v1/marketdata/{pair}/ticker
func (c *ValrClient) GetMarketSummary(pair string) (*MarketSummary, error) {
	resp := new(MarketSummary)
	_, err := c.client.R().
		SetResult(resp).
		Get(fmt.Sprintf("%s/%s/ticker", marketPrefix, pair))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetServerTime returns VALR server time.
// GET /api/v1/marketdata/serverTime
func (c *ValrClient) GetServerTime() (*ServerTimeResponse, error) {
	resp := new(ServerTimeResponse)
	_, err := c.client.R().
		SetResult(resp).
		Get(marketPrefix + "/serverTime")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// —————————————————————————————————————
// 2. PRIVATE / ACCOUNT
// —————————————————————————————————————

// GetBalances returns all your wallet balances (ZAR, BTC, ETH, etc.).
// GET /api/v1/account/balances
func (c *ValrClient) GetBalances() ([]Balance, error) {
	resp := new(ValrBalancesResponse)
	_, err := c.client.R().
		SetResult(resp).
		Get(accountPrefix + "/balances")
	if err != nil {
		return nil, err
	}
	return resp.Balances, nil
}

// GetTradeHistory returns your last N trades on a pair.
// GET /api/v1/account/{pair}/tradeHistory?limit={limit}
func (c *ValrClient) GetTradeHistory(pair string, limit int) ([]TradeHistoryEntry, error) {
	resp := new(TradeHistoryResponse)
	req := c.client.R().SetResult(resp)
	req = req.SetPathParam("pair", pair).
		SetQueryParam("limit", fmt.Sprintf("%d", limit))
	_, err := req.Get(fmt.Sprintf("%s/%s/tradeHistory", accountPrefix, pair))
	if err != nil {
		return nil, err
	}
	return resp.Trades, nil
}

// —————————————————————————————————————
// 3. TRADING
// —————————————————————————————————————

// PlaceLimitOrder submits a LIMIT order.
// POST /api/v1/orders/limit
func (c *ValrClient) PlaceLimitOrder(reqBody PlaceLimitOrderRequest) (*PlaceLimitOrderResponse, error) {
	resp := new(PlaceLimitOrderResponse)
	_, err := c.client.R().
		SetBody(reqBody).
		SetResult(resp).
		Post(orderPrefix + "/limit")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PlaceMarketOrder submits a MARKET order.
// POST /api/v1/orders/market
func (c *ValrClient) PlaceMarketOrder(reqBody PlaceMarketOrderRequest) (*PlaceMarketOrderResponse, error) {
	resp := new(PlaceMarketOrderResponse)
	_, err := c.client.R().
		SetBody(reqBody).
		SetResult(resp).
		Post(orderPrefix + "/market")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetOrderStatus fetches details of a single order by ID.
// GET /api/v1/orders/{orderId}
func (c *ValrClient) GetOrderStatus(orderID string) (*OrderStatusResponse, error) {
	resp := new(OrderStatusResponse)
	_, err := c.client.R().
		SetResult(resp).
		Get(fmt.Sprintf("%s/%s", orderPrefix, orderID))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CancelOrder cancels an existing order by ID.
// DELETE /api/v1/orders/{orderId}
func (c *ValrClient) CancelOrder(orderID string) error {
	_, err := c.client.R().
		Delete(fmt.Sprintf("%s/%s", orderPrefix, orderID))
	return err
}
