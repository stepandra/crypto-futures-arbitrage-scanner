package exchanges

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type vestSubscribeRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

type vestPingRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     int           `json:"id"`
}

type vestDepthMessage struct {
	Channel string `json:"channel"`
	Data    struct {
		Bids [][]string `json:"bids"`
		Asks [][]string `json:"asks"`
	} `json:"data"`
}

func convertToVestSymbol(symbol string) string {
	if _, ok := normalizeTONSymbol(symbol); !ok {
		return ""
	}
	return "TON-PERP"
}

func convertFromVestSymbol(vestSymbol string) string {
	canon, ok := normalizeTONSymbol(vestSymbol)
	if !ok {
		return ""
	}
	return canon
}

func parseVestDepthMessage(payload []byte) (symbol string, bestBid float64, bestAsk float64, ts int64, ok bool) {
	var msg vestDepthMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return "", 0, 0, 0, false
	}

	channelSym := msg.Channel
	if idx := strings.Index(channelSym, "@"); idx != -1 {
		channelSym = channelSym[:idx]
	}

	symbol = convertFromVestSymbol(channelSym)
	if symbol == "" {
		return "", 0, 0, 0, false
	}

	bestBid = math.SmallestNonzeroFloat64
	for _, level := range msg.Data.Bids {
		if len(level) < 1 {
			continue
		}
		px, err := strconv.ParseFloat(level[0], 64)
		if err != nil {
			continue
		}
		if px > bestBid {
			bestBid = px
		}
	}

	bestAsk = math.MaxFloat64
	for _, level := range msg.Data.Asks {
		if len(level) < 1 {
			continue
		}
		px, err := strconv.ParseFloat(level[0], 64)
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

	ts = time.Now().UnixMilli()
	return symbol, bestBid, bestAsk, ts, true
}

func ConnectVestFutures(symbols []string, priceChan chan<- PriceData, orderbookChan chan<- OrderbookData, tradeChan chan<- TradeData) {
	wsURL := "wss://ws-prod.hz.vestmarkets.com/ws-api?version=1.0"

	tonSyms := filterTONSymbols(symbols)
	if len(tonSyms) == 0 {
		log.Printf("Vest skipped: no TON symbol provided (symbols=%v)", symbols)
		return
	}

	var params []string
	vestSym := convertToVestSymbol(tonCanonicalSymbol)
	if vestSym != "" {
		params = append(params, vestSym+"@depth")
	}
	if len(params) == 0 {
		log.Printf("Vest: unable to map TON to Vest market; sleeping")
		time.Sleep(30 * time.Second)
		return
	}

	for {
		headers := http.Header{}
		headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			if resp != nil {
				log.Printf("Vest connection error: %v, Status: %s", err, resp.Status)
			} else {
				log.Printf("Vest connection error: %v", err)
			}
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Connected to Vest WebSocket")

		subReq := vestSubscribeRequest{Method: "SUBSCRIBE", Params: params, ID: 1}
		if err := conn.WriteJSON(subReq); err != nil {
			log.Printf("Vest subscription error: %v", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}



		// Start Ping loop
		go func(c *websocket.Conn) {
			ticker := time.NewTicker(20 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					pingReq := vestPingRequest{Method: "PING", Params: []interface{}{}, ID: 0}
					if err := c.WriteJSON(pingReq); err != nil {
						return
					}
				}
			}
		}(conn)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Vest read error: %v", err)
				conn.Close()
				break
			}

			stdSymbol, bestBid, bestAsk, ts, ok := parseVestDepthMessage(msg)
			if !ok {
				continue
			}

			orderbookChan <- OrderbookData{
				Symbol:    stdSymbol,
				Source:    "vest_futures",
				BestBid:   bestBid,
				BestAsk:   bestAsk,
				Timestamp: ts,
			}
		}

		time.Sleep(2 * time.Second)
	}
}
