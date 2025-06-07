// main.go
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fcvgraan/arbitrage_bot/http_client/luno"
	"github.com/fcvgraan/arbitrage_bot/http_client/valr"
	"github.com/go-resty/resty/v2"
)

const (
	lunoFee  = 0.0025      // 0.25% taker fee on Luno
	valrFee  = 0.0010      // 0.10% taker fee on VALR
	symbol   = "XBTZAR"    // trading pair
	interval = time.Second // polling interval
)

func main() {
	// Load API credentials from environment
	lunoKey, lunoSec := os.Getenv("LUNO_KEY"), os.Getenv("LUNO_SECRET")
	valrKey, valrSec := os.Getenv("VALR_API_KEY"), os.Getenv("VALR_API_SECRET")
	if lunoKey == "" || lunoSec == "" || valrKey == "" || valrSec == "" {
		log.Fatal("Please set LUNO_KEY, LUNO_SECRET, VALR_API_KEY, VALR_API_SECRET")
	}

	// Initialize Resty clients
	rcL := resty.New().
		SetHostURL("https://api.luno.com").
		SetBasicAuth(lunoKey, lunoSec).
		SetTimeout(10 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)

	rcV := resty.New().
		SetHostURL("https://api.valr.com").
		SetBasicAuth(valrKey, valrSec).
		SetTimeout(10 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)

	// Instantiate exchange clients
	lClient := luno.NewLunoClient(rcL)
	vClient := valr.NewValrClient(rcV)

	// Open (or create) CSV file for output
	file, err := os.OpenFile("spread_data.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	// Write header if file is empty
	if fi, _ := file.Stat(); fi.Size() == 0 {
		header := []string{
			"timestamp",
			"luno_bid", "luno_bid_vol", "luno_ask", "luno_ask_vol",
			"valr_bid", "valr_bid_qty", "valr_ask", "valr_ask_qty",
			"raw_spread", "net_spread", "trade_vol", "potential_profit",
		}
		if err := writer.Write(header); err != nil {
			log.Fatalf("Error writing CSV header: %v", err)
		}
		writer.Flush()
	}

	ticker := time.NewTicker(interval)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-ticker.C:
			now := time.Now().UTC().Format(time.RFC3339)

			// 1) Fetch Luno order book
			lOb, err := lClient.GetOrderBook(symbol)
			if err != nil {
				log.Printf("Luno error: %v", err)
				continue
			}
			lBidP, _ := strconv.ParseFloat(lOb.Bids[0][0], 64)
			lBidV, _ := strconv.ParseFloat(lOb.Bids[0][1], 64)
			lAskP, _ := strconv.ParseFloat(lOb.Asks[0][0], 64)
			lAskV, _ := strconv.ParseFloat(lOb.Asks[0][1], 64)

			// 2) Fetch VALR order book
			vOb, err := vClient.GetOrderBook("BTCZAR", false)
			if err != nil {
				log.Printf("VALR error: %v", err)
				continue
			}
			vBidP, _ := strconv.ParseFloat(vOb.Bids[0].Price, 64)
			vBidV, _ := strconv.ParseFloat(vOb.Bids[0].Quantity, 64)
			vAskP, _ := strconv.ParseFloat(vOb.Asks[0].Price, 64)
			vAskV, _ := strconv.ParseFloat(vOb.Asks[0].Quantity, 64)

			// 3) Compute spreads and fees
			rawSpread := vBidP - lAskP
			feeBuy := lAskP * lunoFee
			feeSell := vBidP * valrFee
			netSpread := rawSpread - (feeBuy + feeSell)

			// 4) Determine tradable volume
			tradeVol := lAskV
			if vBidV < tradeVol {
				tradeVol = vBidV
			}
			potentialProfit := netSpread * tradeVol

			// 5) Write a record to CSV
			record := []string{
				now,
				fmt.Sprintf("%.2f", lBidP),
				fmt.Sprintf("%.6f", lBidV),
				fmt.Sprintf("%.2f", lAskP),
				fmt.Sprintf("%.6f", lAskV),
				fmt.Sprintf("%.2f", vBidP),
				fmt.Sprintf("%.6f", vBidV),
				fmt.Sprintf("%.2f", vAskP),
				fmt.Sprintf("%.6f", vAskV),
				fmt.Sprintf("%.2f", rawSpread),
				fmt.Sprintf("%.2f", netSpread),
				fmt.Sprintf("%.6f", tradeVol),
				fmt.Sprintf("%.2f", potentialProfit),
			}
			if err := writer.Write(record); err != nil {
				log.Printf("Error writing record: %v", err)
			}
			writer.Flush()
		}
	}

	log.Println("Shutting down, CSV data saved.")
}
