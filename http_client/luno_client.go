package http_client

import (
	"strconv"

	"github.com/go-resty/resty/v2"
)

// LunoClient wraps a Resty client and exposes all Luno endpoints.
type LunoClient struct {
	baseURL string
	client  *resty.Client
}

// NewLunoClient returns a LunoClient configured to talk to Luno's REST API.
func NewLunoClient(restyClient *resty.Client) *LunoClient {
	return &LunoClient{
		baseURL: "https://api.luno.com",
		client:  restyClient,
	}
}

// —————————————————————————————————————
// 1. PUBLIC / MARKET DATA
// —————————————————————————————————————

// GetTicker returns the latest ticker for a given pair.
// GET /api/1/ticker?pair=XBTZAR
func (c *LunoClient) GetTicker(pair string) (*TickerResponse, error) {
	resp := new(TickerResponse)
	_, err := c.client.R().
		SetQueryParam("pair", pair).
		SetResult(resp).
		Get(c.baseURL + "/api/1/ticker")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAllTickers returns tickers for all active pairs.
// GET /api/1/tickers
func (c *LunoClient) GetAllTickers() (*AllTickersResponse, error) {
	resp := new(AllTickersResponse)
	_, err := c.client.R().
		SetResult(resp).
		Get(c.baseURL + "/api/1/tickers")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetOrderBook returns the order book for the given pair.
// GET /api/1/orderbook?pair=XBTZAR
func (c *LunoClient) GetOrderBook(pair string) (*OrderBookResponse, error) {
	resp := new(OrderBookResponse)
	_, err := c.client.R().
		SetQueryParam("pair", pair).
		SetResult(resp).
		Get(c.baseURL + "/api/1/orderbook")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTrades returns up to 100 of the most recent trades for pair since 'since' timestamp.
// GET /api/1/trades?pair=XBTZAR&since=TIMESTAMP
func (c *LunoClient) GetTrades(pair string, since int64) (*TradesResponse, error) {
	resp := new(TradesResponse)
	req := c.client.R().SetResult(resp)
	if pair != "" {
		req = req.SetQueryParam("pair", pair)
	}
	if since > 0 {
		req = req.SetQueryParam("since", strconv.FormatInt(since, 10))
	}
	_, err := req.Get(c.baseURL + "/api/1/trades")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// —————————————————————————————————————
// 2. PRIVATE / ACCOUNT & TRANSACTIONS
// —————————————————————————————————————

// CreateAccount creates a new account (wallet) for the specified currency.
// POST /api/1/create_account
func (c *LunoClient) CreateAccount(name, currency string) (*CreateAccountResponse, error) {
	payload := CreateAccountRequest{Name: name, Currency: currency}
	resp := new(CreateAccountResponse)
	_, err := c.client.R().
		SetBody(payload).
		SetResult(resp).
		Post(c.baseURL + "/api/1/create_account")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetBalances returns all your account balances.
// GET /api/1/balance
func (c *LunoClient) GetBalances() ([]Account, error) {
	resp := new(BalancesResponse)
	_, err := c.client.R().
		SetResult(resp).
		Get(c.baseURL + "/api/1/balance")
	if err != nil {
		return nil, err
	}
	return resp.Balance, nil
}

// GetTransactions returns recent transactions for a given account ID.
// GET /api/1/list_transactions?id=ACCOUNT_ID&min_row=&max_row=
func (c *LunoClient) GetTransactions(accountID string, minRow, maxRow int) ([]TransactionEntry, error) {
	resp := new(TransactionsResponse)
	req := c.client.R().SetResult(resp).
		SetQueryParam("id", accountID)
	if minRow != 0 {
		req = req.SetQueryParam("min_row", strconv.Itoa(minRow))
	}
	if maxRow != 0 {
		req = req.SetQueryParam("max_row", strconv.Itoa(maxRow))
	}
	_, err := req.Get(c.baseURL + "/api/1/list_transactions")
	if err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// —————————————————————————————————————
// 3. TRADING
// —————————————————————————————————————

// PostMarketOrder creates a new market order. (BUY spends ZAR, SELL spends BTC).
// POST /api/1/marketorder
func (c *LunoClient) PostMarketOrder(orderType, volume, pair string) (*PostMarketOrderResponse, error) {
	payload := PostMarketOrderRequest{
		Type:   orderType,
		Volume: volume,
		Pair:   pair,
	}
	resp := new(PostMarketOrderResponse)
	_, err := c.client.R().
		SetBody(payload).
		SetResult(resp).
		Post(c.baseURL + "/api/1/marketorder")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// StopOrder cancels an existing order by its ID.
// POST /api/1/stoporder
func (c *LunoClient) StopOrder(orderID string) error {
	payload := StopOrderRequest{ID: orderID}
	_, err := c.client.R().
		SetBody(payload).
		Post(c.baseURL + "/api/1/stoporder")
	return err
}

// GetOrder fetches details of a single order.
// GET /api/1/getorder?id=ORDER_ID
func (c *LunoClient) GetOrder(orderID string) (*GetOrderResponse, error) {
	resp := new(GetOrderResponse)
	_, err := c.client.R().
		SetQueryParam("id", orderID).
		SetResult(resp).
		Get(c.baseURL + "/api/1/getorder")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTradesList returns your past trades for a pair (oldest first).
// GET /api/1/list_trades?pair=XBTZAR&since=TIMESTAMP&limit=NUM
func (c *LunoClient) GetTradesList(pair string, since int64, limit int) (*TradesListResponse, error) {
	resp := new(TradesListResponse)
	req := c.client.R().SetResult(resp)
	if pair != "" {
		req = req.SetQueryParam("pair", pair)
	}
	if since > 0 {
		req = req.SetQueryParam("since", strconv.FormatInt(since, 10))
	}
	if limit > 0 {
		req = req.SetQueryParam("limit", strconv.Itoa(limit))
	}
	_, err := req.Get(c.baseURL + "/api/1/list_trades")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// —————————————————————————————————————
// 4. WALLET & WITHDRAWALS
// —————————————————————————————————————

// GetReceiveAddress returns your default (or a specified) receive address for an asset.
// GET /api/1/receive_address?asset=XBT&address=OPTIONAL
func (c *LunoClient) GetReceiveAddress(asset, address string) (*ReceiveAddressResponse, error) {
	resp := new(ReceiveAddressResponse)
	req := c.client.R().SetResult(resp).
		SetQueryParam("asset", asset)
	if address != "" {
		req = req.SetQueryParam("address", address)
	}
	_, err := req.Get(c.baseURL + "/api/1/receive_address")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateReceiveAddress allocates a new receive address for an asset.
// POST /api/1/create_receive_address
func (c *LunoClient) CreateReceiveAddress(asset string) (*CreateReceiveAddressResponse, error) {
	payload := map[string]string{"asset": asset}
	resp := new(CreateReceiveAddressResponse)
	_, err := c.client.R().
		SetBody(payload).
		SetResult(resp).
		Post(c.baseURL + "/api/1/create_receive_address")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Send transfers crypto from your account to another address/email.
// POST /api/1/send
func (c *LunoClient) Send(amount, currency, address, description, message string) (*SendResponse, error) {
	payload := SendRequest{
		Amount:      amount,
		Currency:    currency,
		Address:     address,
		Description: description,
		Message:     message,
	}
	resp := new(SendResponse)
	_, err := c.client.R().
		SetBody(payload).
		SetResult(resp).
		Post(c.baseURL + "/api/1/send")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetWithdrawalRequests returns your past withdrawal requests.
// GET /api/1/list_withdrawals
func (c *LunoClient) GetWithdrawalRequests() ([]WithdrawalRequest, error) {
	resp := new(WithdrawalRequestsResponse)
	_, err := c.client.R().
		SetResult(resp).
		Get(c.baseURL + "/api/1/list_withdrawals")
	if err != nil {
		return nil, err
	}
	return resp.Requests, nil
}

// RequestWithdrawal creates a new withdrawal request.
// POST /api/1/withdrawal
func (c *LunoClient) RequestWithdrawal(wtype, amount, beneficiaryID string) (*RequestWithdrawalResponse, error) {
	payload := RequestWithdrawalRequest{
		Type:          wtype,
		Amount:        amount,
		BeneficiaryID: beneficiaryID,
	}
	resp := new(RequestWithdrawalResponse)
	_, err := c.client.R().
		SetBody(payload).
		SetResult(resp).
		Post(c.baseURL + "/api/1/withdrawal")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetWithdrawalStatus fetches the status of a withdrawal by its ID.
// GET /api/1/get_withdrawal?id=ID
func (c *LunoClient) GetWithdrawalStatus(id string) (*WithdrawalStatusResponse, error) {
	resp := new(WithdrawalStatusResponse)
	_, err := c.client.R().
		SetQueryParam("id", id).
		SetResult(resp).
		Get(c.baseURL + "/api/1/get_withdrawal")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CancelWithdrawalRequest cancels a pending withdrawal.
// DELETE /api/1/cancel_withdrawal?id=ID
func (c *LunoClient) CancelWithdrawalRequest(id string) error {
	_, err := c.client.R().
		SetQueryParam("id", id).
		Delete(c.baseURL + "/api/1/cancel_withdrawal")
	return err
}

// —————————————————————————————————————
// 5. QUOTES
// —————————————————————————————————————

// CreateQuote creates a new quote to BUY/SELL an amount.
// POST /api/1/quotes
func (c *LunoClient) CreateQuote(qType, amount, pair string) (*CreateQuoteResponse, error) {
	payload := CreateQuoteRequest{Type: qType, Amount: amount, Pair: pair}
	resp := new(CreateQuoteResponse)
	_, err := c.client.R().
		SetBody(payload).
		SetResult(resp).
		Post(c.baseURL + "/api/1/quotes")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetQuote fetches details of an existing quote.
// GET /api/1/get_quote?id=QUOTE_ID
func (c *LunoClient) GetQuote(id string) (*GetQuoteResponse, error) {
	resp := new(GetQuoteResponse)
	_, err := c.client.R().
		SetQueryParam("id", id).
		SetResult(resp).
		Get(c.baseURL + "/api/1/get_quote")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ExerciseQuote exercises a quote (executes the trade).
// POST /api/1/exercise_quote
func (c *LunoClient) ExerciseQuote(id string) (*GetQuoteResponse, error) {
	payload := map[string]string{"id": id}
	resp := new(GetQuoteResponse)
	_, err := c.client.R().
		SetBody(payload).
		SetResult(resp).
		Post(c.baseURL + "/api/1/exercise_quote")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DiscardQuote discards a quote so it cannot be exercised.
// POST /api/1/discard_quote
func (c *LunoClient) DiscardQuote(id string) (*GetQuoteResponse, error) {
	payload := map[string]string{"id": id}
	resp := new(GetQuoteResponse)
	_, err := c.client.R().
		SetBody(payload).
		SetResult(resp).
		Post(c.baseURL + "/api/1/discard_quote")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// —————————————————————————————————————
// 6. FEES
// —————————————————————————————————————

// GetFeeInfo returns your maker/taker fees and 30d volume for a pair.
// GET /api/1/fee_info?pair=XBTZAR
func (c *LunoClient) GetFeeInfo(pair string) (*FeeInfoResponse, error) {
	resp := new(FeeInfoResponse)
	_, err := c.client.R().
		SetQueryParam("pair", pair).
		SetResult(resp).
		Get(c.baseURL + "/api/1/fee_info")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
