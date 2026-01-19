package exchanges

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type extendedOrderbookEnvelope struct {
	Type string `json:"type"`
	Data struct {
		Market     string `json:"m"`
		MarketLong string `json:"market"`
		Bids       []struct {
			Price string `json:"p"`
			Qty   string `json:"q"`
		} `json:"b"`
		Asks []struct {
			Price string `json:"p"`
			Qty   string `json:"q"`
		} `json:"a"`
	} `json:"data"`
	TS  int64 `json:"ts"`
	Seq int64 `json:"seq"`
}

type extendedMarketsResponse struct {
	Status string `json:"status"`
	Data   []struct {
		Name            string          `json:"name"`
		AssetName       string          `json:"assetName"`
		CollateralAsset string          `json:"collateralAssetName"`
		TradingStatus   string          `json:"status"`
		UIName          string          `json:"uiName"`
		MarketStats     json.RawMessage `json:"marketStats"`
		TradingConfig   json.RawMessage `json:"tradingConfig"`
	} `json:"data"`
	Error json.RawMessage `json:"error"`
}

func convertToExtendedMarket(symbol string) string {
	if _, ok := normalizeTONSymbol(symbol); ok {
		// Fallback market identifier if discovery fails.
		return "TON-USD"
	}
	return ""
}

func convertFromExtendedMarket(market string) string {
	canon, ok := normalizeTONSymbol(market)
	if !ok {
		return ""
	}
	return canon
}

func parseExtendedOrderbookMessage(payload []byte) (symbol string, bestBid float64, bestAsk float64, ts int64, ok bool) {
	var env extendedOrderbookEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		log.Printf("Extended parse error: %v", err)
		return "", 0, 0, 0, false
	}

	market := env.Data.Market
	if market == "" {
		market = env.Data.MarketLong
	}

	stdSymbol := convertFromExtendedMarket(market)
	if stdSymbol == "" {
		log.Printf("Extended unknown market: %s or %s", env.Data.Market, env.Data.MarketLong)
		return "", 0, 0, 0, false
	}

	bestBid = math.SmallestNonzeroFloat64
	for _, lvl := range env.Data.Bids {
		px, err := strconv.ParseFloat(lvl.Price, 64)
		if err != nil {
			continue
		}
		if px > bestBid {
			bestBid = px
		}
	}

	bestAsk = math.MaxFloat64
	for _, lvl := range env.Data.Asks {
		px, err := strconv.ParseFloat(lvl.Price, 64)
		if err != nil {
			continue
		}
		if px < bestAsk {
			bestAsk = px
		}
	}

	if bestBid == math.SmallestNonzeroFloat64 || bestAsk == math.MaxFloat64 {
		return "", 0, 0, 0, false
	}

	if env.TS > 0 {
		ts = env.TS
	} else {
		ts = time.Now().UnixMilli()
	}

	return stdSymbol, bestBid, bestAsk, ts, true
}

func extendedUserAgent() string {
	return "crypto-futures-arbitrage-scanner/1.0"
}

func fetchExtendedMarkets(restBaseURL string) (map[string]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := strings.TrimRight(restBaseURL, "/") + "/info/markets"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", extendedUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var decoded extendedMarketsResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	m := make(map[string]string)
	bestMarket := ""
	bestRank := -1
	for _, row := range decoded.Data {
		asset := strings.ToUpper(row.AssetName)
		collateral := strings.ToUpper(row.CollateralAsset)
		if asset != "TON" {
			continue
		}

		rank := -1
		switch collateral {
		case "USDT":
			rank = 30
		case "USDC":
			rank = 20
		case "USD":
			rank = 10
		}
		if rank < 0 {
			continue
		}

		if rank > bestRank {
			bestRank = rank
			bestMarket = row.Name
		}
	}

	if bestMarket != "" {
		m[tonCanonicalSymbol] = bestMarket
	}

	return m, nil
}

func connectExtendedMarketOrderbook(stdSymbol, market, wsBaseURL string, orderbookChan chan<- OrderbookData) {
	headers := http.Header{}
	headers.Set("User-Agent", extendedUserAgent())

	wsURL := fmt.Sprintf("%s/orderbooks/%s?depth=1", strings.TrimRight(wsBaseURL, "/"), market)

	backoff := 2 * time.Second
	maxBackoff := 60 * time.Second

	for {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			log.Printf("Extended connection error (%s/%s): %v (retrying in %s)", stdSymbol, market, err, backoff)
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
		log.Printf("Connected to Extended orderbook stream (%s/%s)", stdSymbol, market)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Extended read error (%s/%s): %v", stdSymbol, market, err)
				conn.Close()
				break
			}
			// log.Printf("Extended received: %s", string(msg))

			parsedSymbol, bestBid, bestAsk, ts, ok := parseExtendedOrderbookMessage(msg)
			if !ok {
				continue
			}
			if parsedSymbol == "" {
				parsedSymbol = stdSymbol
			}

			orderbookChan <- OrderbookData{
				Symbol:    parsedSymbol,
				Source:    "extended_futures",
				BestBid:   bestBid,
				BestAsk:   bestAsk,
				Timestamp: ts,
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func ConnectExtendedFutures(symbols []string, priceChan chan<- PriceData, orderbookChan chan<- OrderbookData, tradeChan chan<- TradeData) {
	restBaseURL := "https://api.starknet.extended.exchange/api/v1"
	wsBaseURL := "wss://api.starknet.extended.exchange/stream.extended.exchange/v1"

	tonSyms := filterTONSymbols(symbols)
	if len(tonSyms) == 0 {
		log.Printf("Extended skipped: no TON symbol provided (symbols=%v)", symbols)
		return
	}

	marketMap, err := fetchExtendedMarkets(restBaseURL)
	if err != nil {
		log.Printf("Extended market discovery error: %v", err)
		marketMap = make(map[string]string)
	}

	market := marketMap[tonCanonicalSymbol]
	if market == "" {
		market = convertToExtendedMarket(tonCanonicalSymbol)
	}
	if market == "" {
		log.Printf("Extended: unable to resolve TON market; sleeping")
		time.Sleep(30 * time.Second)
		return
	}

	go connectExtendedMarketOrderbook(tonCanonicalSymbol, market, wsBaseURL, orderbookChan)

	select {}
}
