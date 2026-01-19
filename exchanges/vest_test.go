package exchanges

import "testing"

func TestConvertVestSymbol(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{in: "tonusdt", want: "TON-PERP"},
		{in: "TON-USD", want: "TON-PERP"},
		{in: "TON/USD", want: "TON-PERP"},
		{in: "TON-USDT", want: "TON-PERP"},
		{in: "TONUSDC", want: "TON-PERP"},
		{in: "TON-PERP", want: "TON-PERP"},
		{in: "BTCUSDT", want: ""},
		{in: "", want: ""},
	}

	for _, tc := range cases {
		if got := convertToVestSymbol(tc.in); got != tc.want {
			t.Fatalf("convertToVestSymbol(%q): got %q want %q", tc.in, got, tc.want)
		}
	}

	fromCases := []struct {
		in   string
		want string
	}{
		{in: "TON-PERP", want: "TONUSDT"},
		{in: "TON-PERPUSDT", want: "TONUSDT"},
		{in: "TONPERP", want: "TONUSDT"},
		{in: "TON-USD", want: "TONUSDT"},
		{in: "BTC-PERP", want: ""},
	}

	for _, tc := range fromCases {
		if got := convertFromVestSymbol(tc.in); got != tc.want {
			t.Fatalf("convertFromVestSymbol(%q): got %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseVestDepthMessage(t *testing.T) {
	payload := []byte(`{"channel":"TON-PERP@depth","data":{"bids":[["100","1"],["99","2"]],"asks":[["101","1"],["102","2"]]}}`)
	gotSym, gotBid, gotAsk, _, ok := parseVestDepthMessage(payload)
	if !ok {
		t.Fatalf("expected ok")
	}
	if gotSym != "TONUSDT" {
		t.Fatalf("symbol: got %q", gotSym)
	}
	if gotBid != 100 {
		t.Fatalf("bid: got %v", gotBid)
	}
	if gotAsk != 101 {
		t.Fatalf("ask: got %v", gotAsk)
	}
}
