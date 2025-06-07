package luno

// OrderBookResponse represents Luno's order book JSON.
// Bids/Asks are slices of [price, volume].
type OrderBookResponse struct {
	Pair      string     `json:"pair"`
	Timestamp int64      `json:"timestamp"`
	Bids      [][]string `json:"bids"` // [price, volume]
	Asks      [][]string `json:"asks"` // [price, volume]
}
