package exchanges

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const DeDustPoolsURL = "https://api.dedust.io/v2/pools"

type DeDustAsset struct {
	Type     string `json:"type"`
	Address  string `json:"address"`
	Metadata struct {
		Decimals int `json:"decimals"`
	} `json:"metadata"`
}

type DeDustPool struct {
	Assets   []DeDustAsset `json:"assets"`
	Reserves []string      `json:"reserves"`
}

func ConnectDeDust(priceChan chan<- PriceData) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	usdtAddr := "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"

	ticker := time.NewTicker(2 * time.Second) // Poll frequently (2s)
	defer ticker.Stop()

    log.Println("DeDust: Connected and polling for deepest liquidity pool...")

	for range ticker.C {
		resp, err := client.Get(DeDustPoolsURL)
		if err != nil {
			log.Printf("DeDust: Error fetching pools: %v", err)
			continue
		}

		var pools []DeDustPool
		if err := json.NewDecoder(resp.Body).Decode(&pools); err != nil {
			log.Printf("DeDust: Error decoding pools: %v", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

        var bestPrice float64
        var maxLiquidity float64
        found := false

		for _, p := range pools {
			if len(p.Assets) != 2 {
				continue
			}

			var tonReserve, usdtReserve float64
			isTon := false
			isUsdt := false

			for i, asset := range p.Assets {
				if asset.Type == "native" {
					isTon = true
					if len(p.Reserves) > i {
						r, _ := strconv.ParseFloat(p.Reserves[i], 64)
						tonReserve = r / 1e9 // TON has 9 decimals
					}
				} else if asset.Type == "jetton" && asset.Address == usdtAddr {
					isUsdt = true
					if len(p.Reserves) > i {
						r, _ := strconv.ParseFloat(p.Reserves[i], 64)
						usdtReserve = r / 1e6 // USDT has 6 decimals
					}
				}
			}

			if isTon && isUsdt && tonReserve > 0 {
                // Determine liquidity score (simple approximation: USdT reserve)
                liquidity := usdtReserve
                if liquidity > maxLiquidity {
                    maxLiquidity = liquidity
                    bestPrice = usdtReserve / tonReserve
                    found = true
                }
			}
		}

        if found {
            priceChan <- PriceData{
                Symbol:    "TONUSDT",
                Source:    "DeDust",
                Price:     bestPrice,
                Timestamp: time.Now().UnixMilli(),
            }
            // Optional: verbose log to prove updates
            // log.Printf("DeDust: Deepest Pool Price: %.4f (TVL: %.2f USDT)", bestPrice, maxLiquidity)
        }
	}
}
