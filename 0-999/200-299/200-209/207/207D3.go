package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var keywordWeights map[string][3]float64

func init() {
	keywordWeights = make(map[string][3]float64)

	addWords([]string{
		"president", "prime", "minister", "government", "parliament", "senate", "duma",
		"kremlin", "administration", "governor", "mayor", "cabinet", "spokesman",
		"council", "authority", "official", "diplomat", "embassy", "foreign", "policy",
		"politics", "political", "opposition", "coalition", "congress", "white", "house",
		"election", "campaign", "vote", "voter", "constitution", "referendum",
		"justice", "court", "judge", "attorney", "prosecutor", "law", "legal", "bill",
		"ministry", "defense", "security", "army", "military", "nato", "alliance",
		"kremlin's", "federal", "state", "regional", "republic", "senator", "ambassador",
	}, 0, 2.6)

	addWords([]string{
		"sanction", "sanctions", "treaty", "negotiation", "talks", "agreement",
		"delegation", "summit", "embargo", "resolution", "referendum",
	}, 0, 1.6)

	addWords([]string{
		"match", "matches", "game", "games", "team", "club", "coach", "manager", "player",
		"forward", "defender", "goalkeeper", "captain", "lineup", "substitution", "tactics",
		"stadium", "arena", "fans", "supporters", "crowd", "spectators",
		"goal", "goals", "score", "scored", "scorer", "assist", "penalty", "shootout",
		"win", "won", "victory", "defeat", "draw", "tie", "season", "league", "cup",
		"championship", "tournament", "playoff", "series", "round", "semifinal", "final",
		"athlete", "athletes", "runner", "sprinter", "marathon", "cyclist", "skier",
		"race", "racing", "start", "finish", "lap",
		"hockey", "football", "soccer", "basketball", "volleyball", "baseball", "golf",
		"handball", "rugby", "tennis", "grand", "slam", "set", "serve", "ace",
		"boxing", "fighter", "fight", "rounds", "knockout", "ufc", "mma", "wrestling",
		"medal", "olympic", "olympics", "games", "world", "cup",
	}, 1, 2.4)

	addWords([]string{
		"coach", "contract", "transfer", "loan", "signing", "training", "injury", "injured",
		"physio", "referee", "umpire", "pitch", "court", "rink", "ice", "ball", "puck",
	}, 1, 1.5)

	addWords([]string{
		"concert", "festival", "tour", "performance", "premiere", "theater", "theatre",
		"opera", "ballet", "orchestra", "choir", "artist", "painter", "sculptor",
		"museum", "gallery", "exhibition", "culture", "cultural", "literature",
		"author", "novelist", "book", "novel", "poet", "song", "album", "singer", "band",
	}, 1, 1.7)

	addWords([]string{
		"market", "markets", "trading", "trade", "business", "company", "companies", "firm",
		"corporation", "corporate", "holding", "group", "conglomerate", "enterprise",
		"startup", "venture", "investor", "investors", "investment", "investments",
		"fund", "funds", "capital", "capitalization", "ipo", "share", "shares", "stock",
		"stocks", "bond", "bonds", "dividend", "portfolio", "index", "indices",
		"finance", "financial", "bank", "banker", "banking", "credit", "loan", "debt",
		"mortgage", "insurance", "pension",
		"budget", "deficit", "surplus", "fiscal", "tax", "taxes", "tariff", "tariffs",
		"economy", "economic", "economics", "macroeconomic", "gdp", "cpi", "inflation",
		"unemployment", "employment", "wage", "wages", "salary", "salaries",
		"profit", "profits", "loss", "losses", "revenue", "revenues", "turnover",
		"income", "earnings", "forecast", "outlook",
		"price", "prices", "pricing", "cost", "costs", "expenses",
		"oil", "gas", "energy", "pipeline", "refinery", "petrochemical",
		"metallurgy", "steel", "coal", "mining", "extraction",
		"agriculture", "farmer", "farmers", "harvest", "crop", "grain", "wheat",
		"retail", "wholesale", "store", "chain", "mall", "supermarket",
		"logistics", "shipping", "freight", "port", "terminal", "warehouse",
		"export", "exports", "import", "imports", "customs", "duty",
		"currency", "currencies", "ruble", "rubles", "rubl", "dollar", "dollars",
		"euro", "yen", "yuan", "pound", "franc", "peso", "bitcoin", "crypto",
		"airline", "flight", "tickets", "hotel", "tourism", "travel", "tourist",
	}, 2, 2.8)

	addWords([]string{
		"factory", "plant", "production", "manufacturing", "assembly", "output",
		"capacity", "supply", "demand", "logistics", "delivery", "shipment",
		"merger", "acquisition", "takeover", "buyout", "partnership", "contract",
		"licence", "license", "royalty",
	}, 2, 1.8)

	addWords([]string{
		"inc", "corp", "co", "ltd", "llc", "plc", "sa", "ag", "nv", "ab",
	}, 2, 1.4)

	addSpecificWord("gdp", [3]float64{-0.4, -0.4, 4.5})
	addSpecificWord("cpi", [3]float64{-0.3, -0.3, 3.2})
	addSpecificWord("ipo", [3]float64{-0.3, -0.3, 3.5})
	addSpecificWord("ooo", [3]float64{-0.2, -0.2, 1.5})
}

func addWords(words []string, class int, weight float64) {
	for _, w := range words {
		word := strings.ToLower(w)
		entry := keywordWeights[word]
		for i := 0; i < 3; i++ {
			if i == class {
				entry[i] += weight
			} else {
				entry[i] -= weight * 0.12
			}
		}
		keywordWeights[word] = entry
	}
}

func addSpecificWord(word string, values [3]float64) {
	lw := strings.ToLower(word)
	entry := keywordWeights[lw]
	for i := 0; i < 3; i++ {
		entry[i] += values[i]
	}
	keywordWeights[lw] = entry
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	reader.Buffer(make([]byte, 0, 1024), 4*1024*1024)
	var builder strings.Builder
	for reader.Scan() {
		builder.WriteString(reader.Text())
		builder.WriteByte('\n')
	}
	text := builder.String()
	tokens := tokenize(text)
	scores := evaluate(tokens, text)

	best := 0
	for i := 1; i < 3; i++ {
		if scores[i] > scores[best]+1e-9 {
			best = i
		}
	}
	fmt.Println(best + 1)
}

func tokenize(s string) []string {
	var tokens []string
	var current strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) {
			current.WriteRune(unicode.ToLower(r))
		} else if unicode.IsDigit(r) {
			if current.Len() > 0 && (current.Len() < 3 || containsDigit(current.String())) {
				current.WriteRune(r)
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
				current.WriteRune(r)
			}
		} else {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

func containsDigit(word string) bool {
	for _, r := range word {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func evaluate(tokens []string, raw string) [3]float64 {
	scores := [3]float64{-2.0, -2.0, -2.0}
	freq := make(map[string]int)
	for _, tok := range tokens {
		freq[tok]++
	}

	for word, count := range freq {
		if weight, ok := keywordWeights[word]; ok {
			for i := 0; i < 3; i++ {
				scores[i] += weight[i] * float64(count)
			}
		}
		if strings.HasSuffix(word, "tion") && len(word) > 6 {
			scores[0] -= 0.05 * float64(count)
			scores[2] += 0.1 * float64(count)
		}
		if strings.HasSuffix(word, "ing") && len(word) > 5 {
			scores[1] += 0.02 * float64(count)
		}
		if word == "vs" || word == "vs." {
			scores[1] += 1.2 * float64(count)
		}
		if word == "percent" || word == "percentage" {
			scores[2] += 1.1 * float64(count)
		}
		if word == "peace" || word == "conflict" {
			scores[0] += 0.8 * float64(count)
		}
	}

	bonus := heuristics(raw, freq)
	for i := 0; i < 3; i++ {
		scores[i] += bonus[i]
	}
	return scores
}

func heuristics(raw string, freq map[string]int) [3]float64 {
	var bonus [3]float64
	lower := strings.ToLower(raw)

	if containsCurrencySymbol(lower) {
		bonus[2] += 4.0
	}
	if strings.Contains(lower, "billion") || strings.Contains(lower, "million") || strings.Contains(lower, "mln") || strings.Contains(lower, "bn") {
		bonus[2] += 1.5
	}
	if strings.Contains(lower, "percent") || strings.Contains(lower, "%") {
		bonus[2] += 0.8
	}
	if hasScorePattern(lower) {
		bonus[1] += 2.5
	}
	if strings.Contains(lower, "lineup") || strings.Contains(lower, "fixtures") {
		bonus[1] += 1.2
	}
	if strings.Contains(lower, "diplomat") || strings.Contains(lower, "kremlin") || strings.Contains(lower, "ministry") {
		bonus[0] += 1.5
	}
	if freq["budget"] > 0 || freq["tax"] > 0 {
		bonus[2] += 0.9
	}
	if freq["summit"] > 0 || freq["talks"] > 0 {
		bonus[0] += 0.8
	}
	if freq["stadium"] > 0 || freq["fans"] > 0 {
		bonus[1] += 0.8
	}
	return bonus
}

func containsCurrencySymbol(text string) bool {
	const symbols = "$€£¥₽₴₺₹"
	for _, r := range text {
		if strings.ContainsRune(symbols, r) {
			return true
		}
	}
	return strings.Contains(text, "usd") ||
		strings.Contains(text, "eur") ||
		strings.Contains(text, "rub") ||
		strings.Contains(text, "rmb") ||
		strings.Contains(text, "cny")
}

func hasScorePattern(text string) bool {
	for i := 1; i+1 < len(text); i++ {
		if text[i] == ':' && isDigit(text[i-1]) && isDigit(text[i+1]) {
			return true
		}
		if text[i] == '-' && isDigit(text[i-1]) && isDigit(text[i+1]) {
			return true
		}
	}
	return false
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
