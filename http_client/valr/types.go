package valr

// OrderBookEntry represents one level in VALR's book.
type OrderBookEntry struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

// OrderBookResponse holds VALR's JSON book.
// Bids sorted desc, Asks sorted asc.
type OrderBookResponse struct {
	Pair      string           `json:"pair"`
	Timestamp int64            `json:"timestamp"`
	Bids      []OrderBookEntry `json:"bids"`
	Asks      []OrderBookEntry `json:"asks"`
}
