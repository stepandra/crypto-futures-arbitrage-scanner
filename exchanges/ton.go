package exchanges

import "strings"

const tonCanonicalSymbol = "TONUSDT"

func normalizeTONSymbol(symbol string) (string, bool) {
	u := strings.ToUpper(strings.TrimSpace(symbol))
	if u == "" {
		return "", false
	}

	u = strings.NewReplacer("/", "-", "_", "-", ":", "-").Replace(u)

	if u == "TON" {
		return tonCanonicalSymbol, true
	}
	if u == tonCanonicalSymbol {
		return tonCanonicalSymbol, true
	}

	if strings.HasPrefix(u, "TON") {
		suffix := strings.TrimPrefix(u, "TON")
		suffix = strings.TrimPrefix(suffix, "-")
		switch suffix {
		case "USDT", "USDC", "USD", "PERP", "PERPUSDT":
			return tonCanonicalSymbol, true
		}
	}

	parts := strings.Split(u, "-")
	if len(parts) >= 1 && parts[0] == "TON" {
		if len(parts) >= 2 {
			switch parts[1] {
			case "USD", "USDT", "USDC", "PERP", "PERPUSDT":
				return tonCanonicalSymbol, true
			}
		}
		if len(parts) >= 3 && parts[1] == "USD" && parts[2] == "PERP" {
			return tonCanonicalSymbol, true
		}
	}

	return "", false
}

func filterTONSymbols(symbols []string) []string {
	seen := make(map[string]struct{}, 1)
	out := make([]string, 0, 1)
	for _, s := range symbols {
		canon, ok := normalizeTONSymbol(s)
		if !ok {
			continue
		}
		if _, exists := seen[canon]; exists {
			continue
		}
		seen[canon] = struct{}{}
		out = append(out, canon)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
