package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/fcvgraan/arbitrage_bot/http_client"
)

func main() {
	// Load API credentials from environment
	apiKey := os.Getenv("LUNO_KEY")
	apiSecret := os.Getenv("LUNO_SECRET")
	if apiKey == "" || apiSecret == "" {
		log.Fatal("Please set LUNO_KEY and LUNO_SECRET environment variables")
	}

	// Create a Resty client with a 10s timeout and simple retry policy
	restyClient := resty.New().
		SetTimeout(10*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1*time.Second).
		SetBasicAuth(apiKey, apiSecret).
		SetHeader("Content-Type", "application/json")

	// Initialize our Luno client wrapper
	luno := http_client.NewLunoClient(restyClient)

	// Example: Get the BTC/ZAR ticker
	tickerResp, err := luno.GetTicker("XBTZAR")
	if err != nil {
		log.Fatalf("Error fetching ticker: %v", err)
	}
	fmt.Printf("Ticker XBTZAR: Bid=%s, Ask=%s, Last=%s\n",
		tickerResp.Bid, tickerResp.Ask, tickerResp.Last)

	// Example: Get your account balances (private endpoint)
	balances, err := luno.GetBalances()
	if err != nil {
		log.Fatalf("Error fetching balances: %v", err)
	}
	fmt.Println("Your balances:")
	for _, bal := range balances {
		fmt.Printf("- %s: Available=%s, Reserved=%s\n",
			bal.Currency, bal.Available, bal.Reserved)
	}
}
