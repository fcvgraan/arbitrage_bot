package valr

// OrderBookEntry is one price level in VALR's book.
type OrderBookEntry struct {
	Price    string `json:"price"`    // e.g. "951500.00"
	Quantity string `json:"quantity"` // e.g. "0.0345"
}

// OrderBookResponse mirrors GET /api/v1/public/{pair}/orderbook
type OrderBookResponse struct {
	Pair      string           `json:"pair"`      // e.g. "BTCZAR"
	Timestamp int64            `json:"timestamp"` // server time in ms
	Bids      []OrderBookEntry `json:"bids"`      // sorted desc by price
	Asks      []OrderBookEntry `json:"asks"`      // sorted asc by price
}
