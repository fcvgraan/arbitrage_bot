package http_client

// ——————————————————————————————————————————————
// PUBLIC / MARKET DATA ENDPOINTS
// ——————————————————————————————————————————————

type TickerResponse struct {
	Pair      string `json:"pair"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	Volume    string `json:"volume"`
	LastTrade string `json:"last_trade"`
	Timestamp int64  `json:"timestamp"`
}

type AllTickersResponse struct {
	Tickers []struct {
		Pair      string `json:"pair"`
		Bid       string `json:"bid"`
		Ask       string `json:"ask"`
		Volume    string `json:"volume"`
		LastTrade string `json:"last_trade"`
		Timestamp int64  `json:"timestamp"`
	} `json:"tickers"`
}

type OrderBookEntry struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
}

type OrderBookResponse struct {
	Pair   string         `json:"pair"`
	Tid    int64          `json:"timestamp"` // Note: Luno calls this `timestamp`
	Bids   [][]string     `json:"bids"`      // each entry: [price, volume]
	Asks   [][]string     `json:"asks"`
	TopBid OrderBookEntry `json:"-"`
	TopAsk OrderBookEntry `json:"-"`
	// We can post-process Bids/Asks into TopBid/TopAsk, if desired.
}

type TradeEntry struct {
	Timestamp int64  `json:"timestamp"`
	Price     string `json:"price"`
	Volume    string `json:"volume"`
	Side      string `json:"side"` // BUY or SELL
}

type TradesResponse struct {
	Trades []TradeEntry `json:"trades"`
}

// ——————————————————————————————————————————————
// PRIVATE / ACCOUNT & TRADING ENDPOINTS
// ——————————————————————————————————————————————

type Account struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Reserved  string `json:"reserved"`
	Ghost     bool   `json:"ghost"`
}

type BalancesResponse struct {
	Balance []Account `json:"balance"`
}

type TransactionEntry struct {
	// Fields vary by transaction type; we include the most common
	//   - type: "DEPOSIT", "EXCHANGE_SELL", etc.
	//   - status: "PENDING", "COMPLETED", etc.
	//   - details: may include "amount", "fee", etc.
	Timestamp int64  `json:"timestamp"`
	Address   string `json:"address"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
	Network   string `json:"network"`
	Fee       string `json:"fee"`
	ID        int64  `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
}

type TransactionsResponse struct {
	Entries []TransactionEntry `json:"transactions"`
}

// CreateAccount: POST /api/1/create_account
type CreateAccountRequest struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
}
type CreateAccountResponse struct {
	AccountID   int64  `json:"id"`
	Currency    string `json:"currency"`
	CreatedAt   int64  `json:"created_at"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Seq         int64  `json:"sequence"`
	ReferenceID string `json:"reference_id"`
}

// ——————————————————————————————————————————————
// TRADING ENDPOINTS
// ——————————————————————————————————————————————

type PostMarketOrderRequest struct {
	Type   string `json:"type"`   // "BUY" or "SELL"
	Volume string `json:"volume"` // For BUY: amount in ZAR; for SELL: amount in BTC
	Pair   string `json:"pair"`   // e.g. "XBTZAR"
}

type PostMarketOrderResponse struct {
	OrderID   int64  `json:"order_id"`
	Type      string `json:"type"`
	Volume    string `json:"volume"`
	Pair      string `json:"pair"`
	Price     string `json:"price"`
	Tid       int64  `json:"timestamp"`
	Status    string `json:"status"`
	Fee       string `json:"fee"`
	Remaining string `json:"remaining_volume"`
	Funds     string `json:"funds"`
}

// StopOrder: POST /api/1/stop_order
type StopOrderRequest struct {
	ID string `json:"id"`
}

// GetOrder: GET /api/1/get_order
type GetOrderResponse struct {
	OrderID      int64  `json:"order_id"`
	Type         string `json:"type"`
	Volume       string `json:"volume"`
	Remaining    string `json:"remaining_volume"`
	Price        string `json:"price"`
	Pair         string `json:"pair"`
	Status       string `json:"status"`
	Timestamp    int64  `json:"timestamp"`
	FeePaid      string `json:"fee_paid"`
	Tid          int64  `json:"balance_after"`
	PostOnly     bool   `json:"post_only"`
	TriggerPrice string `json:"trigger_price,omitempty"`
}

// GetTradesList: GET /api/1/list_trades
type TradesListResponse struct {
	Trades []struct {
		ID      int64  `json:"id"`
		OrderID int64  `json:"order_id"`
		Pair    string `json:"pair"`
		Type    string `json:"type"`
		Price   string `json:"price"`
		Volume  string `json:"volume"`
		FeePaid string `json:"fee_paid"`
		Tid     int64  `json:"timestamp"`
		Side    string `json:"side"`
		Funds   string `json:"funds"`
	} `json:"trades"`
}

// ——————————————————————————————————————————————
// WALLET & WITHDRAWAL ENDPOINTS
// ——————————————————————————————————————————————

type ReceiveAddressResponse struct {
	Address       string `json:"address"`
	TotalReceived string `json:"total_received"`
	Unconfirmed   string `json:"total_unconfirmed"`
}

type CreateReceiveAddressResponse struct {
	Address       string `json:"address"`
	TotalReceived string `json:"total_received"`
	Unconfirmed   string `json:"total_unconfirmed"`
}

type SendRequest struct {
	Amount      string `json:"amount"`   // decimal string, e.g. "0.001"
	Currency    string `json:"currency"` // "XBT" or "ETH"
	Address     string `json:"address"`
	Description string `json:"description,omitempty"`
	Message     string `json:"message,omitempty"`
}

type SendResponse struct {
	TransactionID int64  `json:"id"`
	Status        string `json:"status"`
	FeePaid       string `json:"fee_paid"`
}

// ——————————————————————————————————————————————
// QUOTE ENDPOINTS
// ——————————————————————————————————————————————

type CreateQuoteRequest struct {
	Type   string `json:"type"`   // "BUY" or "SELL"
	Amount string `json:"amount"` // amount in base currency, e.g. "0.01" BTC
	Pair   string `json:"pair"`   // e.g. "XBTZAR"
}

type CreateQuoteResponse struct {
	QuoteID   int64  `json:"quote_id"`
	Type      string `json:"type"`
	Price     string `json:"price"`
	Fee       string `json:"fee"`
	Volume    string `json:"volume"`
	Pair      string `json:"pair"`
	ExpiresAt int64  `json:"expires_at"`
}

type GetQuoteResponse CreateQuoteResponse

// ExerciseQuote / DiscardQuote responses are identical to GetQuoteResponse

// ——————————————————————————————————————————————
// FEE & WITHDRAWAL REQUEST ENDPOINTS
// ——————————————————————————————————————————————

type FeeInfoResponse struct {
	TakerFee          string `json:"taker_fee"`
	MakerFee          string `json:"maker_fee"`
	DailyVolume       string `json:"daily_volume"`
	NextMakerDiscount string `json:"next_maker_discount"`
}

type WithdrawalRequest struct {
	ID         int64  `json:"id"`
	Amount     string `json:"amount"`
	Fee        string `json:"fee"`
	NetAmount  string `json:"net_amount"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	CreatedAt  int64  `json:"created_at"`
	Processing int64  `json:"processing_at"`
}

type WithdrawalRequestsResponse struct {
	Requests []WithdrawalRequest `json:"withdrawals"`
}

type RequestWithdrawalRequest struct {
	Type          string `json:"type"` // e.g. "ZAR_EFT"
	Amount        string `json:"amount"`
	BeneficiaryID string `json:"beneficiary_id,omitempty"`
}

type RequestWithdrawalResponse struct {
	ID         int64  `json:"id"`
	Amount     string `json:"amount"`
	Fee        string `json:"fee"`
	NetAmount  string `json:"net_amount"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	CreatedAt  int64  `json:"created_at"`
	Processing int64  `json:"processing_at"`
}

type WithdrawalStatusResponse RequestWithdrawalResponse
