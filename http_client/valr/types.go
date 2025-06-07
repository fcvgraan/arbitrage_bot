package valr

// OrderBookEntry is one price‐level in VALR's book.
type OrderBookEntry struct {
	Price    string `json:"price"`    // e.g. "590000.00"
	Quantity string `json:"quantity"` // e.g. "0.05"
}

// OrderBookResponse is VALR's JSON for BTCZAR.
// GET https://api.valr.com/api/v1/marketdata/BTCZAR/orderbook?full={true|false}
type OrderBookResponse struct {
	Pair      string           `json:"pair"`
	Timestamp int64            `json:"timestamp"`
	Bids      []OrderBookEntry `json:"bids"` // sorted descending
	Asks      []OrderBookEntry `json:"asks"` // sorted ascending
}
