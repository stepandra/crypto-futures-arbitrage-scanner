package exchanges

import (
	"reflect"
	"testing"
)

func TestNormalizeTONSymbol(t *testing.T) {
	cases := []struct {
		in     string
		want   string
		wantOK bool
	}{
		{in: "TONUSDT", want: tonCanonicalSymbol, wantOK: true},
		{in: "tonusdt", want: tonCanonicalSymbol, wantOK: true},
		{in: "TON", want: tonCanonicalSymbol, wantOK: true},
		{in: "TON-USDT", want: tonCanonicalSymbol, wantOK: true},
		{in: "TON/USDT", want: tonCanonicalSymbol, wantOK: true},
		{in: "TON_USDT", want: tonCanonicalSymbol, wantOK: true},
		{in: "TON-USD", want: tonCanonicalSymbol, wantOK: true},
		{in: "TON/USDC", want: tonCanonicalSymbol, wantOK: true},
		{in: "TON-PERP", want: tonCanonicalSymbol, wantOK: true},
		{in: "TONPERP", want: tonCanonicalSymbol, wantOK: true},
		{in: "TON-USD-PERP", want: tonCanonicalSymbol, wantOK: true},
		{in: "BTCUSDT", want: "", wantOK: false},
		{in: "", want: "", wantOK: false},
		{in: "   ", want: "", wantOK: false},
	}

	for _, tc := range cases {
		got, ok := normalizeTONSymbol(tc.in)
		if ok != tc.wantOK {
			t.Fatalf("normalizeTONSymbol(%q): ok=%v want %v", tc.in, ok, tc.wantOK)
		}
		if got != tc.want {
			t.Fatalf("normalizeTONSymbol(%q): got %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestFilterTONSymbols(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{name: "no-ton", in: []string{"BTCUSDT", "ETHUSDT"}, want: nil},
		{name: "canonical-present", in: []string{"BTCUSDT", "TONUSDT", "ETHUSDT"}, want: []string{tonCanonicalSymbol}},
		{name: "alias-present", in: []string{"BTCUSDT", "TON-USD", "ETHUSDT"}, want: []string{tonCanonicalSymbol}},
		{name: "many-aliases-dedupe", in: []string{"TON-USD", "TON/USDT", "TONUSDT"}, want: []string{tonCanonicalSymbol}},
		{name: "handles-extra-aliases", in: []string{"TON-USDC", "TON-USD-PERP"}, want: []string{tonCanonicalSymbol}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := filterTONSymbols(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("filterTONSymbols(%v): got %v want %v", tc.in, got, tc.want)
			}
		})
	}
}
