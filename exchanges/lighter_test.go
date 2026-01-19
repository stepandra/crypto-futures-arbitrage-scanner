package exchanges

import "testing"

func TestConvertToLighterSymbol(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{in: "tonusdt", want: "TON"},
		{in: "TON-USD", want: "TON"},
		{in: "TON/USD", want: "TON"},
		{in: "TON-USDC", want: "TON"},
		{in: "TON-PERP", want: "TON"},
		{in: "BTCUSDT", want: ""},
		{in: "", want: ""},
	}

	for _, tc := range cases {
		if got := convertToLighterSymbol(tc.in); got != tc.want {
			t.Fatalf("convertToLighterSymbol(%q): got %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseLighterOrderBookMessage(t *testing.T) {
	payload := []byte(`{"channel":"order_book:0","offset":1,"order_book":{"code":0,"asks":[{"price":"101","size":"1"}],"bids":[{"price":"100","size":"2"}],"offset":1,"nonce":1,"begin_nonce":1},"timestamp":10,"type":"update/order_book"}`)
	marketID, bid, ask, ts, ok := parseLighterOrderBookMessage(payload)
	if !ok {
		t.Fatalf("expected ok")
	}
	if marketID != 0 {
		t.Fatalf("marketID: got %d", marketID)
	}
	if bid != 100 {
		t.Fatalf("bid: got %v", bid)
	}
	if ask != 101 {
		t.Fatalf("ask: got %v", ask)
	}
	if ts != 10 {
		t.Fatalf("ts: got %v", ts)
	}
}

func TestParseLighterMarketID(t *testing.T) {
	if got := parseLighterMarketID("order_book:123"); got != 123 {
		t.Fatalf("colon: got %d", got)
	}
	if got := parseLighterMarketID("order_book/456"); got != 456 {
		t.Fatalf("slash: got %d", got)
	}
}
