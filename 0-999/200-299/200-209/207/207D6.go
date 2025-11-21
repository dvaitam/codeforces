package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

var keywordWeights map[string][3]float64

func init() {
	keywordWeights = make(map[string][3]float64)

	addWords(0, 3.5,
		"president", "presidential", "minister", "ministry", "prime", "premier", "government", "governor",
		"parliament", "congress", "senate", "senator", "deputy", "cabinet", "kremlin", "duma",
		"election", "electoral", "candidate", "campaign", "coalition", "opposition", "republic",
		"constitution", "constitutional", "referendum", "treaty", "summit", "negotiation", "agreement",
		"policy", "foreign", "diplomat", "diplomatic", "embassy", "ambassador", "consulate", "visa",
		"sanction", "customs", "border", "security", "defense", "defence", "army", "navy", "airforce",
		"military", "soldier", "troop", "brigade", "regiment", "commander", "general", "war", "conflict",
		"battle", "frontline", "occupation", "separatist", "terrorism", "terrorist", "extremist", "radical",
		"police", "prosecutor", "investigation", "arrest", "verdict", "trial", "court", "justice", "legal",
		"law", "laws", "legislative", "bill", "decree", "ordinance", "resolution", "mayor", "council",
		"municipal", "federal", "state", "republican", "democrat", "chancellor", "monarch", "king", "queen",
		"emperor", "empress", "prime-minister")

	addWords(0, 2.0,
		"rights", "human", "freedom", "civic", "activist", "protest", "demonstration", "rally",
		"security", "intelligence", "nato", "alliance", "bloc", "union", "strategy", "strategic",
		"border", "checkpoint", "migration", "refugee", "asylum", "occupy", "annexation", "diplomacy",
		"geopolitics", "embargo", "sanctions")

	addWords(2, 3.5,
		"market", "trade", "trading", "commerce", "commercial", "company", "companies", "business",
		"enterprise", "firm", "corporation", "holding", "concern", "factory", "plant", "workshop",
		"industry", "industrial", "production", "manufacture", "manufacturing", "product", "products",
		"supply", "demand", "logistics", "warehouse", "shipping", "transport", "delivery", "export",
		"exporter", "import", "importer", "tariff", "customs", "supply", "distribution", "pipeline",
		"terminal", "retail", "wholesale", "store", "shop", "sales", "sale", "discount", "client",
		"customer", "contract", "agreement", "deal", "merger", "acquisition", "ipo", "listing",
		"investment", "investor", "investors", "fund", "funds", "financing", "finance", "financial",
		"economy", "economic", "economics", "economist", "macro", "micro", "gdp", "inflation",
		"deflation", "recession", "crisis", "budget", "deficit", "surplus", "tax", "taxes", "taxation",
		"revenue", "profit", "profits", "loss", "losses", "earnings", "income", "salary", "wage",
		"wages", "dividend", "share", "shares", "shareholder", "stock", "stocks", "exchange", "bourse",
		"bond", "bonds", "securities", "futures", "option", "currency", "currencies", "dollar", "dollars",
		"euro", "yen", "yuan", "ruble", "rouble", "pound", "bitcoin", "crypto", "cryptocurrency",
		"blockchain", "wallet", "bank", "banks", "banking", "loan", "loans", "mortgage", "credit",
		"deposit", "deposits", "savings", "interest", "rate", "rates", "capital", "capitalization",
		"valuation", "payroll", "employment", "unemployment", "jobless", "labor", "labour", "strike",
		"union", "pension", "insurance", "premium", "policyholder", "commodity", "commodities",
		"metal", "metals", "steel", "aluminum", "copper", "nickel", "gold", "silver", "platinum",
		"oil", "gas", "fuel", "diesel", "petrol", "coal", "uranium", "energy", "power", "electricity",
		"grid", "generation", "refinery", "drilling", "field", "well", "rig", "seaport", "airport",
		"railway", "logistic", "shipping", "fleet", "vehicle", "auto", "automobile", "car", "cars",
		"truck", "bus", "plane", "aircraft", "airline", "ticket", "tourism", "hotel", "real", "estate",
		"property", "housing", "construction", "developer", "mortgage", "lease", "rental", "rent",
		"agriculture", "agricultural", "farmer", "farm", "harvest", "crop", "grain", "wheat", "corn",
		"barley", "soy", "sugar", "cotton", "livestock", "cattle", "poultry", "fishery", "fishing",
		"food", "beverage", "restaurant", "cafe")

	addWords(1, 3.0,
		"sport", "sports", "football", "soccer", "hockey", "basketball", "volleyball", "handball",
		"baseball", "rugby", "tennis", "table", "badminton", "golf", "boxing", "martial", "wrestling",
		"judo", "karate", "gymnastics", "athletics", "swimming", "skiing", "biathlon", "skating",
		"snowboard", "cycling", "race", "racing", "motor", "motorsport", "formula", "grand", "prix",
		"nascar", "rally", "marathon", "triathlon", "rowing", "sailing", "chess", "e-sport", "esport",
		"match", "matches", "game", "games", "fixture", "fixtures", "cup", "league", "championship",
		"champion", "champions", "tournament", "playoff", "season", "round", "stage", "group",
		"quarterfinal", "semifinal", "final", "title", "medal", "bronze", "silver", "gold", "record",
		"victory", "win", "wins", "defeat", "draw", "coach", "manager", "trainer", "team", "teams",
		"club", "clubs", "player", "players", "striker", "forward", "midfielder", "defender",
		"goalkeeper", "goaltender", "goal", "goals", "assist", "score", "scored", "points", "fans",
		"supporters", "stadium", "arena", "squad", "lineup", "injury", "contract", "transfer",
		"loaned", "referee", "umpire", "olympic", "olympics", "paralympic")

	addWords(1, 2.5,
		"science", "scientific", "scientist", "research", "study", "studies", "laboratory", "lab",
		"university", "universities", "college", "campus", "faculty", "academy", "institute",
		"professor", "academic", "student", "students", "school", "schools", "education", "educational",
		"lesson", "classroom", "teacher", "teachers", "pupil", "pupils", "curriculum", "exam", "exams",
		"olympiad", "scholarship", "innovation", "innovative", "technology", "technologies", "tech",
		"technological", "startup", "start-up", "it", "software", "hardware", "computer", "computers",
		"server", "cloud", "platform", "app", "apps", "application", "applications", "program",
		"programming", "algorithm", "artificial", "intelligence", "ai", "machine", "learning", "data",
		"analytics", "robot", "robotics", "automation", "digital", "online", "internet", "network",
		"cyber", "security", "hacker", "malware", "virus", "social", "media", "website", "portal",
		"browser", "email", "mobile", "smartphone", "tablet", "gadget", "device", "sensor",
		"nanotechnology", "biotech", "biotechnology", "genetics", "genome", "dna", "molecule",
		"physics", "chemistry", "biology", "astronomy", "astrophysics", "space", "cosmos", "satellite",
		"spacecraft", "rocket", "launch", "cosmodrome", "astronaut", "cosmonaut", "orbital", "planet",
		"galaxy", "telescope", "meteor", "asteroid", "probe", "medicine", "medical", "doctor",
		"doctors", "nurse", "nurses", "clinic", "hospital", "hospitals", "patient", "patients",
		"health", "healthcare", "treatment", "therapy", "surgery", "diagnosis", "symptom", "vaccine",
		"vaccination", "epidemic", "pandemic", "infection", "infectious", "virus", "covid", "flu",
		"disease", "illness", "microbe", "pharmacy", "drug", "drugs", "pharmaceutical")

	addWords(1, 2.0,
		"culture", "cultural", "art", "arts", "artist", "artists", "painting", "gallery", "museum",
		"heritage", "tradition", "festival", "concert", "music", "musical", "musician", "band",
		"singer", "song", "songs", "album", "orchestra", "symphony", "opera", "ballet", "theater",
		"theatre", "stage", "performance", "show", "series", "episode", "season", "cinema", "film",
		"films", "movie", "movies", "director", "producer", "screenplay", "premiere", "actor", "actress",
		"casting", "shooting", "literature", "book", "books", "novel", "poem", "poet", "writer",
		"author", "publishing", "magazine", "journal", "press", "media", "broadcast", "television",
		"radio", "newspaper", "interview", "blog", "fashion", "style", "designer", "design", "trend",
		"lifestyle", "family", "child", "children", "kid", "kids", "parent", "parents", "youth",
		"society", "community", "volunteer", "charity", "festival", "holiday")
}

func addWords(class int, weight float64, words ...string) {
	for _, w := range words {
		word := strings.ToLower(strings.ReplaceAll(w, " ", ""))
		if word == "" {
			continue
		}
		entry := keywordWeights[word]
		entry[class] += weight
		keywordWeights[word] = entry
	}
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return
	}
	rawText := string(data)
	lower := strings.ToLower(rawText)

	tokens := splitWords(lower)
	counts := make(map[string]int, len(tokens))
	for _, tok := range tokens {
		if tok == "" {
			continue
		}
		counts[tok]++
	}

	var scores [3]float64
	signal := false
	for tok, cnt := range counts {
		if weights, ok := lookupWeight(tok); ok {
			signal = true
			for i := 0; i < 3; i++ {
				if weights[i] != 0 {
					scores[i] += weights[i] * float64(cnt)
				}
			}
		}
	}

	signal = applyHeuristics(lower, rawText, counts, &scores) || signal

	best := 0
	for i := 1; i < 3; i++ {
		if scores[i] > scores[best]+1e-9 {
			best = i
		}
	}

	if !signal {
		fmt.Println(3)
		return
	}
	fmt.Println(best + 1)
}

func lookupWeight(token string) ([3]float64, bool) {
	if weight, ok := keywordWeights[token]; ok {
		return weight, true
	}
	if len(token) > 4 && strings.HasSuffix(token, "ies") {
		if weight, ok := keywordWeights[token[:len(token)-3]+"y"]; ok {
			return weight, true
		}
	}
	suffixes := []string{"ing", "ed", "es", "s"}
	for _, suf := range suffixes {
		if len(token) > len(suf)+2 && strings.HasSuffix(token, suf) {
			base := token[:len(token)-len(suf)]
			if weight, ok := keywordWeights[base]; ok {
				return weight, true
			}
		}
	}
	return [3]float64{}, false
}

func splitWords(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r))
	})
}

func applyHeuristics(lower, raw string, counts map[string]int, scores *[3]float64) bool {
	signal := false
	bump := func(class int, delta float64) {
		if delta == 0 {
			return
		}
		scores[class] += delta
		signal = true
	}

	phrases := []struct {
		text  string
		class int
		score float64
	}{
		{"prime minister", 0, 8},
		{"white house", 0, 5},
		{"state duma", 0, 4},
		{"human rights", 0, 3},
		{"world cup", 1, 6},
		{"champions league", 1, 6},
		{"grand slam", 1, 4},
		{"formula one", 1, 5},
		{"winter games", 1, 4},
		{"summer games", 1, 4},
		{"artificial intelligence", 1, 4},
		{"space station", 1, 4},
		{"space agency", 1, 3},
		{"central bank", 2, 6},
		{"stock exchange", 2, 6},
		{"oil price", 2, 5},
		{"gas price", 2, 5},
		{"percent", 2, 3},
		{"per cent", 2, 3},
		{"billion", 2, 2.5},
		{"million", 2, 2},
	}
	for _, p := range phrases {
		if strings.Contains(lower, p.text) {
			bump(p.class, p.score)
		}
	}

	sportsHits := countHits(counts,
		"match", "matches", "goal", "goals", "points", "assist", "team", "teams", "club", "clubs",
		"coach", "stadium", "fans", "season", "tournament", "cup", "league", "championship",
		"olympic", "medal", "race", "champion", "driver", "lap")
	if sportsHits > 0 {
		bump(1, float64(sportsHits))
	}

	scienceHits := countHits(counts,
		"science", "scientific", "research", "study", "technology", "technologies", "tech", "digital",
		"online", "internet", "computer", "software", "hardware", "startup", "robot", "space",
		"rocket", "satellite", "medicine", "medical", "doctor", "hospital", "vaccine", "virus",
		"education", "school", "student", "university")
	if scienceHits > 0 {
		bump(1, 0.7*float64(scienceHits))
	}

	economyHits := countHits(counts,
		"market", "trade", "company", "business", "investment", "investor", "bank", "loan", "credit",
		"budget", "deficit", "tax", "tariff", "profit", "loss", "revenue", "salary", "wage", "factory",
		"plant", "industry", "export", "import", "price", "oil", "gas", "energy", "currency",
		"dollar", "euro", "ruble", "rouble", "yuan", "yen", "bitcoin", "crypto")
	if economyHits > 0 {
		bump(2, 0.8*float64(economyHits))
	}

	digitCount := 0
	currencySymbols := 0
	for _, r := range raw {
		if unicode.IsDigit(r) {
			digitCount++
		}
		switch r {
		case '$', '€', '£', '¥', '₽':
			currencySymbols++
		}
	}
	if currencySymbols > 0 {
		bump(2, float64(currencySymbols)*4)
	}
	if digitCount > len(raw)/25 {
		bump(2, 3)
	}
	if strings.Contains(raw, "%") {
		bump(2, 2)
	}

	socialHits := countHits(counts,
		"family", "children", "child", "kids", "parents", "school", "education", "teacher", "hospital",
		"clinic", "patient", "health", "culture", "museum", "concert", "festival", "book", "author")
	if socialHits > 0 {
		bump(1, 0.6*float64(socialHits))
	}

	if strings.Contains(lower, "court") || strings.Contains(lower, "lawsuit") || strings.Contains(lower, "lawyer") {
		bump(0, 3)
	}
	if strings.Contains(lower, "mobilized") || strings.Contains(lower, "regiment") {
		bump(0, 2)
	}

	if strings.Contains(lower, "merger") || strings.Contains(lower, "acquisition") || strings.Contains(lower, "ipo") {
		bump(2, 5)
	}
	if strings.Contains(lower, "start-up") || strings.Contains(lower, "startup") {
		bump(2, 2.5)
	}

	return signal
}

func countHits(counts map[string]int, words ...string) int {
	total := 0
	for _, w := range words {
		total += counts[w]
	}
	return total
}
