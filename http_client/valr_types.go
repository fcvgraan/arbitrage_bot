package http_client

// ——————————————————————————————————————————————
// 1. PUBLIC / MARKET DATA TYPES
// ——————————————————————————————————————————————

type CurrencyPair struct {
	Pair          string `json:"pair"`          // e.g. "BTCZAR"
	BaseCurrency  string `json:"baseCurrency"`  // e.g. "BTC"
	QuoteCurrency string `json:"quoteCurrency"` // e.g. "ZAR"
	Status        string `json:"status"`        // "ACTIVE" or "INACTIVE"
}

type CurrencyPairsResponse struct {
	Pairs []CurrencyPair `json:"pairs"`
}

type ValrOrderBookEntry struct {
	Price    string `json:"price"`    // e.g. "590000.00"
	Quantity string `json:"quantity"` // e.g. "0.05"
}

type ValrOrderBookResponse struct {
	Pair      string           `json:"pair"`
	Timestamp int64            `json:"timestamp"`
	Bids      []OrderBookEntry `json:"bids"` // sorted descending
	Asks      []OrderBookEntry `json:"asks"` // sorted ascending
}

type MarketSummary struct {
	Pair            string `json:"pair"`
	LastTradedPrice string `json:"lastTradedPrice"`
	LowestAskPrice  string `json:"lowestAskPrice"`
	HighestBidPrice string `json:"highestBidPrice"`
	DayLow          string `json:"dayLow"`
	DayHigh         string `json:"dayHigh"`
	Volume          string `json:"volume"`
	QuoteVolume24Hr string `json:"baseVolume"` // sometimes called “baseVolume”
}

type MarketSummaryResponse struct {
	Summaries []MarketSummary `json:"marketSummary"`
}

type ServerTimeResponse struct {
	ServerTime int64 `json:"time"`
}

// ——————————————————————————————————————————————
// 2. PRIVATE / ACCOUNT TYPES
// ——————————————————————————————————————————————

type Balance struct {
	Currency         string `json:"currency"`         // e.g. "BTC", "ZAR"
	AvailableBalance string `json:"availableBalance"` // e.g. "0.005"
	TotalBalance     string `json:"totalBalance"`     // e.g. "0.010"
}

type ValrBalancesResponse struct {
	Balances []Balance `json:"balances"`
}

type OrderStatusResponse struct {
	OrderID         string `json:"orderId"`
	CustomerOrderID string `json:"customerOrderId,omitempty"`
	CurrencyPair    string `json:"currencyPair"`
	Side            string `json:"side"`     // "BUY" or "SELL"
	Type            string `json:"type"`     // "LIMIT", "MARKET", etc.
	Price           string `json:"price"`    // price per coin (for LIMIT)
	Quantity        string `json:"quantity"` // amount in base (e.g. BTC)
	QuantityFilled  string `json:"quantityFilled"`
	Status          string `json:"status"` // "NEW", "PARTIALLY_FILLED", "FILLED", "CANCELLED"
	Timestamp       int64  `json:"timestamp"`
}

type PlaceLimitOrderRequest struct {
	CurrencyPair    string `json:"pair"`            // e.g. "BTCZAR"
	Side            string `json:"side"`            // "BUY" or "SELL"
	Price           string `json:"price"`           // per coin, in quote currency
	Quantity        string `json:"quantity"`        // in base currency
	PostOnly        bool   `json:"postOnly"`        // optional; default false
	CustomerOrderID string `json:"customerOrderId"` // optional; unique across your orders
}

type PlaceLimitOrderResponse struct {
	OrderID         string `json:"orderId"`
	CustomerOrderID string `json:"customerOrderId,omitempty"`
	CurrencyPair    string `json:"currencyPair"`
	Side            string `json:"side"`
	Type            string `json:"type"`
	Price           string `json:"price"`
	Quantity        string `json:"quantity"`
	Timestamp       int64  `json:"timestamp"`
	Status          string `json:"status"`
}

type PlaceMarketOrderRequest struct {
	CurrencyPair    string `json:"pair"`            // e.g. "BTCZAR"
	Side            string `json:"side"`            // "BUY" or "SELL"
	Amount          string `json:"amount"`          // if SIDE=BUY, amount in quote (ZAR); if SIDE=SELL, amount in base (BTC)
	CustomerOrderID string `json:"customerOrderId"` // optional
}

type PlaceMarketOrderResponse struct {
	OrderID         string `json:"orderId"`
	CustomerOrderID string `json:"customerOrderId,omitempty"`
	CurrencyPair    string `json:"currencyPair"`
	Side            string `json:"side"`
	Type            string `json:"type"`
	Price           string `json:"price"`    // average executed price
	Quantity        string `json:"quantity"` // filled qty
	Timestamp       int64  `json:"timestamp"`
	Status          string `json:"status"`
}

type TradeHistoryEntry struct {
	OrderID         string `json:"orderId"`
	CustomerOrderID string `json:"customerOrderId,omitempty"`
	Pair            string `json:"pair"`
	Side            string `json:"side"`
	Price           string `json:"price"`
	Quantity        string `json:"quantity"`
	Timestamp       int64  `json:"timestamp"`
}

type TradeHistoryResponse struct {
	Trades []TradeHistoryEntry `json:"trades"`
}
