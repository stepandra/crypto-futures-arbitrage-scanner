package exchanges

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type variationalMetadataResponse struct {
	Listings []struct {
		Ticker    string `json:"ticker"`
		MarkPrice string `json:"mark_price"`
		Quotes    struct {
			Size1k struct {
				Bid string `json:"bid"`
				Ask string `json:"ask"`
			} `json:"size_1k"`
		} `json:"quotes"`
	} `json:"listings"`
}

func convertToVariationalTicker(symbol string) string {
	if _, ok := normalizeTONSymbol(symbol); ok {
		return "TON"
	}
	return ""
}

func parseVariationalTopOfBook(meta variationalMetadataResponse, symbol string) (bestBid float64, bestAsk float64, ok bool) {
	ticker := convertToVariationalTicker(symbol)
	if ticker == "" {
		return 0, 0, false
	}

	for _, l := range meta.Listings {
		if strings.ToUpper(l.Ticker) != ticker {
			continue
		}

		bid, bidErr := strconv.ParseFloat(l.Quotes.Size1k.Bid, 64)
		ask, askErr := strconv.ParseFloat(l.Quotes.Size1k.Ask, 64)
		if bidErr == nil && askErr == nil && bid > 0 && ask > 0 {
			return bid, ask, true
		}
		return 0, 0, false
	}

	return 0, 0, false
}

func parseVariationalMark(meta variationalMetadataResponse, symbol string) (mark float64, ok bool) {
	ticker := convertToVariationalTicker(symbol)
	if ticker == "" {
		return 0, false
	}

	for _, l := range meta.Listings {
		if strings.ToUpper(l.Ticker) != ticker {
			continue
		}
		mark, err := strconv.ParseFloat(l.MarkPrice, 64)
		if err != nil || mark <= 0 {
			return 0, false
		}
		return mark, true
	}

	return 0, false
}

func ConnectVariationalFutures(symbols []string, priceChan chan<- PriceData, orderbookChan chan<- OrderbookData, tradeChan chan<- TradeData) {
	url := "https://omni-client-api.prod.ap-northeast-1.variational.io/metadata/stats"

	tonSyms := filterTONSymbols(symbols)
	if len(tonSyms) == 0 {
		log.Printf("Variational skipped: no TON symbol provided (symbols=%v)", symbols)
		return
	}

	supportedSymbols := []string{tonCanonicalSymbol}

	client := &http.Client{Timeout: 10 * time.Second}

	backoff := 2 * time.Second
	maxBackoff := 60 * time.Second

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Variational request build error: %v", err)
			time.Sleep(backoff)
			if backoff < maxBackoff {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Variational request error: %v (retrying in %s)", err, backoff)
			time.Sleep(backoff)
			if backoff < maxBackoff {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			log.Printf("Variational request error: unexpected status %s (retrying in %s)", resp.Status, backoff)
			time.Sleep(backoff)
			if backoff < maxBackoff {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
			continue
		}

		var meta variationalMetadataResponse
		err = json.NewDecoder(resp.Body).Decode(&meta)
		resp.Body.Close()
		if err != nil {
			log.Printf("Variational decode error: %v (retrying in %s)", err, backoff)
			time.Sleep(backoff)
			if backoff < maxBackoff {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
			continue
		}

		backoff = 2 * time.Second

		now := time.Now().UnixMilli()
		for _, sym := range supportedSymbols {
			bestBid, bestAsk, ok := parseVariationalTopOfBook(meta, sym)
			if ok {
				orderbookChan <- OrderbookData{Symbol: sym, Source: "variational_perps", BestBid: bestBid, BestAsk: bestAsk, Timestamp: now}
				continue
			}

			mark, ok := parseVariationalMark(meta, sym)
			if ok {
				priceChan <- PriceData{Symbol: sym, Source: "variational_perps", Price: mark, Timestamp: now}
			}
		}

		time.Sleep(2 * time.Second)
	}
}
