package http_client

import (
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

// These tests require real Luno API keys in environment:
// export LUNO_KEY="your_key_id"
// export LUNO_SECRET="your_key_secret"
func TestGetTicker(t *testing.T) {
	apiKey := os.Getenv("LUNO_KEY")
	apiSecret := os.Getenv("LUNO_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping Luno tests; LUNO_KEY/LUNO_SECRET not set")
	}

	// Initialize a Resty client just like in main.go
	restyClient := NewRestyClientForTests(apiKey, apiSecret)
	luno := NewLunoClient(restyClient)

	// Fetch ticker XBTZAR
	ticker, err := luno.GetTicker("XBTZAR")
	if err != nil {
		t.Fatalf("GetTicker failed: %v", err)
	}
	if ticker.Pair != "XBTZAR" {
		t.Errorf("Expected pair XBTZAR, got %s", ticker.Pair)
	}
	if ticker.Bid == "" || ticker.Ask == "" {
		t.Errorf("Bid or Ask was empty: bid=%s ask=%s", ticker.Bid, ticker.Ask)
	}
	t.Logf("Got ticker: Bid=%s Ask=%s", ticker.Bid, ticker.Ask)
}

func TestGetBalances(t *testing.T) {
	apiKey := os.Getenv("LUNO_KEY")
	apiSecret := os.Getenv("LUNO_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping Luno tests; LUNO_KEY/LUNO_SECRET not set")
	}

	restyClient := NewRestyClientForTests(apiKey, apiSecret)
	luno := NewLunoClient(restyClient)

	balances, err := luno.GetBalances()
	if err != nil {
		t.Fatalf("GetBalances failed: %v", err)
	}
	if len(balances) == 0 {
		t.Errorf("Expected at least one balance, got zero")
	}
	for _, bal := range balances {
		if bal.Currency == "" {
			t.Errorf("Balance had empty currency field")
		}
	}
}

// NewRestyClientForTests is a helper to initialize a Resty client with BasicAuth.
func NewRestyClientForTests(key, secret string) *resty.Client {
	client := resty.New().
		SetBasicAuth(key, secret).
		SetHeader("Content-Type", "application/json").
		SetTimeout(10 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)

	return client
}
