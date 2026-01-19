package exchanges

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type lighterOrderBooksResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	OrderBooks []struct {
		Symbol     string `json:"symbol"`
		MarketID   int    `json:"market_id"`
		MarketType string `json:"market_type"`
		Status     string `json:"status"`
	} `json:"order_books"`
}

type lighterSubscribeMessage struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Auth    string `json:"auth,omitempty"`
}

type lighterOrderBookUpdate struct {
	Channel   string `json:"channel"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	OrderBook struct {
		Asks []struct {
			Price string `json:"price"`
			Size  string `json:"size"`
		} `json:"asks"`
		Bids []struct {
			Price string `json:"price"`
			Size  string `json:"size"`
		} `json:"bids"`
	} `json:"order_book"`
}

func convertToLighterSymbol(symbol string) string {
	if _, ok := normalizeTONSymbol(symbol); ok {
		return "TON"
	}
	return ""
}

func parseLighterMarketID(channel string) int {
	channel = strings.TrimSpace(channel)
	if channel == "" {
		return -1
	}
	if strings.HasPrefix(channel, "order_book:") {
		idStr := strings.TrimPrefix(channel, "order_book:")
		id, err := strconv.Atoi(idStr)
		if err == nil {
			return id
		}
	}
	if strings.HasPrefix(channel, "order_book/") {
		idStr := strings.TrimPrefix(channel, "order_book/")
		id, err := strconv.Atoi(idStr)
		if err == nil {
			return id
		}
	}
	return -1
}

func parseLighterOrderBookMessage(payload []byte) (marketID int, bestBid float64, bestAsk float64, ts int64, ok bool) {
	var msg lighterOrderBookUpdate
	if err := json.Unmarshal(payload, &msg); err != nil {
		log.Printf("Lighter parse error: %v", err)
		return 0, 0, 0, 0, false
	}
	if msg.Type != "update/order_book" {
		return 0, 0, 0, 0, false
	}

	marketID = parseLighterMarketID(msg.Channel)
	if marketID < 0 {
		return 0, 0, 0, 0, false
	}

	bestBid = math.SmallestNonzeroFloat64
	for _, lvl := range msg.OrderBook.Bids {
		px, err := strconv.ParseFloat(lvl.Price, 64)
		if err != nil {
			continue
		}
		if px > bestBid {
			bestBid = px
		}
	}

	bestAsk = math.MaxFloat64
	for _, lvl := range msg.OrderBook.Asks {
		px, err := strconv.ParseFloat(lvl.Price, 64)
		if err != nil {
			continue
		}
		if px < bestAsk {
			bestAsk = px
		}
	}

	if bestBid == math.SmallestNonzeroFloat64 || bestAsk == math.MaxFloat64 {
		return 0, 0, 0, 0, false
	}

	if msg.Timestamp > 0 {
		ts = msg.Timestamp
	} else {
		ts = time.Now().UnixMilli()
	}

	return marketID, bestBid, bestAsk, ts, true
}

func fetchLighterMarketMap(baseURL string) (map[string]int, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := strings.TrimRight(baseURL, "/") + "/api/v1/orderBooks"
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var decoded lighterOrderBooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	symbolToID := make(map[string]int)
	for _, ob := range decoded.OrderBooks {
		if strings.ToLower(ob.MarketType) != "perp" {
			continue
		}
		if strings.ToLower(ob.Status) != "active" {
			continue
		}
		s := strings.ToUpper(strings.TrimSpace(ob.Symbol))
		if s == "" {
			continue
		}

		base := s
		if parts := strings.SplitN(base, "-", 2); len(parts) > 0 {
			base = parts[0]
		}
		if parts := strings.SplitN(base, "/", 2); len(parts) > 0 {
			base = parts[0]
		}
		if base == "" {
			continue
		}

		if base != "TON" {
			continue
		}

		symbolToID[base] = ob.MarketID
		symbolToID[s] = ob.MarketID
	}

	return symbolToID, nil
}

func lighterAuthToken() string {
	if v := strings.TrimSpace(os.Getenv("LIGHTER_READONLY_AUTH")); v != "" {
		return v
	}
	if v := strings.TrimSpace(os.Getenv("LIGHTER_AUTH")); v != "" {
		return v
	}
	return ""
}

func ConnectLighterFutures(symbols []string, priceChan chan<- PriceData, orderbookChan chan<- OrderbookData, tradeChan chan<- TradeData) {
	restBaseURL := "https://mainnet.zklighter.elliot.ai"
	wsURL := "wss://mainnet.zklighter.elliot.ai/stream"

	tonSyms := filterTONSymbols(symbols)
	if len(tonSyms) == 0 {
		log.Printf("Lighter skipped: no TON symbol provided (symbols=%v)", symbols)
		return
	}

	symbolToID, err := fetchLighterMarketMap(restBaseURL)
	if err != nil {
		log.Printf("Lighter market map fetch error: %v", err)
		symbolToID = map[string]int{}
	}

	selectedIDs := make(map[int]string)
	id, ok := symbolToID["TON"]
	if !ok {
		id, ok = symbolToID["TON-USDT"]
	}
	if ok {
		selectedIDs[id] = tonCanonicalSymbol
	}
	if len(selectedIDs) == 0 {
		log.Printf("Lighter: unable to resolve TON market id; sleeping")
		time.Sleep(30 * time.Second)
		return
	}

	auth := lighterAuthToken()

	for {
		headers := http.Header{}
		headers.Set("User-Agent", "crypto-futures-arbitrage-scanner/1.0")
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			log.Printf("Lighter connection error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Set a handler for ping messages to ensure we reply with pong


		// Set a handler for ping messages to ensure we reply with pong
		conn.SetPingHandler(func(appData string) error {
			// log.Printf("Lighter: Received ping") 
			return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
		})

		log.Printf("Connected to Lighter WebSocket")

		for id := range selectedIDs {
			sub := lighterSubscribeMessage{Type: "subscribe", Channel: fmt.Sprintf("order_book/%d", id)}
			if auth != "" {
				sub.Auth = auth
			}
			if err := conn.WriteJSON(sub); err != nil {
				log.Printf("Lighter subscribe error (%d): %v", id, err)
			}
		}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Lighter read error: %v", err)
				conn.Close()
				break
			}


			marketID, bestBid, bestAsk, ts, ok := parseLighterOrderBookMessage(msg)
			if !ok {
				continue
			}

			stdSymbol := selectedIDs[marketID]
			if stdSymbol == "" {
				continue
			}

			orderbookChan <- OrderbookData{Symbol: stdSymbol, Source: "lighter_futures", BestBid: bestBid, BestAsk: bestAsk, Timestamp: ts}
		}

		time.Sleep(2 * time.Second)
	}
}
