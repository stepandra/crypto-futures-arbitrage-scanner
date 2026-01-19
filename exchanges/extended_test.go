package exchanges

import "testing"

func TestConvertExtendedMarket(t *testing.T) {
	if got := convertToExtendedMarket("TONUSDT"); got == "" {
		t.Fatalf("expected non-empty market")
	}
	if got := convertToExtendedMarket("TON-USD"); got == "" {
		t.Fatalf("expected non-empty market")
	}
	if got := convertToExtendedMarket("BTCUSDT"); got != "" {
		t.Fatalf("expected empty market, got %q", got)
	}

	fromCases := []struct {
		in   string
		want string
	}{
		{in: "TON-USD", want: "TONUSDT"},
		{in: "TON-USDT", want: "TONUSDT"},
		{in: "TON-USDC", want: "TONUSDT"},
		{in: "TON-PERP", want: "TONUSDT"},
		{in: "BTC-USD", want: ""},
	}
	for _, tc := range fromCases {
		if got := convertFromExtendedMarket(tc.in); got != tc.want {
			t.Fatalf("convertFromExtendedMarket(%q): got %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseExtendedOrderbookMessage(t *testing.T) {
	cases := []struct {
		name   string
		market string
	}{
		{name: "usd", market: "TON-USD"},
		{name: "usdt", market: "TON-USDT"},
		{name: "usdc", market: "TON-USDC"},
		{name: "perp", market: "TON-PERP"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := []byte(`{"type":"SNAPSHOT","data":{"m":"` + tc.market + `","b":[{"p":"2.10","q":"1"}],"a":[{"p":"2.20","q":"2"}]},"ts":123,"seq":1}`)
			sym, bid, ask, ts, ok := parseExtendedOrderbookMessage(payload)
			if !ok {
				t.Fatalf("expected ok")
			}
			if sym != "TONUSDT" {
				t.Fatalf("symbol: got %q", sym)
			}
			if bid != 2.10 {
				t.Fatalf("bid: got %v", bid)
			}
			if ask != 2.20 {
				t.Fatalf("ask: got %v", ask)
			}
			if ts != 123 {
				t.Fatalf("ts: got %v", ts)
			}
		})
	}
}
