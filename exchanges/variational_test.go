package exchanges

import (
	"encoding/json"
	"testing"
)

func TestConvertVariationalTicker(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{in: "TONUSDT", want: "TON"},
		{in: "ton-usdt", want: "TON"},
		{in: "TON-USD", want: "TON"},
		{in: "TON/USD", want: "TON"},
		{in: "TON-PERP", want: "TON"},
		{in: "BTCUSDT", want: ""},
		{in: "", want: ""},
	}
	for _, tc := range cases {
		if got := convertToVariationalTicker(tc.in); got != tc.want {
			t.Fatalf("convertToVariationalTicker(%q): got %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseVariationalTopOfBook(t *testing.T) {
	payload := []byte(`{"listings":[{"ticker":"TON","mark_price":"100","quotes":{"size_1k":{"bid":"99","ask":"101"}}}]}`)

	var meta variationalMetadataResponse
	if err := json.Unmarshal(payload, &meta); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	bid, ask, ok := parseVariationalTopOfBook(meta, "TONUSDT")
	if !ok {
		t.Fatalf("expected ok")
	}
	if bid != 99 || ask != 101 {
		t.Fatalf("unexpected bid/ask: %v/%v", bid, ask)
	}
}

func TestParseVariationalMark(t *testing.T) {
	payload := []byte(`{"listings":[{"ticker":"TON","mark_price":"123.45","quotes":{"size_1k":{"bid":"0","ask":"0"}}}]}`)

	var meta variationalMetadataResponse
	if err := json.Unmarshal(payload, &meta); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	mark, ok := parseVariationalMark(meta, "TONUSDT")
	if !ok {
		t.Fatalf("expected ok")
	}
	if mark != 123.45 {
		t.Fatalf("mark: got %v", mark)
	}
}
