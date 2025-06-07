package luno

// OrderBookResponse is Luno's raw JSON for the BTC/ZAR orderbook.
type OrderBookResponse struct {
	Pair      string     `json:"pair"`
	Timestamp int64      `json:"timestamp"`
	Bids      [][]string `json:"bids"` // each entry: [price, volume]
	Asks      [][]string `json:"asks"` // each entry: [price, volume]
}
