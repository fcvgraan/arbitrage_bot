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
	lunoFee  = 0.0025 // 0.25% taker fee on Luno
	valrFee  = 0.0010 // 0.10% taker fee on VALR
	interval = time.Second
)

func main() {
	lKey, lSecret := os.Getenv("LUNO_KEY"), os.Getenv("LUNO_SECRET")
	vKey, vSecret := os.Getenv("VALR_API_KEY"), os.Getenv("VALR_API_SECRET")
	if lKey == "" || lSecret == "" || vKey == "" || vSecret == "" {
		log.Fatal("Set LUNO_KEY, LUNO_SECRET, VALR_API_KEY, VALR_API_SECRET")
	}

	rcL := resty.New().
		SetBaseURL("https://api.luno.com").
		SetBasicAuth(lKey, lSecret).
		SetTimeout(10 * time.Second).
		SetRetryCount(3).SetRetryWaitTime(1 * time.Second)

	rcV := resty.New().
		SetBaseURL("https://api.valr.com").
		SetBasicAuth(vKey, vSecret).
		SetTimeout(10 * time.Second).
		SetRetryCount(3).SetRetryWaitTime(1 * time.Second)

	lClient := luno.NewLunoClient(rcL)
	vClient := valr.NewValrClient(rcV)

	pairs := []struct {
		asset    string
		lunoPair string
		valrPair string
	}{
		{"BTC", "XBTZAR", "BTCZAR"},
		{"ETH", "ETHZAR", "ETHZAR"},
		{"XRP", "XRPZAR", "XRPZAR"},
	}

	oppFile, err := os.OpenFile("opportunities.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Opening opportunities CSV: %v", err)
	}
	defer oppFile.Close()
	oppW := csv.NewWriter(oppFile)
	if fi, _ := oppFile.Stat(); fi.Size() == 0 {
		oppW.Write([]string{
			"asset",
			"timestamp",
			"direction",
			"price_buy", "vol_buy",
			"price_sell", "vol_sell",
			"net_spread", "trade_vol", "profit",
		})
		oppW.Flush()
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
			ts := time.Now().UTC().Format(time.RFC3339)

			for _, p := range pairs {
				lob, err := lClient.GetOrderBook(p.lunoPair)
				if err != nil || len(lob.Bids) == 0 || len(lob.Asks) == 0 {
					continue
				}
				lBidP, _ := strconv.ParseFloat(lob.Bids[0].Price, 64)
				lBidV, _ := strconv.ParseFloat(lob.Bids[0].Volume, 64)
				lAskP, _ := strconv.ParseFloat(lob.Asks[0].Price, 64)
				lAskV, _ := strconv.ParseFloat(lob.Asks[0].Volume, 64)

				vob, err := vClient.GetOrderBook(p.valrPair)
				if err != nil || len(vob.Bids) == 0 || len(vob.Asks) == 0 {
					continue
				}
				vBidP, _ := strconv.ParseFloat(vob.Bids[0].Price, 64)
				vBidV, _ := strconv.ParseFloat(vob.Bids[0].Quantity, 64)
				vAskP, _ := strconv.ParseFloat(vob.Asks[0].Price, 64)
				vAskV, _ := strconv.ParseFloat(vob.Asks[0].Quantity, 64)

				// Forward: buy on Luno (ask) → sell on VALR (bid)
				rawFwd := vBidP - lAskP
				netFwd := rawFwd - (lAskP*lunoFee + vBidP*valrFee)
				volFwd := lAskV
				if vBidV < volFwd {
					volFwd = vBidV
				}
				profFwd := netFwd * volFwd
				if profFwd > 0 {
					oppW.Write([]string{
						p.asset,
						ts,
						"Luno→VALR",
						fmt.Sprintf("%.2f", lAskP),
						fmt.Sprintf("%.6f", volFwd),
						fmt.Sprintf("%.2f", vBidP),
						fmt.Sprintf("%.6f", volFwd),
						fmt.Sprintf("%.2f", netFwd),
						fmt.Sprintf("%.6f", volFwd),
						fmt.Sprintf("%.2f", profFwd),
					})
					oppW.Flush()
				}

				// Reverse: buy on VALR (ask) → sell on Luno (bid)
				rawRev := lBidP - vAskP
				netRev := rawRev - (vAskP*valrFee + lBidP*lunoFee)
				volRev := vAskV
				if lBidV < volRev {
					volRev = lBidV
				}
				profRev := netRev * volRev
				if profRev > 0 {
					oppW.Write([]string{
						p.asset,
						ts,
						"VALR→Luno",
						fmt.Sprintf("%.2f", vAskP),
						fmt.Sprintf("%.6f", volRev),
						fmt.Sprintf("%.2f", lBidP),
						fmt.Sprintf("%.6f", volRev),
						fmt.Sprintf("%.2f", netRev),
						fmt.Sprintf("%.6f", volRev),
						fmt.Sprintf("%.2f", profRev),
					})
					oppW.Flush()
				}
			}
		}
	}

	log.Println("Done — opportunities.csv updated.")
}
