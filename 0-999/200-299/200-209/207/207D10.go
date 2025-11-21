package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

var cultureKeywords = map[string]float64{
	"art": 3, "arts": 2, "artist": 3, "artists": 3,
	"artistic": 2.5, "painting": 3.5, "paintings": 3.5, "gallery": 3,
	"galleries": 3, "museum": 3, "museums": 3, "culture": 3.5,
	"cultural": 3.5, "heritage": 2.5, "music": 3, "musical": 3,
	"musician": 3, "musicians": 3, "concert": 3, "concerts": 3,
	"opera": 3.5, "ballet": 3, "theatre": 3, "theater": 3,
	"cinema": 2.5, "film": 2.5, "films": 2.5, "movie": 2.5,
	"movies": 2.5, "director": 2, "directors": 2, "festival": 2.5,
	"festivals": 2.5, "literature": 3.5, "literary": 3, "novel": 3,
	"novels": 3, "poet": 3, "poets": 3, "poetry": 3, "book": 2,
	"books": 2, "author": 2.5, "authors": 2.5, "story": 2,
	"stories": 2, "drama": 2.5, "dance": 2.5, "dancing": 2.5,
	"sculpture": 3, "sculptures": 3, "exhibition": 3, "exhibitions": 3,
	"choir": 2.5, "symphony": 3, "orchestra": 3, "composer": 3,
	"composers": 3, "design": 2, "fashion": 2.5, "style": 1.5,
	"photography": 2.5, "photo": 1.5, "photos": 1.5, "stage": 1.5,
	"screenplay": 2.5, "screenwriter": 2.5, "painting": 3.5,
}

var politicsKeywords = map[string]float64{
	"president": 4, "presidential": 4, "prime": 1.5, "minister": 4,
	"ministry": 3.5, "government": 4.5, "parliament": 4, "congress": 4,
	"senate": 3.5, "senator": 3, "state": 2, "states": 2,
	"policy": 3.5, "policies": 3.5, "politics": 4, "political": 4,
	"party": 3, "parties": 3, "election": 4.5, "elections": 4.5,
	"vote": 3.5, "voting": 3.5, "law": 2.5, "laws": 2.5, "constitutional": 3.5,
	"constitution": 3.5, "reform": 3, "authority": 2.5, "administration": 2.5,
	"official": 2.5, "officials": 2.5, "diplomat": 3, "diplomacy": 3.5,
	"diplomatic": 3.5, "foreign": 3, "relations": 2, "security": 3,
	"military": 3, "army": 3, "troops": 3, "soldiers": 2.5,
	"war": 3.5, "conflict": 3.5, "border": 3, "defense": 3.5,
	"defence": 3.5, "strategy": 2.5, "strategic": 2.5, "federal": 2.5,
	"council": 2.5, "commission": 2, "committee": 2, "republic": 2.5,
	"royal": 2, "king": 2.5, "queen": 2.5, "parliamentary": 3,
	"cabinet": 2.5, "governor": 2.5, "ambassador": 3, "campaign": 3.5,
	"coalition": 3, "democracy": 3.5, "democratic": 3.5, "republican": 3.5,
	"monarchy": 3, "opposition": 3.5, "regime": 3, "sanction": 4,
	"sanctions": 4, "agreement": 3, "treaty": 3, "negotiation": 3,
	"negotiations": 3, "national": 2, "nation": 2, "sovereignty": 3,
	"referendum": 3.5, "duma": 3, "kremlin": 4, "kreml": 4,
	"presidium": 3, "resolution": 2.5, "deputy": 2.5, "ministries": 3,
}

var tradeKeywords = map[string]float64{
	"trade": 5, "trading": 4.5, "market": 4.5, "markets": 4.5,
	"economy": 5, "economic": 4.5, "economics": 4.5, "financial": 4,
	"finance": 4, "bank": 4, "banks": 4, "banking": 4,
	"company": 3.5, "companies": 3.5, "business": 3.5, "businesses": 3.5,
	"enterprise": 3, "enterprises": 3, "industry": 3.5, "industrial": 3.5,
	"manufacture": 3.5, "manufacturing": 3.5, "factory": 3, "factories": 3,
	"production": 3.5, "product": 2.5, "products": 2.5, "export": 4,
	"exports": 4, "import": 4, "imports": 4, "investment": 4,
	"investments": 4, "investor": 3.5, "investors": 3.5, "capital": 3.5,
	"budget": 3.5, "budgets": 3.5, "gross": 2, "gdp": 4,
	"inflation": 4, "tax": 3.5, "taxes": 3.5, "tariff": 3.5,
	"tariffs": 3.5, "price": 3, "prices": 3, "cost": 2.5,
	"costs": 2.5, "sale": 2.5, "sales": 2.5, "retail": 3,
	"wholesale": 3, "profit": 3.5, "profits": 3.5, "loss": 2.5,
	"losses": 2.5, "revenue": 3.5, "revenues": 3.5, "income": 3,
	"incomes": 3, "loan": 3, "loans": 3, "credit": 3, "credits": 3,
	"currency": 3.5, "currencies": 3.5, "dollar": 3.5, "dollars": 3.5,
	"euro": 3.5, "euros": 3.5, "yen": 3, "pound": 3,
	"fund": 3, "funds": 3, "stock": 4, "stocks": 4, "exchange": 3,
	"bond": 3, "bonds": 3, "oil": 3.5, "gas": 3, "energy": 3.5,
	"fuel": 2.5, "power": 2.5, "raw": 2, "material": 2,
	"materials": 2, "resource": 2, "resources": 2, "agriculture": 3,
	"agricultural": 3, "crop": 2, "crops": 2, "harvest": 2,
	"shipping": 2.5, "logistics": 2.5, "transport": 2.5, "supply": 3,
	"demand": 3, "merger": 3.5, "acquisition": 3.5, "share": 3,
	"shares": 3, "shareholder": 3, "dividend": 3, "valuation": 3,
	"insurance": 2.5, "economist": 3, "analyst": 3,
}

type fragment struct {
	subject int
	part    string
	weight  float64
}

var fragments = []fragment{
	{0, "art house", 4}, {0, "folk art", 3}, {0, "short story", 2.5},
	{0, "film festival", 4}, {0, "literary prize", 4},
	{1, "prime minister", 5}, {1, "head of state", 4},
	{1, "foreign ministry", 4}, {1, "security council", 4},
	{1, "armed forces", 3.5}, {1, "ruling party", 4},
	{1, "opposition leader", 4}, {1, "peace talks", 3.5},
	{1, "human rights", 3.5}, {1, "state duma", 4},
	{2, "stock exchange", 5}, {2, "central bank", 5},
	{2, "foreign investment", 4}, {2, "economic growth", 4},
	{2, "price index", 3.5}, {2, "gross domestic product", 5},
	{2, "oil and gas", 4}, {2, "supply chain", 3.5},
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return
	}
	text := strings.ToLower(string(data))
	tokens := tokenize(text)

	scores := [3]float64{0.5, 0.5, 1} // light prior for trade-related docs

	for _, token := range tokens {
		if w, ok := cultureKeywords[token]; ok {
			scores[0] += w
		}
		if w, ok := politicsKeywords[token]; ok {
			scores[1] += w
		}
		if w, ok := tradeKeywords[token]; ok {
			scores[2] += w
		}
		if looksNumeric(token) {
			scores[2] += 0.4
		}
		if strings.HasSuffix(token, "ism") || strings.HasSuffix(token, "ist") {
			scores[1] += 0.6
		}
		if strings.HasSuffix(token, "ing") && len(token) > 6 {
			scores[2] += 0.2
		}
	}

	for _, frag := range fragments {
		if strings.Contains(text, frag.part) {
			scores[frag.subject] += frag.weight
		}
	}

	if strings.Contains(text, "%") || strings.Contains(text, "$") {
		scores[2] += 1.5
	}

	result := 0
	for i := 1; i < 3; i++ {
		if scores[i] > scores[result] {
			result = i
		}
	}
	fmt.Println(result + 1)
}

func tokenize(s string) []string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
	}
	return strings.Fields(b.String())
}

func looksNumeric(token string) bool {
	if len(token) == 0 {
		return false
	}
	for _, r := range token {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
