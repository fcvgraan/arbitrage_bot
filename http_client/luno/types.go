package luno

// OrderBookEntry represents one level in the Luno order book.
type OrderBookEntry struct {
	Price  string `json:"price"`  // e.g. "950000.00"
	Volume string `json:"volume"` // e.g. "0.0123"
}

// OrderBookResponse is the raw JSON payload from GET /api/1/orderbook.
type OrderBookResponse struct {
	Pair      string           `json:"pair"`      // e.g. "XBTZAR"
	Timestamp int64            `json:"timestamp"` // server time in ms
	Sequence  int64            `json:"sequence"`  // incremental sequence
	Bids      []OrderBookEntry `json:"bids"`      // sorted desc by price
	Asks      []OrderBookEntry `json:"asks"`      // sorted asc by price
}
