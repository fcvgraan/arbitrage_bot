package http_client

import (
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

// These tests require real VALR API keys in environment:
// export VALR_API_KEY="your_key_id"
// export VALR_API_SECRET="your_key_secret"
func TestGetCurrencyPairs(t *testing.T) {
	apiKey := os.Getenv("VALR_API_KEY")
	apiSecret := os.Getenv("VALR_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping VALR tests; VALR_API_KEY/VALR_API_SECRET not set")
	}

	restyClient := NewRestyClientForTests(apiKey, apiSecret)
	valr := NewValrClient(restyClient)

	pairsResp, err := valr.GetCurrencyPairs()
	if err != nil {
		t.Fatalf("GetCurrencyPairs failed: %v", err)
	}
	if len(pairsResp.Pairs) == 0 {
		t.Errorf("Expected at least one currency pair, got zero")
	}
	found := false
	for _, p := range pairsResp.Pairs {
		if p.Pair == "BTCZAR" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to see BTCZAR in currency pairs, but did not")
	}
}

func TestGetOrderBook(t *testing.T) {
	apiKey := os.Getenv("VALR_API_KEY")
	apiSecret := os.Getenv("VALR_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping VALR tests; VALR_API_KEY/VALR_API_SECRET not set")
	}

	restyClient := NewRestyClientForTests(apiKey, apiSecret)
	valr := NewValrClient(restyClient)

	ob, err := valr.GetOrderBook("BTCZAR", false)
	if err != nil {
		t.Fatalf("GetOrderBook failed: %v", err)
	}
	if len(ob.Bids) == 0 || len(ob.Asks) == 0 {
		t.Errorf("Expected non-empty bids/asks, got bids=%d, asks=%d", len(ob.Bids), len(ob.Asks))
	}
	t.Logf("Top Bid: %s (vol=%s), Top Ask: %s (vol=%s)", ob.Bids[0].Price, ob.Bids[0].Quantity, ob.Asks[0].Price, ob.Asks[0].Quantity)
}

func TestGetBalances(t *testing.T) {
	apiKey := os.Getenv("VALR_API_KEY")
	apiSecret := os.Getenv("VALR_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping VALR tests; VALR_API_KEY/VALR_API_SECRET not set")
	}

	restyClient := NewRestyClientForTests(apiKey, apiSecret)
	valr := NewValrClient(restyClient)

	bals, err := valr.GetBalances()
	if err != nil {
		t.Fatalf("GetBalances failed: %v", err)
	}
	if len(bals) == 0 {
		t.Errorf("Expected at least one balance entry, got zero")
	}
	for _, b := range bals {
		if b.Currency == "" {
			t.Errorf("Found balance with empty currency")
		}
	}
}

func TestPlaceAndCancelLimitOrder(t *testing.T) {
	apiKey := os.Getenv("VALR_API_KEY")
	apiSecret := os.Getenv("VALR_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping VALR tests; VALR_API_KEY/VALR_API_SECRET not set")
	}

	restyClient := NewRestyClientForTests(apiKey, apiSecret)
	valr := NewValrClient(restyClient)

	// NOTE: This test actually places an order on your live/paper account.
	//       Make sure your key has minimal balance or is set to "paper" if possible.
	limReq := PlaceLimitOrderRequest{
		CurrencyPair:    "BTCZAR",
		Side:            "BUY",
		Price:           "1", // obviously very low; likely it will stay OPEN
		Quantity:        "0.0001",
		CustomerOrderID: "testOrder123",
	}
	limResp, err := valr.PlaceLimitOrder(limReq)
	if err != nil {
		t.Fatalf("PlaceLimitOrder failed: %v", err)
	}
	if limResp.Status != "NEW" && limResp.Status != "PARTIALLY_FILLED" {
		t.Logf("Order returned status %q; proceeding to cancel anyway", limResp.Status)
	}

	// Cancel the order immediately
	err = valr.CancelOrder(limResp.OrderID)
	if err != nil {
		t.Errorf("CancelOrder failed: %v", err)
	}
}

// NewRestyClientForTests is a helper to initialize a Resty client with BasicAuth.
func NewRestyClientForTests(key, secret string) *resty.Client {
	return resty.New().
		SetHostURL("https://api.valr.com").
		SetBasicAuth(key, secret).
		SetHeader("Content-Type", "application/json").
		SetTimeout(10 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)
}
